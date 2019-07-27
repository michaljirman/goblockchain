package wallet

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"

	"github.com/rs/zerolog/log"

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

func (w Wallet) Address() []byte {
	pubHash := PublicKeyHash(w.PublicKey)
	versionedHash := append([]byte{version}, pubHash...)
	checksum := Checksum(versionedHash)

	fullHash := append(versionedHash, checksum...)
	address := Base58Encode(fullHash)

	log.Debug().Msgf("pub key: %x", w.PublicKey)
	log.Debug().Msgf("pub hash: %x", pubHash)
	log.Debug().Msgf("address: %s", address)
	return address
}

// Validates address by comparing checksum within an address with expected / computed checksum for the payload.
//
// Address:    14LErwM2aHhdsDym6PkyutyG9ZSm51UHXc
// 			fmt.Printf("%x", wallet.Base58Decode([]byte("14LErwM2aHhdsDym6PkyutyG9ZSm51UHXc")))
// FullHash:   00248bd9e7a51b7dd07aba9766a7c62d50207902802bc6c767
// Version:    00
// PubKeyHash: 248bd9e7a51b7dd07aba9766a7c62d5020790280
// Checksum:   2bc6c767
func ValidateAddress(address string) bool {
	pubKeyHash := Base58Decode([]byte(address))
	actualChecksum := pubKeyHash[len(pubKeyHash)-checksumLength:]
	version := pubKeyHash[0]
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-checksumLength]
	targetChecksum := Checksum(append([]byte{version}, pubKeyHash...))
	return bytes.Compare(actualChecksum, targetChecksum) == 0
}

func NewKeyPair() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()

	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic().Err(err)
	}
	pub := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)
	return *private, pub
}

func MakeWallet() *Wallet {
	private, public := NewKeyPair()
	wallet := Wallet{private, public}
	return &wallet
}

func PublicKeyHash(pubKey []byte) []byte {
	pubHash := sha256.Sum256(pubKey)
	hasher := ripemd160.New()
	_, err := hasher.Write(pubHash[:])
	if err != nil {
		log.Panic().Err(err)
	}
	publicRipMD := hasher.Sum(nil)
	return publicRipMD
}

func Checksum(payload []byte) []byte {
	firstHash := sha256.Sum256(payload)
	secondHash := sha256.Sum256(firstHash[:])
	return secondHash[:checksumLength]
}
