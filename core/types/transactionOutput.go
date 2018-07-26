package types

import (
	"bytes"
	"common"
)

type TxOuts []*TxOut

type TxOut struct {
	Value      int64
	PubKeyHash common.Hash
}

func GenesisTxOut() *TxOut {
	phk := common.Address2PubKeyHash(common.GenesisHexAddress)
	return &TxOut{
		Value:      common.Subsidy,
		PubKeyHash: phk,
	}
}

func NewTxOut(value int64, address string) *TxOut {
	out := &TxOut{Value: value}
	out.Lock(address)
	return out
}

func (out *TxOut) Lock(address string) {
	out.PubKeyHash = common.Address2PubKeyHash(address)
}

func (out *TxOut) MatchPubKeyHash(pubKeyHash common.Hash) bool {
	return bytes.Compare(out.PubKeyHash[:], pubKeyHash[:]) == 0
}
