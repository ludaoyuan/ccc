package common

import (
	"crypto/sha256"
	"fmt"

	"golang.org/x/crypto/ripemd160"
)

const version = byte(0x00)
const addressChecksumLen = 4

type Address [AddressLen]byte

func (addr *Address) SetBytes(newHash []byte) error {
	nhlen := len(newHash)
	if nhlen != AddressLen {
		return fmt.Errorf("invalid hash length of %v, want %v", nhlen,
			AddressLen)
	}
	copy(addr[:], newHash)

	return nil
}

func Address2PubKeyHash(address string) Hash {
	pubKeyHash := Base58Decode([]byte(address))
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-addressChecksumLen]

	var hash Hash
	hash.SetBytes(pubKeyHash)

	return hash
}

func PubKey2Address(pubKey Hash) Address {
	pubKeyHash, _ := HashPubKey(pubKey)

	versionedPayload := append([]byte{version}, pubKeyHash...)
	checksum := Checksum(versionedPayload)

	fullPayload := append(versionedPayload, checksum...)
	address := Base58Encode(fullPayload)

	var addr Address
	addr.SetBytes(address)

	return addr
}

func Checksum(payload []byte) []byte {
	firstSHA := sha256.Sum256(payload)
	secondSHA := sha256.Sum256(firstSHA[:])

	return secondSHA[:addressChecksumLen]
}

// hash160: RIPEMD160(SHA256(PubKey))
func HashPubKey(pubKey Hash) ([]byte, error) {
	publicSHA256 := sha256.Sum256(pubKey[:])

	RIPEMD160Hasher := ripemd160.New()
	_, err := RIPEMD160Hasher.Write(publicSHA256[:])
	if err != nil {
		return nil, err
	}
	publicRIPEMD160 := RIPEMD160Hasher.Sum(nil)

	return publicRIPEMD160, nil
}
