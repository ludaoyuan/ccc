package types

import (
	"bytes"
	"encoding/gob"
	"log"
)

type TxOut struct {
	// 单位聪
	Value      uint32
	PubKeyHash [32]byte
}

// TODO: 第一个地址
var genesisTxOut = &TxOut{
	Value: Subsidy,
	// PubKeyHash: [32]byte(""),
}

func (out *TxOut) IsLockedWithKey(pubKeyHash []byte) bool {
	return bytes.Compare(out.PubKeyHash[:], pubKeyHash[:]) == 0
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
