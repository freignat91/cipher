package rsa

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"time"
)

var (
	zero = big.NewInt(0)
	one  = big.NewInt(1)
	two  = big.NewInt(2)
)

func NewDecimal(decimal string) *big.Int {
	nn := big.NewInt(0)
	fmt.Sscan(decimal, nn)
	return nn
}

func PowModulo(b *big.Int, e *big.Int, m *big.Int) *big.Int {
	pp := big.NewInt(1)
	for i, j := 0, e.BitLen(); i < j; i++ {
		if e.Bit(i) != 0 {
			pp.Mod(pp.Mul(pp, b), m)
		}
		b.Mod(b.Mul(b, b), m)
	}
	return pp
}

func PowModulo2(n *big.Int, exp *big.Int, nn *big.Int) *big.Int {
	if nn.Cmp(zero) == 0 {
		return n
	}
	pp := big.NewInt(1)
	nb := big.NewInt(0)
	nb.Add(n, zero)
	expb := big.NewInt(0)
	expb.Add(exp, zero)
	for expb.Cmp(zero) > 0 {
		if expb.Bit(0) != 0 {
			pp.Mul(pp, nb)
			pp.Mod(pp, nn)
		}
		nb.Mul(nb, nb)
		nb.Mod(nb, nn)
		expb.Div(expb, two)
	}
	return pp
}

func GetRandom(size int) *big.Int {
	n := big.NewInt(0)
	if size == 0 {
		return n
	}
	for n.BitLen() != size {
		b := make([]byte, size/8)
		rand.Read(b)
		n.SetBytes(b)
	}
	return n
}

func GetRandomPrime(size int, verbose bool, debug bool) *big.Int {
	n := GetRandom(size)
	p := GetNextPrime(n, verbose, debug)
	return p
}

func GetNextPrime(n *big.Int, verbose bool, debug bool) *big.Int {
	prime := false
	if debug {
		verbose = false
	}
	if n.Bit(0) == 0 {
		n.Add(n, one)
	}
	for !prime {
		n.Add(n, two)
		t0 := time.Now()
		prime = IsPrime(n)
		if verbose {
			fmt.Printf(".")
		}
		if debug {
			fmt.Printf("%t (%dms): %s\n", prime, time.Now().Sub(t0).Nanoseconds()/1000000, n)
		}
	}
	if verbose {
		fmt.Println("")
	}
	return n
}

func IsPrime(n *big.Int) bool {
	if n.Cmp(zero) == 0 {
		return false
	}
	if n.Cmp(one) == 0 {
		return false
	}
	if n.Bit(0) == 0 {
		return false
	}
	radixList := []int64{3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47, 53, 59, 61, 67, 71, 73, 79, 83, 89, 97} //, 2591, 7879, 9041, 10663, 11587, 12637};
	for _, rr := range radixList {
		nb := big.NewInt(0)
		nb.Add(zero, n)
		radix := big.NewInt(rr)
		if rr <= 7 {
			if nb.Mod(n, radix).Cmp(zero) == 0 {
				return false
			}
		}
		nb.Add(zero, n)
		if !isPrimeForRadix(nb, radix) {
			return false
		}
	}
	return true
}

func isPrimeForRadix(n *big.Int, radix *big.Int) bool {
	//Test Miller-Rabin
	if n.Cmp(radix) == 0 {
		return true
	}
	dd := big.NewInt(0)
	dd = dd.Sub(n, one)
	hh := 0
	for dd.Bit(0) == 0 {
		hh++
		dd.Div(dd, two)
	}
	nn1 := big.NewInt(0)
	nn1 = nn1.Sub(n, one)
	xx := big.NewInt(0)
	xx = PowModulo(radix, dd, n)
	if xx.Cmp(one) == 0 {
		return true
	}
	if xx.Cmp(nn1) == 0 {
		return true
	}
	ii := 1
	for ii != hh {
		xx = PowModulo(xx, two, n)
		if xx.Cmp(one) == 0 {
			return false
		}
		if xx.Cmp(nn1) == 0 {
			return true
		}
		ii++
	}
	return false
}

/*
//Test Miller-Rabin
func (n *BigInt) IsPrimeMillerRabin(radix *big.Int, dd *big.Int, hh int) bool {
        //compute: radix^dd % this
        xx := NewCopy(radix)
        nn1 := NewCopy(n)
        nn1.SubInt(1)
        two := NewInt(2)
        xx.PowModulo(dd, n)
        if xx.EqualsInt(1) {
                return true
        }
        if xx.Equals(nn1) {
                return true
        }
        ii := 1
        for ii != hh {
                two.SetInt(2)
                xx.PowModulo(two, n)
                if xx.EqualsInt(1) {
                        return false
                }
                if xx.Equals(nn1) {
                        return true
                }
                ii++
        }
        return false
}

func (n *BigInt) IsPrimeFermat(radix *big.Int) bool {
        //compute: radix^dd % this
        nn1 := NewCopy(n)
        nn1.SubInt(1)
        radix.PowModulo(nn1, n)
        if radix.EqualsInt(1) {
                return true
        }
        return false
}


// return x when x * b mod n = 1
func InverseModulo(b *big.Int, n *big.Int) *big.Int {
	b = b.Add(n)
	tmp := big.NewInt(0)
	tmp2 := big.NewInt(0)
	n0 := NewCopy(n)
	b0 := NewCopy(b)
	t0 := NewInt(0)
	t := NewInt(1)
	q := NewCopy(n0)
	q.Div(b0)
	tmp.Set(q)
	tmp.Mul(b0)
	r := NewCopy(n0)
	r.Sub(tmp)
	for r.GreaterInt(0) {
		tmp.Set(q)
		tmp.Mul(t)
		tmp2.Set(t0)
		tmp2.Sub(tmp)
		if tmp2.GreaterOrEqualInt(0) {
			tmp2.Modulo(n)
		} else {
			tmp2.Negation()
			tmp.Set(tmp2)
			tmp2.Set(n)
			tmp.Modulo(n)
			tmp2.Sub(tmp)
		}
		t0.Set(t)
		t.Set(tmp2)
		n0.Set(b0)
		b0.Set(r)
		q.Set(n0)
		q.Div(b0)
		tmp.Set(q)
		tmp.Mul(b0)
		r.Set(n0)
		r.Sub(tmp)
	}
	if !b0.EqualsInt(1) {
		return New()
	}
	return t
}
*/
