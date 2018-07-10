package types

type TxIn struct {
	// 引用交易的Hash
	ParentTxHash     [32]byte
	ParentTxOutIndex int64
	SignatureKey     []byte
	PubKeyHash       []byte
}

var genesisTxIn = &TxIn{
	ParentTxOutIndex: 0,
	// PubKeyHash:       []byte{238, 182, 209, 217, 143, 250, 69, 157, 152, 113, 95, 125, 203, 91, 25, 141, 20, 1, 37, 114, 61, 74, 208, 106, 167, 174, 95, 154, 148, 75, 39, 145, 173, 73, 247, 116, 169, 179, 164, 111, 66, 196, 188, 231, 46, 4, 182, 21, 151, 251, 184, 251, 227, 135, 225, 195, 45, 106, 162, 253, 6, 20, 116, 189},
}
