package wallet

import (
	"bytes"
	"common"
	"crypto/ecdsa"
	"encoding/hex"
)

const version = byte(0x00)
const addressChecksumLen = 4

type Wallet struct {
	PrivateKey ecdsa.PrivateKey
}

func NewWallet() (*Wallet, error) {
	private, err := common.GenerateKey()
	if err != nil {
		return nil, err
	}

	wallet := Wallet{private}
	return &wallet, nil
}

func (w *Wallet) PublicKey() common.Key {
	public := w.PrivateKey.PublicKey
	keyBytes := append(public.X.Bytes(), public.Y.Bytes()...)

	var key common.Key
	key.SetBytes(keyBytes)
	return key
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
	public := w.PublicKey()
	pubKeyHash, _ := common.HashPubKey(public)
	var hash common.Hash
	hash.SetBytes(pubKeyHash)
	return hash
}

func ValidateAddress(address string) bool {
	pubKeyHash := common.Base58Decode([]byte(address))
	actualChecksum := pubKeyHash[len(pubKeyHash)-addressChecksumLen:]
	version := pubKeyHash[0]
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-addressChecksumLen]
	targetChecksum := common.Checksum(append([]byte{version}, pubKeyHash...))

	return bytes.Compare(actualChecksum, targetChecksum) == 0
}
