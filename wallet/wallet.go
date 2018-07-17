package wallet

import (
	"common"
	"core"
	"core/types"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"log"
	"time"
)

// Wallet stores private and public keys
type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}

// NewWallet creates and returns a Wallet
func NewWallet() (*Wallet, error) {
	private, public, err := common.NewKeyPair()
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	wallet := Wallet{private, public}

	return &wallet, nil
}

func (w Wallet) PubKeyHash() ([]byte, error) {
	return common.HashPubKey(w.PublicKey)
}

// GetAddress returns wallet address
func (w Wallet) GetAddress() ([]byte, error) {
	pubKeyHash, err := common.HashPubKey(w.PublicKey)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	versionedPayload := append([]byte{common.Version}, pubKeyHash...)
	checksum := common.Checksum(versionedPayload)

	fullPayload := append(versionedPayload, checksum...)
	address := common.Base58Encode(fullPayload)

	return address, err
}

func (w Wallet) CreateTx(chain *core.Blockchain, to []byte, amount uint32, utxo *core.UTXOSet) (*types.Transaction, error) {
	var inputs []*types.TxIn
	var outputs []*types.TxOut

	pubKeyHash, err := w.PubKeyHash()
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	acc, validOutputs := utxo.FindTxOutsOfAmount(pubKeyHash, amount)

	if acc < amount {
		err := errors.New("ERROR: Insufficient balance")
		log.Println(err.Error())
		return nil, err
	}

	for txhashstr, outs := range validOutputs {
		txhash, _ := hex.DecodeString(txhashstr)
		for _, out := range outs {
			input := &types.TxIn{common.ToHash32(txhash), int64(out), nil, w.PublicKey}
			inputs = append(inputs, input)
		}
	}

	myAddr, err := w.GetAddress()
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	outputs = append(outputs, types.NewTxOut(amount, to))
	if acc > amount {
		outputs = append(outputs, types.NewTxOut(acc-amount, myAddr))
	}

	tx := &types.Transaction{LockTime: uint32(time.Now().Unix()), TxIn: inputs, TxOut: outputs}
	tx.TxHash, err = tx.Hash()
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	chain.SignTransaction(tx, w.PrivateKey)

	return tx, nil
}
