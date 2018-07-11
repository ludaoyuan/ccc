package types

import (
	"bytes"
	"encoding/gob"
	"log"
)

type TxOut struct {
	// 单位聪
	Value      uint32
	PubKeyHash []byte
}

// TODO: 第一个地址
var genesisTxOut = &TxOut{
	Value:      Subsidy,
	PubKeyHash: []byte{105, 197, 202, 132, 29, 223, 157, 165, 31, 193, 245, 40, 157, 167, 253, 183, 171, 119, 47, 183},
}

func (out *TxOut) IsLockedWithKey(pubKeyHash []byte) bool {
	return bytes.Compare(out.PubKeyHash[:], pubKeyHash[:]) == 0
}

func (out *TxOut) Lock(address []byte) {
	key := Base58Decode(address)
	key := key[1 : len(key)-4]
	out.PubKeyHash = key
}

func NewTxOut(value uint32, address []byte) *TxOut {
	out := &TxOut{value, address}
	out.Lock(address)
	return out
}

type TxOuts struct {
	Outs []*TxOut
}

func DecodeToTxOuts(stream []byte) (*TxOuts, error) {
	var outs TxOuts

	decoder := gob.NewDecoder(bytes.NewReader(stream))
	err := decoder.Decode(&outs)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return &outs, nil

}

func (outs *TxOuts) EncodeToBytes() ([]byte, error) {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)

	err := encoder.Encode(outs)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}
