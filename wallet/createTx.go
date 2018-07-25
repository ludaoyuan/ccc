package wallet

import (
	"core/types"
	"crypto/ecdsa"
	"crypto/rand"
	"errors"
	"fmt"
	"log"
)

func (ws *WalletSvr) UpdateByTx(tx *types.Transaction) {
	newOuts := make(types.TxOuts, 0, len(tx.TxOut))
	for _, out := range tx.TxOut {
		newOuts = append(newOuts, out)
	}

	ws.utxos[tx.HexHash()] = newOuts

	if !tx.IsCoinbase() {
		for _, in := range tx.TxIn {
			inParentHash := in.HexParentTxHash()
			outs := ws.utxos[inParentHash]
			updateOuts := make(types.TxOuts, 0, len(outs))
			for idx, out := range outs {
				if int64(idx) != in.PreviousOutPoint.Index {
					updateOuts = append(updateOuts, out)
				}
			}

			if len(updateOuts) != 0 {
				delete(ws.utxos, inParentHash)
			} else {
				ws.utxos[inParentHash] = updateOuts
			}
		}
	}
}

func (ws *WalletSvr) UpdateByBlock(block *types.Block) {
	for _, tx := range block.Transactions() {
		ws.UpdateByTx(tx)
	}
}

func (ws *WalletSvr) CreateTx(from, to string, amount int64) (*types.Transaction, error) {
	acc, outPoints := ws.FindSpendableOuts(from, amount)
	if acc < amount {
		err := errors.New("ERROR: Insufficient balance")
		log.Println(err.Error())
		return nil, err
	}

	ins := make([]*types.TxIn, 0, len(outPoints))
	for _, point := range outPoints {
		in := types.NewTxIn(point.TxHash, point.Index, ws.myWallet.PublicKey)
		ins = append(ins, in)
	}

	outs := make(types.TxOuts, 0)
	outs = append(outs, types.NewTxOut(amount, to))
	if acc > amount {
		outs = append(outs, types.NewTxOut(acc-amount, from))
	}

	tx := types.NewTransction(0, ins, outs)

	ws.SignTx(tx)

	return tx, nil
}

func (ws *WalletSvr) SignTx(tx *types.Transaction) error {
	if tx.IsCoinbase() {
		return nil
	}

	parentTxs, err := ws.chain.FindParentTransactions(tx)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	for _, in := range tx.TxIn {
		parentTx, ok := parentTxs[in.ParentHashString()]
		if !ok || (ok && parentTx.TxHash == ZeroHash) {
			err := errors.New("ERROR: preious Transaction error")
			log.Println(err.Error())
			return err
		}
	}

	txCopy := tx.Copy()

	for i, in := range txCopy.TxIn {
		parentTx := parentTxs[in.ParentHashString()]
		txCopy.TxIn[i].SignatureKey = nil
		txCopy.TxIn[i].PubKeyHash = parentTx.TxOut[in.ParentTxOutIndex].PubKeyHash

		signData := fmt.Sprintf("%x\n", txCopy)

		r, s, err := ecdsa.Sign(rand.Reader, &sig, []byte(signData))
		if err != nil {
			log.Println(err.Error())
			return err
		}
		signature := append(r.Bytes(), s.Bytes()...)

		tx.TxIn[i].SignatureKey = signature
		txCopy.TxIn[i].PubKeyHash = nil
	}
	return nil
}
