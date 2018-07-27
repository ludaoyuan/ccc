package types

import (
	"bytes"
	"common"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"math/big"
	"time"
)

type Transactions []*Transaction

type Transaction struct {
	Version  int32
	TxHash   common.Hash
	TxIn     []*TxIn
	TxOut    []*TxOut
	LockTime int64
}

func GenesisTransactions() Transactions {
	tx := &Transaction{
		Version:  0,
		TxIn:     make([]*TxIn, 0, 1),
		TxOut:    make([]*TxOut, 0, 1),
		LockTime: time.Now().Unix(),
	}

	tx.TxIn = append(tx.TxIn, GenesisTxIn())
	tx.TxOut = append(tx.TxOut, GenesisTxOut())

	tx.TxHash = tx.Hash()
	return Transactions{tx}
}

func NewTransction(version int32, in []*TxIn, out []*TxOut) *Transaction {
	tx := &Transaction{
		Version:  version,
		TxIn:     in,
		TxOut:    out,
		LockTime: time.Now().Unix(),
	}
	tx.TxHash = tx.Hash()
	return tx
}

func (tx *Transaction) Hash() common.Hash {
	var hash common.Hash
	txCopy := *tx
	txCopy.TxHash = common.ZeroHash

	var buf bytes.Buffer
	_ = txCopy.Encode(&buf)

	hash = sha256.Sum256(buf.Bytes())
	return common.Hash(hash)
}

func (tx *Transaction) Copy() *Transaction {
	newTx := Transaction{
		Version:  tx.Version,
		TxHash:   tx.TxHash,
		TxIn:     make([]*TxIn, 0, len(tx.TxIn)),
		TxOut:    make([]*TxOut, 0, len(tx.TxOut)),
		LockTime: tx.LockTime,
	}

	for _, oldTxIn := range tx.TxIn {
		// newTx.TxIn = append(newTx.TxIn, &TxIn{ParentTxHash: oldTxIn.ParentTxHash, ParentTxOutIndex: oldTxIn.ParentTxOutIndex, PublicKey: oldTxIn.PublicKey, Signature: oldTxIn.Signature})
		newTx.TxIn = append(newTx.TxIn, &TxIn{oldTxIn.PreviousOutPoint, oldTxIn.Signature, oldTxIn.PublicKey})
	}

	for _, oldTxOut := range tx.TxOut {
		newTx.TxOut = append(newTx.TxOut, &TxOut{PubKeyHash: oldTxOut.PubKeyHash, Value: oldTxOut.Value})
	}

	return &newTx
}

func (tx *Transaction) Encode(w io.Writer) error {
	gob.Register(Transaction{})
	enc := gob.NewEncoder(w)
	return enc.Encode(*tx)
}

func (tx *Transaction) Decode(r io.Reader) error {
	gob.Register(Transaction{})
	dec := gob.NewDecoder(r)
	return dec.Decode(tx)
}

func (tx *Transaction) HexHash() string {
	return hex.EncodeToString(tx.TxHash[:])
}

func (tx *Transaction) IsCoinbase() bool {
	in := tx.TxIn[0]
	return len(tx.TxIn) == 1 && in.PreviousOutPoint.TxHash.IsNil()
}

func (tx *Transaction) Verify(parentTxs map[common.Hash]*Transaction) bool {
	txCopy := tx.Copy()
	curve := elliptic.P256()

	for i, in := range tx.TxIn {
		parentTx := parentTxs[in.ParentTxHash()]
		txCopy.TxIn[i].Signature = nil
		// txCopy.TxIn[i].PublicKey = parentTx.TxOut[in.PreviousOutPoint.Index].PubKeyHash
		copy(txCopy.TxIn[i].PublicKey[:], parentTx.TxOut[in.PreviousOutPoint.Index].PubKeyHash[:])

		r := big.Int{}
		s := big.Int{}
		sigLen := len(in.Signature)
		r.SetBytes(in.Signature[:(sigLen / 2)])
		s.SetBytes(in.Signature[(sigLen / 2):])

		x := big.Int{}
		y := big.Int{}
		keyLen := len(in.PublicKey)
		x.SetBytes(in.PublicKey[:(keyLen / 2)])
		y.SetBytes(in.PublicKey[(keyLen / 2):])

		verifyData := fmt.Sprintf("%x\n", txCopy)

		rawPubKey := ecdsa.PublicKey{Curve: curve, X: &x, Y: &y}
		if ecdsa.Verify(&rawPubKey, []byte(verifyData), &r, &s) == false {
			log.Println("ecdsa Verify Error")
			return false
		}
		txCopy.TxIn[i].PublicKey = common.ZeroKey
	}

	return true
}

func (tx *Transaction) FromAddr() common.Address {
	pubKey := tx.TxIn[0].PublicKey
	return common.PubKey2Address(pubKey)
}
