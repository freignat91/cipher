package rsa

import (
	"crypto/rand"
	"fmt"
	"io"
	"io/ioutil"
	"math/big"
	"os"
	"strings"
	"time"
)

type PublicKey struct {
	nn *big.Int
	ee *big.Int
}

type PrivateKey struct {
	nn *big.Int
	dd *big.Int
}

func CreateRSAKey(keyBitSize int, verbose bool, debug bool) (*PublicKey, *PrivateKey, error) {
	if keyBitSize%64 != 0 {
		return nil, nil, fmt.Errorf("number of bits should be a multiple of 64")
	}
	pSize := keyBitSize / 2
	if verbose {
		fmt.Printf("Compute RSA keys size: %d bits\n", pSize*2)
	}

	//find two random prime size sqrt(keyBitsize)
	p1 := big.NewInt(0)
	p2 := big.NewInt(0)
	found := 0

	go func() {
		for p1.BitLen() != pSize {
			p1 = GetRandomPrime(pSize, verbose, debug)
			if p1.BitLen() == pSize {
				found++
				if verbose {
					fmt.Printf("prime (%dbits): %s\n", p1.BitLen(), p1)
				}
			}
			if verbose && p2.BitLen() != pSize {
				fmt.Println("Bad one, recompute it")
			}
		}
	}()

	go func() {
		for p2.BitLen() != pSize {
			p2 = GetRandomPrime(pSize, verbose, debug)
			if p2.BitLen() == pSize {
				found++
				if verbose {
					fmt.Printf("prime (%dbits): %s\n", p2.BitLen(), p2)
				}
			}
			if verbose && p2.BitLen() != pSize {
				fmt.Println("Bad one, recompute it")
			}
		}
	}()

	for found < 2 {
		time.Sleep(1 * time.Second)
	}

	//Compute n
	nn := big.NewInt(0)
	nn.Mul(p1, p2)
	if verbose {
		fmt.Printf("size=%d: n=%s\n", nn.BitLen(), nn)
	}

	//compute phi
	phi := big.NewInt(0)
	phi.Mul(p1.Sub(p1, one), p2.Sub(p2, one))
	if verbose {
		fmt.Printf("phi=%s\n", phi)
	}

	//compute e
	//ee := big.NewInt(13)
	ee := GetRandom(keyBitSize / 4)
	tmp := big.NewInt(0)
	for {
		ee = GetNextPrime(ee, verbose, false)
		if verbose {
			fmt.Printf("test e=%s\n", ee)
		}
		if tmp.Mod(phi, ee).Cmp(zero) != 0 {
			break
		}
	}
	if verbose {
		fmt.Printf("e=%s\n", ee)
	}

	//dd := InverseModulo(ee, phi)
	dd := big.NewInt(0)
	dd.ModInverse(ee, phi)
	if verbose {
		fmt.Printf("d=%s\n", dd)
	}
	return &PublicKey{nn: nn, ee: ee}, &PrivateKey{nn: nn, dd: dd}, nil
}

func EncryptFile(sourcePath string, targetPath string, keyPath string) error {
	publicKey, errp := GetPublicKey(keyPath)
	if errp != nil {
		return errp
	}
	bufferSize := publicKey.nn.BitLen()/8 - 1
	//fmt.Printf("key size: %d\n", bufferSize+1)
	filei, errf := os.OpenFile(sourcePath, os.O_RDWR, 0666)
	if errf != nil {
		return errf
	}
	defer filei.Close()
	fileo, errf := os.Create(targetPath)
	if errf != nil {
		return errf
	}
	defer fileo.Close()
	data := make([]byte, bufferSize, bufferSize)
	nn := 0
	lastN := 0
	for {
		data = data[:cap(data)]
		n, err := filei.Read(data)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		data = data[:n]
		//fmt.Printf("%d: read: (%d):%v\n\n", nn, len(data), data)
		datac, _ := publicKey.Encrypt(data, bufferSize+1)
		//fmt.Printf("%d: enc: (%d):%v\n", nn, len(datac), datac)
		if _, err := fileo.Write(datac); err != nil {
			return err
		}
		lastN = n
		nn++
	}
	//fmt.Printf("end: %d\n", lastN)
	if _, err := fileo.Write([]byte{byte(lastN % 256), byte(lastN / 256)}); err != nil {
		return err
	}
	return nil
}

func DecryptFile(sourcePath string, targetPath string, keyPath string) error {
	privateKey, errp := GetPrivateKey(keyPath)
	if errp != nil {
		return errp
	}
	bufferSize := privateKey.nn.BitLen()/8 - 1
	//fmt.Printf("key size: %d\n", bufferSize+1)
	filei, errf := os.OpenFile(sourcePath, os.O_RDWR, 0666)
	if errf != nil {
		return errf
	}
	defer filei.Close()
	fileo, errf := os.Create(targetPath)
	if errf != nil {
		return errf
	}
	defer fileo.Close()
	prevData := make([]byte, bufferSize+1, bufferSize+1)
	data := make([]byte, bufferSize+1, bufferSize+1)
	datac := make([]byte, 0, 0)
	nn := 0
	for {
		data = data[:cap(data)]
		n, err := filei.Read(data)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		data = data[:n]
		//fmt.Printf("%d: read: (%d):%v\n", nn, len(data), data)
		if nn > 0 {
			datac, _ = privateKey.Decrypt(prevData, bufferSize)
			if len(data) == 2 {
				slen := int(data[0]) + int(data[1])*256
				datac = datac[bufferSize-slen:]
				//fmt.Printf("%d: dec-end: (%d):%v\n\n", nn-1, len(datac), datac)
			}
			//fmt.Printf("%d: write: (%d):%v\n\n", nn-1, len(datac), datac)
			if _, err := fileo.Write(datac); err != nil {
				return err
			}
		}
		prevData = prevData[:n]
		for i, val := range data {
			prevData[i] = val
		}
		nn++
	}
	return nil
}

func GetKeys(path string) (*PublicKey, *PrivateKey, error) {
	publicKey, erru := GetPublicKey(fmt.Sprintf("%s.pub", path))
	if erru != nil {
		return nil, nil, erru
	}
	privateKey, erri := GetPrivateKey(fmt.Sprintf("%s.key", path))
	if erri != nil {
		return nil, nil, erri
	}
	return publicKey, privateKey, nil
}

func SaveKeys(path string, publicKey *PublicKey, privateKey *PrivateKey) error {
	if err := ioutil.WriteFile(fmt.Sprintf("%s.pub", path), []byte(publicKey.ToHexa()), 0666); err != nil {
		return err
	}
	if err := ioutil.WriteFile(fmt.Sprintf("%s.key", path), []byte(privateKey.ToHexa()), 0666); err != nil {
		return err
	}
	return nil
}

func (k *PublicKey) ToHexa() string {
	return fmt.Sprintf("%x-%x", k.nn, k.ee)
}

func (k *PrivateKey) ToHexa() string {
	return fmt.Sprintf("%x-%x", k.nn, k.dd)
}

func GetPublicKey(path string) (*PublicKey, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	keyl := strings.Split(string(data), "-")
	nn := big.NewInt(0)
	fmt.Sscanf(keyl[0], "%x", nn)
	ee := big.NewInt(0)
	fmt.Sscanf(keyl[1], "%x", ee)
	if ee.Cmp(zero) == 0 || nn.Cmp(zero) == 0 {
		return nil, fmt.Errorf("Error reading public key")
	}
	key := &PublicKey{
		nn: nn,
		ee: ee,
	}
	return key, nil
}

func (k *PublicKey) GetRSAKeySize() int {
	return k.nn.BitLen()
}

func (k *PublicKey) Encrypt(data []byte, size int) ([]byte, error) {
	if len(data) > k.nn.BitLen()/8 {
		return nil, fmt.Errorf("data too large. It's can exceed %d bytes for this key", k.nn.BitLen()/8)
	}
	//fmt.Printf("enc data=%d size=%d\n", len(data), size)
	tmp := big.NewInt(0)
	tmp.SetBytes(data)
	ee := big.NewInt(0)
	ee.Abs(k.ee)
	nn := big.NewInt(0)
	nn.Abs(k.nn)
	tmp = PowModulo(tmp, ee, nn)
	dec := tmp.Bytes()
	if len(dec) < size {
		dif := size - len(dec)
		data := make([]byte, size, size)
		for i, _ := range dec {
			data[i+dif] = dec[i]
		}
		//fmt.Printf("enc ret data=%d\n", len(data))
		return data, nil
	}
	//fmt.Printf("enc ret dec=%d\n", len(dec))
	return dec, nil
}

func (k *PrivateKey) GetRSAKeySize() int {
	return k.nn.BitLen()
}

func (k *PrivateKey) Decrypt(data []byte, size int) ([]byte, error) {
	if len(data) > k.nn.BitLen()/8 {
		return nil, fmt.Errorf("data too large. It's can exceed %d bytes for this key", k.nn.BitLen()/8)
	}
	//fmt.Printf("dec data=%d size=%d\n", len(data), size)
	tmp := big.NewInt(0)
	tmp.SetBytes(data)
	dd := big.NewInt(0)
	dd.Abs(k.dd)
	nn := big.NewInt(0)
	nn.Abs(k.nn)
	tmp = PowModulo(tmp, dd, nn)
	dec := tmp.Bytes()

	if len(dec) < size {
		dif := size - len(dec)
		data := make([]byte, size, size)
		for i, _ := range dec {
			data[i+dif] = dec[i]
		}
		return data, nil
		//fmt.Printf("dec ret data=%d\n", len(data))
	}

	//fmt.Printf("dec ret dec=%d\n", len(dec))
	return dec, nil
}

func GetPrivateKey(path string) (*PrivateKey, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	keyl := strings.Split(string(data), "-")
	nn := big.NewInt(0)
	fmt.Sscanf(keyl[0], "%x", nn)
	dd := big.NewInt(0)
	fmt.Sscanf(keyl[1], "%x", dd)
	if dd.Cmp(zero) == 0 || nn.Cmp(zero) == 0 {
		return nil, fmt.Errorf("Error reading private key")
	}
	key := &PrivateKey{
		nn: nn,
		dd: dd,
	}
	return key, nil
}

func TestKeySaveReload() {
	publicKey, privateKey, err := CreateRSAKey(256, false, false)
	if err != nil {
		fmt.Printf("Error creating keys: %v\n", err)
		return
	}
	if err := SaveKeys("test", publicKey, privateKey); err != nil {
		fmt.Printf("Error saving keys: %v\n", err)
		return
	}
	fmt.Printf("public-nn:  %x\n", publicKey.nn)
	fmt.Printf("private-nn: %x\n", privateKey.nn)
	fmt.Printf("public-ee:  %x\n", publicKey.ee)
	fmt.Printf("private-dd: %x\n", privateKey.dd)
	publicKey2, privateKey2, errg := GetKeys("test")
	if errg != nil {
		fmt.Printf("Error reading keys: %v\n", err)
		return
	}
	fmt.Println("-------------------------------------------------------------------")
	fmt.Printf("public-nn:  %x\n", publicKey2.nn)
	fmt.Printf("private-nn: %x\n", privateKey2.nn)
	fmt.Printf("public-ee:  %x\n", publicKey2.ee)
	fmt.Printf("private-dd: %x\n", privateKey2.dd)
	if publicKey2.nn.Cmp(publicKey.nn) != 0 {
		fmt.Printf("Error public key nn not equal: %x != %x\n", publicKey2.nn, publicKey.nn)
		return
	}
	if publicKey2.ee.Cmp(publicKey.ee) != 0 {
		fmt.Printf("Error public key ee not equal: %x != %x\n", publicKey2.ee, publicKey.ee)
		return
	}
	if privateKey2.nn.Cmp(privateKey.nn) != 0 {
		fmt.Printf("Error private key nn not equal: %x != %x\n", privateKey2.nn, privateKey.nn)
		return
	}
	if privateKey2.dd.Cmp(privateKey.dd) != 0 {
		fmt.Printf("Error private key dd not equal: %x != %x\n", privateKey2.dd, privateKey.dd)
		return
	}
	fmt.Println("keys ok")
}

func Test2() {
	//n := NewRandom(2048)
	n := NewDecimal("134525465745756822344523543563467456756734534535643455675678679880548658956787753245634575678978908790097943261452344467468789808790890706646523463456745845696789456734562465785697576894576345523465546785687689658465234524363567468764967365324564665456875664567458656787535635687456734654982763834563651763937")
	fmt.Printf("Size (bits): %d\n", n.BitLen())
	t0 := time.Now()
	GetNextPrime(n, false, false)
	fmt.Printf("time=%d ms\n", time.Now().Sub(t0).Nanoseconds()/1000000)
	fmt.Printf("prime=%s\n", n)
}

func TestEncriptDecript2() {
	fmt.Println("Load keys")
	//publicKey, privateKey, err := CreateRSAKey(256, false)
	publicKey, privateKey, err := GetKeys("k256")
	if err != nil {
		fmt.Printf("Error reading key: %v\n", err)
		return
	}
	size := publicKey.nn.BitLen()/8 - 1
	list := make([]byte, size, size)
	fmt.Printf("start test: size: %d\n", size)
	nn := 0
	for {
		rand.Read(list)
		//fmt.Printf("list: %v\n", list)
		c, _ := publicKey.Encrypt(list, size+1)
		/*
			if len(c) != size+1 {
				fmt.Printf("Error encrypted data size: %d\n", len(c))
				fmt.Printf("list: %v\n", list)
				fmt.Printf(" c: %v\n", c)
				return
			}
		*/
		d, _ := privateKey.Decrypt(c, size)
		if len(list) != len(d) {
			fmt.Printf("Error plain versus decrypt data size\n")
			fmt.Printf("list: %v\n", list)
			fmt.Printf(" c: %v\n", c)
			fmt.Printf(" d: %v\n", d)
			return
		}
		for i, val := range list {
			if i >= len(d) || val != d[i] {
				fmt.Println("Error")
				fmt.Printf("list: %v\n", list)
				fmt.Printf(" c: %v\n", c)
				fmt.Printf(" d: %v\n", d)
				return
			}
		}
		nn++
		if nn%1000 == 0 {
			fmt.Println(nn)
		}
	}
}

func TestEncriptDecript() {
	fmt.Println("Load keys")
	//publicKey, privateKey, err := CreateRSAKey(256, false)
	publicKey, privateKey, err := GetKeys("k256")
	if err != nil {
		fmt.Printf("Error reading key: %v\n", err)
		return
	}
	size := publicKey.nn.BitLen()/8 - 1
	list := []byte{15, 211, 218, 155, 207, 209, 212, 102, 241, 192, 130, 92, 10, 92, 213, 236, 172, 190, 189, 213, 116, 66, 8, 33, 132, 16, 66, 8, 33, 132, 16}
	//list := []byte{200, 166, 141, 215, 211, 66, 55, 245, 7, 183, 83, 15, 192, 57, 118, 110, 186, 145, 209, 127, 28, 54, 180, 68, 13, 122, 155, 105, 108, 110, 239}
	fmt.Printf("list %d: %v\n", len(list), list)
	nn := 0
	//fmt.Printf("list: %v\n", list)
	c, _ := publicKey.Encrypt(list, size+1)
	/*
		if len(c) != size+1 {
			fmt.Printf("Error encrypted data size: %d\n", len(c))
			fmt.Printf("list: %v\n", list)
			fmt.Printf(" c: %v\n", c)
			return
		}
	*/
	fmt.Printf("%d: enc %v\n", len(c), c)
	d, _ := privateKey.Decrypt(c, size)
	fmt.Printf("%d: dec %v\n", len(d), d)
	if len(list) != len(d) {
		fmt.Printf("Error plain versus decrypt data size\n")
		fmt.Printf("list: %v\n", list)
		fmt.Printf(" c: %v\n", c)
		fmt.Printf(" d: %v\n", d)
		return
	}
	for i, val := range list {
		if i >= len(d) || val != d[i] {
			fmt.Println("Error")
			fmt.Printf("list: %v\n", list)
			fmt.Printf(" c: %v\n", c)
			fmt.Printf(" d: %v\n", d)
			return
		}
	}
	nn++
	if nn%1000 == 0 {
		fmt.Println(nn)
	}

}
