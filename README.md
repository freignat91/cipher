# CIPHER

cipher v0.0.1

# Purpose

A project to learn basic of RSA encryption.

Generate big RSA keys (8192, 16384, 32768, ...) and encrypt file with them using RSA encryption directly and not using a faster (but less secured) symmetric key.

This project has only one external dependency, the cobra project to manage command line. It uses math/big library, but compute it-self the needed primes numbers and RSA key various numbers. It's possible to update the code to raise the prime probability to be true until the point you want. it'll be just slower. 



# Install

- prerequisite: have go installed and GOPATH set
- install glide: go get glide
- clone this project in ĜOPATH/src/github.com/freignat91/cipher
- execute: glide update
- execute: make install
- then the command cipher is available

For Ubuntu, you have a pre-build cipher.ubuntu file you can use without cloning and building the projet.

# Usage

## global options

- --help help on command
- -v verbose, display information during command execution
- -debug: display more information during command execution

## cipher createKeys [keyPath] -size [keysize]

This commande generates [keyPath].pub and [keyPath].key keys (public and private) having [keysize] bits long

## cipher encryptFile [sourceFilePath] [targetFilePath] [publicKeyPath]

This command encrypt the file [sourceFilePath] and save the  result in [targetFilePath] using the public key [publicKeyPath]

## cipher decryptFile [sourceFilePath] [targetFilePath] [privateKeyPath]

This command decrypt the file [sourceFilePath] and save the result in [targetFilePath] using the private key [privateKeyPath]


## speed

Using key size from 8192 to 32768 take time:

on a Latitude E6540 under ubuntu 16.10:

average key creation time:
- 2048 bits:   < 1s
- 4096 bits:  ~12 seconds - 20 seconds
- 8192 bits:  ~140 seconds - 300 seconds
- 16284 bits: ~500 seconds - 2 hours
- 32768 bits: ~9 hours - several days

it's possible to use intermediate size, all 64 bits multiple are accepted.

encryption time:
- 2038 bits:      ~60 ko/s
- 4096 bits:      ~19 ko/s
- 8192 bits:      ~6 Ko/s
- 16284 bits:     ~2 ko/s
- 32768 bits:     ~500 oct/s

decryption time:
- 2038 bits:      ~15 ko/s
- 4096 bits:      ~5 Ko/s
- 8192 bits:      ~2.5 Ko/s
- 16284 bits:     ~500 oct/s
- 32768 bits:     ~200 oct/s


ok, it's pretty slow, that's why RSA is more used to encrypt symetric key which is used to encrypt/decrypt file, but it's far more secured.

For security reason, don't share your public and private keys, private and public keys should stay in this context, (encrypt/decrypt your own files), secret.
The encrypt/decryption algorithm of this project don't use any padding scheme. It's not a security issue if the keys stay secret and especially are not used to authenticate, but only to encrypt/decrypt files.





