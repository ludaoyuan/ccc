package types

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"errors"
	"log"
	"math/big"
	"time"
)

// 创建交易
// 验证交易

const Subsidy = 25
const HashSize = 32
const genesisTimestamp = 1530603941

var ZeroHash = [HashSize]byte{}

const defaultTxInOutAlloc = 15

type Transactions []*Transaction

type Transaction struct {
	// 交易Hash
	TxHash [32]byte
	// 交易锁定时间
	LockTime uint32
	// 交易输入
	TxIn []*TxIn
	// 交易输出
	TxOut []*TxOut
}

var genesisTransaction = &Transaction{
	TxHash:   [32]byte{81, 122, 132, 80, 108, 208, 170, 237, 19, 15, 144, 202, 49, 32, 194, 73, 224, 56, 69, 168, 60, 125, 125, 192, 127, 85, 91, 95, 49, 134, 126, 159},
	LockTime: 1531072238,
	TxIn:     []*TxIn{genesisTxIn},
	TxOut:    []*TxOut{genesisTxOut},
}

// 币基交易没有输入
func CreateCoinbaseTX(to []byte) (*Transaction, error) {
	tx := NewTx()

	tx.AddTxIn(&TxIn{
		ParentTxOutIndex: -1,
	})
	tx.AddTxOut(&TxOut{
		Value:      Subsidy,
		PubKeyHash: to,
	})

	var err error
	tx.TxHash, err = tx.Hash()
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return tx, nil
}

// 判断是否为币基交易
func (tx *Transaction) IsCoinbase() bool {
	return len(tx.TxIn) == 1 && len(tx.TxIn[0].ParentTxHash) == 0 && tx.TxIn[0].ParentTxOutIndex == -1
}

func (tx Transaction) TrimmedCopy() *Transaction {
	var ins []*TxIn
	var outs []*TxOut

	for _, in := range tx.TxIn {
		ins = append(ins, &TxIn{in.ParentTxHash, in.ParentTxOutIndex, nil, nil})
	}

	for _, out := range tx.TxOut {
		outs = append(outs, &TxOut{out.Value, out.PubKeyHash})
	}

	txCopy := &Transaction{TxHash: tx.TxHash, LockTime: uint32(time.Now().Unix()), TxIn: ins, TxOut: outs}

	return txCopy
}

// 交易签名
func (tx *Transaction) Sign(sig ecdsa.PrivateKey, parentTxs map[string]*Transaction) error {
	if tx.IsCoinbase() {
		return nil
	}

	for _, in := range tx.TxIn {
		if parentTxs[string(in.ParentTxHash[:])].TxHash == ZeroHash {
			err := errors.New("ERROR: preious Transaction error")
			log.Println(err.Error())
			return err
		}
	}

	txCopy := tx.TrimmedCopy()

	var err error
	for inID, in := range txCopy.TxIn {
		preTx := parentTxs[string(in.ParentTxHash[:])]
		txCopy.TxIn[inID].SignatureKey = nil
		txCopy.TxIn[inID].PubKeyHash = preTx.TxOut[in.ParentTxOutIndex].PubKeyHash
		txCopy.TxHash, err = txCopy.Hash()
		if err != nil {
			log.Println(err.Error())
			return err
		}
		txCopy.TxIn[inID].PubKeyHash = nil

		r, s, err := ecdsa.Sign(rand.Reader, &sig, txCopy.TxHash[:])
		if err != nil {
			log.Println(err.Error())
		}
		signature := append(r.Bytes(), s.Bytes()...)

		copy(tx.TxIn[inID].SignatureKey[:], signature)
	}
	return nil
}

func (tx *Transaction) Verify(parentTxs map[string]*Transaction) bool {
	if tx.IsCoinbase() {
		return true
	}
	for _, in := range tx.TxIn {
		if parentTxs[string(in.ParentTxHash[:])].TxHash == ZeroHash {
			return false
		}
	}

	txCopy := tx.TrimmedCopy()
	curve := elliptic.P256()

	var err error
	for inID, in := range tx.TxIn {
		preTx := parentTxs[string(in.ParentTxHash[:])]
		txCopy.TxIn[inID].SignatureKey = nil
		txCopy.TxIn[inID].PubKeyHash = preTx.TxOut[in.ParentTxOutIndex].PubKeyHash
		txCopy.TxHash, err = txCopy.Hash()
		if err != nil {
			log.Println(err.Error())
			return false
		}
		txCopy.TxIn[inID].PubKeyHash = nil

		r := big.Int{}
		s := big.Int{}
		sigLen := len(in.SignatureKey)
		r.SetBytes(in.SignatureKey[:(sigLen / 2)])
		s.SetBytes(in.SignatureKey[(sigLen / 2):])

		x := big.Int{}
		y := big.Int{}
		keyLen := len(in.PubKeyHash)
		x.SetBytes(in.PubKeyHash[:(keyLen / 2)])
		y.SetBytes(in.PubKeyHash[(keyLen / 2):])

		rawPubKey := ecdsa.PublicKey{curve, &x, &y}
		if ecdsa.Verify(&rawPubKey, txCopy.TxHash[:], &r, &s) == false {
			return false
		}
	}

	return true

}

func (tx *Transaction) EncodeToBytes() ([]byte, error) {
	var buffer bytes.Buffer

	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(tx)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return buffer.Bytes(), nil
}

func DecodeToTransaction(txStream []byte) (*Transaction, error) {
	var tx Transaction

	decoder := gob.NewDecoder(bytes.NewReader(txStream))
	err := decoder.Decode(&tx)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return &tx, nil
}

func (tx *Transaction) Hash() ([32]byte, error) {
	var hash [32]byte

	txCopy := *tx
	txCopy.TxHash = [32]byte{}

	txStream, err := txCopy.EncodeToBytes()
	if err != nil {
		log.Println(err.Error())
		return ZeroHash, err
	}
	hash = sha256.Sum256(txStream)

	return hash, nil
}

func (tx *Transaction) AddTxIn(ti *TxIn) {
	tx.TxIn = append(tx.TxIn, ti)
}

func (tx *Transaction) AddTxOut(to *TxOut) {
	tx.TxOut = append(tx.TxOut, to)
}

func NewTx() *Transaction {
	return &Transaction{
		TxIn:     make([]*TxIn, 0, defaultTxInOutAlloc),
		TxOut:    make([]*TxOut, 0, defaultTxInOutAlloc),
		LockTime: uint32(time.Now().Unix()),
	}
}
