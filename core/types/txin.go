package types

import "encoding/hex"

type TxIn struct {
	// 引用交易的Hash
	ParentTxHash     [32]byte
	ParentTxOutIndex int64
	SignatureKey     []byte
	PubKeyHash       []byte
}

var genesisTxIn = &TxIn{
	ParentTxOutIndex: 0,
}

func (in *TxIn) ParentHashString() string {
	// "0000000000000000000000000000000000000000000000000000000000000000"
	return hex.EncodeToString(in.ParentTxHash[:])
}
