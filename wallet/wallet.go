package wallet

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"log"

	"golang.org/x/crypto/ripemd160"
)

const (
	checksumLength = 4
	version        = byte(0x00)
)

type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}

func GetChecksumLength() int {
	return checksumLength
}

func (w *Wallet) Address() []byte {
	publicHash := PublicKeyHash(w.PublicKey)

	versionedHash := append([]byte{version}, publicHash...)
	checksum := CheckSum(versionedHash)

	fullHash := append(versionedHash, checksum...)
	address := Base58Encode(fullHash)

	fmt.Printf("pub key: %x\n", w.PublicKey)
	fmt.Printf("pub hash: %x\n", publicHash)
	fmt.Printf("address: %s\n", address)

	return address
}

func NewKeyPair() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()

	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic(err)
	}

	pub := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)
	return *private, pub
}

func MakeWallet() *Wallet {
	private, public := NewKeyPair()
	wallet := Wallet{PrivateKey: private, PublicKey: public}

	return &wallet
}

func ValidateAddress(address string) bool {
	publicKeyFullHash := Base58Decode([]byte(address))
	actualChecksum := publicKeyFullHash[len(publicKeyFullHash)-checksumLength:]
	version := publicKeyFullHash[0]
	publicKeyHash := publicKeyFullHash[1 : len(publicKeyFullHash)-checksumLength]

	targetChecksum := CheckSum(append([]byte{version}, publicKeyHash...))

	return bytes.Compare(actualChecksum, targetChecksum) == 0
}

func PublicKeyHash(publicKey []byte) []byte {
	publicHash := sha256.Sum256(publicKey)

	hasher := ripemd160.New()
	_, err := hasher.Write(publicHash[:])
	if err != nil {
		log.Panic(err)
	}

	publicRipMD := hasher.Sum(nil)

	return publicRipMD
}

func CheckSum(payload []byte) []byte {
	firstHash := sha256.Sum256(payload)
	secondHash := sha256.Sum256(firstHash[:])

	return secondHash[:checksumLength]
}
