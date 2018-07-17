package types

import (
	"bytes"
	"common"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"errors"
	"fmt"
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

func Genesishash() [32]byte {
	hash, _ := genesisTransaction.Hash()
	return hash
}

// 币基交易没有输入
func CreateCoinbaseTx(toPubkey []byte) (*Transaction, error) {
	tx := NewTx()

	tx.AddTxIn(&TxIn{
		ParentTxOutIndex: 0,
	})
	// TODO: 应该是pubkey的
	tx.AddTxOut(&TxOut{
		Value:      Subsidy,
		PubKeyHash: toPubkey,
	})

	var err error
	tx.TxHash, err = tx.Hash()
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return tx, nil
}

func (tx *Transaction) TxHashString() string {
	return hex.EncodeToString(tx.TxHash[:])
}

// 交易发起方的公钥
func (tx *Transaction) FromAddr() (common.Address, error) {
	return common.PubKeyToAddress(tx.TxIn[0].PubKeyHash)
}

func (tx *Transaction) Nil() bool {
	return tx.TxHash == ZeroHash
}

// 判断是否为币基交易
func (tx *Transaction) IsCoinbase() bool {
	return len(tx.TxIn) == 1 && len(tx.TxIn[0].ParentTxHash) == 0 && tx.TxIn[0].ParentTxOutIndex == -1
}

func (tx Transaction) TrimmedCopy() *Transaction {
	ins := make([]*TxIn, 0, len(tx.TxIn))
	outs := make([]*TxOut, 0, len(tx.TxOut))

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
		parentTx, ok := parentTxs[in.ParentHashString()]
		// 如果父交易不存在或者父交易 hash错误返回
		if !ok || (ok && parentTx.TxHash == ZeroHash) {
			err := errors.New("ERROR: preious Transaction error")
			log.Println(err.Error())
			return err
		}
	}

	txCopy := tx.TrimmedCopy()

	for inID, in := range txCopy.TxIn {
		parentTx := parentTxs[in.ParentHashString()]
		txCopy.TxIn[inID].SignatureKey = nil
		txCopy.TxIn[inID].PubKeyHash = parentTx.TxOut[in.ParentTxOutIndex].PubKeyHash

		signData := fmt.Sprintf("%x\n", txCopy)
		// signData, err := txCopy.EncodeToBytes()
		// if err != nil {
		// 	log.Println(err.Error())
		// 	return err
		// }

		// 此处签名是否需要直接对二进制进行签名
		r, s, err := ecdsa.Sign(rand.Reader, &sig, []byte(signData))
		if err != nil {
			log.Println(err.Error())
			return err
		}
		signature := append(r.Bytes(), s.Bytes()...)

		tx.TxIn[inID].SignatureKey = signature
		txCopy.TxIn[inID].PubKeyHash = nil
	}
	return nil
}

func (tx *Transaction) Verify(parentTxs map[string]*Transaction) bool {
	if tx.IsCoinbase() {
		return true
	}

	txCopy := tx.TrimmedCopy()
	curve := elliptic.P256()

	for inID, in := range tx.TxIn {
		parentTx := parentTxs[in.ParentHashString()]
		txCopy.TxIn[inID].SignatureKey = nil
		txCopy.TxIn[inID].PubKeyHash = parentTx.TxOut[in.ParentTxOutIndex].PubKeyHash

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

		verifyData := fmt.Sprintf("%x\n", txCopy)
		// verifyData, err := txCopy.EncodeToBytes()
		// if err != nil {
		// 	log.Println(err.Error())
		// 	return false
		// }

		rawPubKey := ecdsa.PublicKey{Curve: curve, X: &x, Y: &y}
		if ecdsa.Verify(&rawPubKey, []byte(verifyData), &r, &s) == false {
			log.Println("ecdsa Verify Error")
			return false
		}
		txCopy.TxIn[inID].PubKeyHash = nil
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
