package wallet

import (
	"common"
	"core"
	"core/types"
	"log"
)

func (ws *WalletSvr) UnspentOutputs(pubKeyHash common.Hash) error {
	spentTXOs := make(map[string][]int64)
	iter := core.NewBlockChainIterator(ws.chainDB, ws.chain.LastBlock().Hash())

	for iter.Next() {
		block := iter.Value()
		for _, tx := range block.Transactions {
			txHash := tx.HexHash()

		Outputs:
			for outIdx, out := range tx.TxOut {
				if spentTXOs[txHash] != nil {
					for _, spentOutIdx := range spentTXOs[txHash] {
						if int(spentOutIdx) == outIdx {
							continue Outputs
						}
					}
				}

				var outs types.TxOuts
				outs = append(outs, out)
				ws.utxos[txHash] = outs
			}

			if tx.IsCoinbase() == false {
				for _, in := range tx.TxIn {
					inParentHash := in.HexParentTxHash()
					spentTXOs[inParentHash] = append(spentTXOs[inParentHash], int64(in.PreviousOutPoint.Index))
				}
			}
		}

		if block.IsGenesisBlock() {
			break
		}
	}

	err := iter.Error()
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

func (ws *WalletSvr) GetBalance(address string) int64 {
	var amount int64
	hashPubKey := common.Address2PubKeyHash(address)
	for _, outs := range ws.utxos {
		for _, out := range outs {
			if out.MatchPubKeyHash(hashPubKey) {
				amount += int64(out.Value)
			}
		}
	}
	return amount
}

func (ws *WalletSvr) ListUTXOs(address string) types.TxOuts {
	newOuts := make(types.TxOuts, 0)
	hashPubKey := common.Address2PubKeyHash(address)
	for _, outs := range ws.utxos {
		for _, out := range outs {
			if out.MatchPubKeyHash(hashPubKey) {
				newOuts = append(newOuts, out)
			}
		}
	}
	return newOuts
}

func (ws *WalletSvr) FindSpendableOuts(address string, amount int64) (int64, types.OutPoints) {
	var total int64
	outPoints := make(types.OutPoints, 0)
	hashPubKey := common.Address2PubKeyHash(address)
	for hexTxHash, outs := range ws.utxos {
		hash := common.HexHash2Hash(hexTxHash)
		for index, out := range outs {
			if out.MatchPubKeyHash(hashPubKey) && total < amount {
				outPoints = append(outPoints, &types.OutPoint{hash, uint32(index)})
			}
		}
	}

	return total, outPoints
}
