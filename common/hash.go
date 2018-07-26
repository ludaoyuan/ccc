package common

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"log"
)

const AddressLen = 20
const HashLen = 32

var (
	ZeroHash Hash
)

type Hash [HashLen]byte

func (h *Hash) Bytes() []byte {
	return h[:]
}

func (h *Hash) IsNil() bool {
	return *h == ZeroHash
}

func (h *Hash) ToHex() string {
	return hex.EncodeToString(h[:])
}
func HexHash2Hash(s string) Hash {
	b, _ := hex.DecodeString(s)
	var hash Hash
	hash.SetBytes(b)

	return hash
}

func (h *Hash) SetBytes(newHash []byte) error {
	nhlen := len(newHash)
	if nhlen != HashLen {
		return fmt.Errorf("invalid hash length of %v, want %v", nhlen,
			HashLen)
	}
	copy(h[:], newHash)

	return nil
}

func (h *Hash) FromHex(str string) error {
	b, _ := hex.DecodeString(str)
	return h.SetBytes(b)
}

func (hash *Hash) IsEqual(target Hash) bool {
	return *hash == target
}

func DoubleHashH(b []byte) Hash {
	first := sha256.Sum256(b)
	return Hash(sha256.Sum256(first[:]))
}

func IntToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}
