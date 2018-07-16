package util

import (
	"bytes"
	"encoding/binary"
	"log"
)

// IntToHex converts an int64 to a byte array
func IntToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Fatalln(err.Error())
	}

	return buff.Bytes()
}

// HexToUint converts a byte array to an uint64
func HexToUint(bs []byte) uint64 {
	buff := bytes.NewBuffer(bs)
	var i uint64
	binary.Read(buff, binary.BigEndian, &i)
	return i
}
