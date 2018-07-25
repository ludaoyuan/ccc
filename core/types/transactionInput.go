package types

import (
	"common"
	"encoding/hex"
)

const MaxPrevOutIndex uint32 = 0xffffffff
const MaxTxInSequenceNum uint32 = 0xffffffff

type OutPoints []*OutPoint

type OutPoint struct {
	TxHash common.Hash
	Index  uint32
}

func GenesisOutPoint() *OutPoint {
	return &OutPoint{
		TxHash: common.ZeroHash,
		Index:  MaxPrevOutIndex,
	}
}

type TxIn struct {
	PreviousOutPoint OutPoint
	Signature        []byte
	PublicKey        common.Hash
}

func GenesisTxIn() *TxIn {
	return &TxIn{
		PreviousOutPoint: *GenesisOutPoint(),
		Signature:        nil,
		PublicKey:        common.ZeroHash,
	}
}

func NewTxIn(parentTxHash common.Hash, parentTxOutIndex uint32, pubKey common.Hash) *TxIn {
	return &TxIn{
		PreviousOutPoint: OutPoint{parentTxHash, parentTxOutIndex},
		PublicKey:        pubKey,
	}
}

func (in *TxIn) HexParentTxHash() string {
	return hex.EncodeToString(in.PreviousOutPoint.TxHash[:])
}

func (in *TxIn) ParentTxHash() common.Hash {
	return in.PreviousOutPoint.TxHash
}
