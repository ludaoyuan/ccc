package wallet

import (
	"bytes"
	"common"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/hex"
	"log"
)

const version = byte(0x00)
const addressChecksumLen = 4

type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey  common.Hash
}

func NewWallet() (*Wallet, error) {
	private, public, err := newKeyPair()
	if err != nil {
		return nil, err
	}

	wallet := Wallet{private, public}
	return &wallet, nil
}

func (w *Wallet) Address() string {
	pubKeyHash := w.PubKeyHash()

	versionedPayload := append([]byte{version}, pubKeyHash[:]...)
	checksum := common.Checksum(versionedPayload)

	fullPayload := append(versionedPayload, checksum...)
	address := common.Base58Encode(fullPayload)

	return hex.EncodeToString(address)
}

func (w *Wallet) PubKeyHash() common.Hash {
	pubKeyHash, _ := common.HashPubKey(w.PublicKey)
	var hash common.Hash
	hash.SetBytes(pubKeyHash)
	return hash
}

func newKeyPair() (ecdsa.PrivateKey, common.Hash, error) {
	curve := elliptic.P256()
	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Println(err.Error())
		return ecdsa.PrivateKey{}, common.ZeroHash, err
	}

	pubKey := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)

	var hash common.Hash
	err = hash.SetBytes(pubKey)
	if err != nil {
		log.Println(err.Error())
		return ecdsa.PrivateKey{}, common.ZeroHash, err
	}

	return *private, hash, nil
}

func ValidateAddress(address string) bool {
	pubKeyHash := common.Base58Decode([]byte(address))
	actualChecksum := pubKeyHash[len(pubKeyHash)-addressChecksumLen:]
	version := pubKeyHash[0]
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-addressChecksumLen]
	targetChecksum := common.Checksum(append([]byte{version}, pubKeyHash...))

	return bytes.Compare(actualChecksum, targetChecksum) == 0
}
