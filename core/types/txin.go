package types

type TxIn struct {
	// 引用交易的Hash
	ParentTxHash     [32]byte
	ParentTxOutIndex int64
	SignatureKey     [32]byte
	PubKeyHash       [32]byte
}

var genesisTxIn = &TxIn{
	ParentTxHash:     ZeroHash,
	ParentTxOutIndex: 0,
	SignatureKey:     ZeroHash,
	PubKeyHash:       ZeroHash,
}
