package common

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
)

const KeyLen = 64

type Key [64]byte

var ZeroKey Key

func GenerateKey() (ecdsa.PrivateKey, error) {
	curve := elliptic.P256()
	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		return ecdsa.PrivateKey{}, err
	}

	return *private, nil
}

func (k *Key) SetBytes(newKey []byte) error {
	nhlen := len(newKey)
	if nhlen != AddressLen {
		return fmt.Errorf("invalid hash length of %v, want %v", nhlen,
			KeyLen)
	}
	copy(k[:], newKey)

	return nil
}
