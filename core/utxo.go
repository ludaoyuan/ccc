package core

import (
	"core/types"
	"log"

	"github.com/syndtr/goleveldb/leveldb"
)

const (
	BlockUTXOPath = "./data/utxo"
)

// 维护本地UTXO
type UTXOSet struct {
	chain *Blockchain
	// 需要维护两个UTXODB 一个确定的, 一个是pending的
	utxodb *leveldb.DB
}

func NewUTXOSet(c *Blockchain) (*UTXOSet, error) {
	db, err := leveldb.OpenFile(BlockUTXOPath, nil)

	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	utxoset := &UTXOSet{
		chain:  c,
		utxodb: db,
	}

	return utxoset, nil
}

// txid --> *Txouts
func (u UTXOSet) CreateUTXOSet() map[string]*types.TxOuts {
	utxos := make(map[string]*types.TxOuts)
	stxos := make(map[string][]int64)
	iter := u.chain.chaindb.NewIterator(nil, nil)

	// 遍历所有链记录, 提取未花费输出, 并返回
	for iter.Next() {
		block, err := u.chain.GetBlockByHash(iter.Key())
		if err != nil {
			log.Println(err.Error())
			continue
		}

		for _, tx := range block.Transactions {
			txID := string(tx.TxHash[:])

			for outIndex, _ := range tx.TxOut {
				if stxos[txID] != nil {
					for _, stxoIndex := range stxos[txID] {
						if int(stxoIndex) == outIndex {
							continue
						}
					}
				}

				outs := utxos[txID]
				outs.Outs = append(outs.Outs, tx.TxOut[outIndex])
				utxos[txID] = outs
			}

			if tx.IsCoinbase() == false {
				for _, in := range tx.TxIn {
					preOutTxID := string(in.ParentTxHash[:])
					stxos[preOutTxID] = append(stxos[preOutTxID], in.ParentTxOutIndex)
				}
			}
		}

		// 创世区块
		if len(block.ParentHash()) == 0 {
			break
		}
	}
	iter.Release()
	return utxos
}

// 凑钱
func (u UTXOSet) MakeARaise(pubKeyHash []byte, amount uint32) (uint32, map[string][]int) {
	raise := make(map[string][]int)
	var acc uint32

	iter := u.utxodb.NewIterator(nil, nil)
	for iter.Next() {
		txID := string(iter.Key())
		outs, err := types.DecodeToTxOuts(iter.Value())
		if err != nil {
			log.Println(err.Error())
			continue
		}

		for outIndex, out := range outs.Outs {
			if out.IsLockedWithKey(pubKeyHash) && acc < amount {
				acc += out.Value
				raise[txID] = append(raise[txID], outIndex)
			}
		}
	}
	iter.Release()

	return acc, raise
}

// 重建UTXO集合, TODO: 判断是否存在,存在先删除
func (u *UTXOSet) Reindex() {
	utxos := u.CreateUTXOSet()
	for txID, outs := range utxos {
		txidStream := []byte(txID)

		outsStream, err := outs.EncodeToBytes()
		if err != nil {
			log.Println(err.Error())
			continue
		}

		err = u.utxodb.Put(txidStream, outsStream, nil)
		if err != nil {
			log.Println(err.Error())
			continue
		}
	}
}

func (u UTXOSet) UpdateByBlock(block *types.Block) error {
	for _, tx := range block.Transactions {
		err := u.UpdateByTx(tx)
		if err != nil {
			log.Println(err.Error())
			continue
		}
	}
	return nil
}

func (u UTXOSet) ToDB(utxos map[string]*types.TxOuts) error {
	for txid, outs := range utxos {
		stream, err := outs.EncodeToBytes()
		if err != nil {
			log.Println(err.Error())
			return err
		}
		err = u.utxodb.Put([]byte(txid), stream, nil)
		if err != nil {
			log.Println(err.Error())
			return err
		}
	}
	return nil
}

func (u UTXOSet) TxToUTXODB(txID [32]byte, outs *types.TxOuts) error {
	outsBytes, err := outs.EncodeToBytes()
	if err != nil {
		log.Println(err.Error())
		return err
	}
	err = u.utxodb.Put(txID[:], outsBytes, nil)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func (u UTXOSet) UpdateByTx(tx *types.Transaction) error {
	if tx.IsCoinbase() {
		newOuts := types.TxOuts{}
		for i := range tx.TxOut {
			newOuts.Outs = append(newOuts.Outs, tx.TxOut[i])
		}

		err := u.TxToUTXODB(tx.TxHash, &newOuts)
		if err != nil {
			log.Println(err.Error())
			return err
		}
		return nil
	}

	for _, in := range tx.TxIn {
		updatedOuts := types.TxOuts{}
		stream, err := u.utxodb.Get(in.ParentTxHash[:], nil)
		if err != nil {
			log.Println(err.Error())
			return err
		}
		outs, err := types.DecodeToTxOuts(stream)
		if err != nil {
			log.Println(err.Error())
			return err
		}

		for outIndex, out := range outs.Outs {
			if outIndex != int(in.ParentTxOutIndex) {
				updatedOuts.Outs = append(updatedOuts.Outs, out)
			}
		}

		if len(updatedOuts.Outs) == 0 {
			err := u.utxodb.Delete(in.ParentTxHash[:], nil)
			if err != nil {
				log.Println(err.Error())
				return err
			} else {
				stream, err := updatedOuts.EncodeToBytes()
				if err != nil {
					log.Println(err.Error())
					return err
				}
				err = u.utxodb.Put(in.PubKeyHash[:], stream, nil)
				if err != nil {
					log.Println(err.Error())
					return err
				}
			}
		}
	}
	return nil
}

func (u UTXOSet) FindUTXOs(pubKeyHash []byte) (*types.TxOuts, error) {
	var UTXOs types.TxOuts

	iter := u.utxodb.NewIterator(nil, nil)

	for iter.Next() {
		outs, err := types.DecodeToTxOuts(iter.Value())
		if err != nil {
			log.Println(err.Error())
			iter.Release()
			return nil, err
		}

		for _, out := range outs.Outs {
			if out.IsLockedWithKey(pubKeyHash) {
				UTXOs.Outs = append(UTXOs.Outs, out)
			}
		}
	}

	iter.Release()
	return &UTXOs, nil
}

func (u UTXOSet) FindTxOutsOfAmount(pubkeyHash []byte, amount uint32) (uint32, map[string][]int) {
	unspentOutputs := make(map[string][]int)
	var accumulated uint32

	iter := u.utxodb.NewIterator(nil, nil)
	defer iter.Release()

	for iter.Next() {
		txID := string(iter.Key())
		outs, err := types.DecodeToTxOuts(iter.Value())
		if err != nil {
			log.Println(err.Error())
			continue
		}

		for outIdx, out := range outs.Outs {
			if out.IsLockedWithKey(pubkeyHash) && accumulated < uint32(amount) {
				accumulated += out.Value
				unspentOutputs[txID] = append(unspentOutputs[txID], outIdx)
			}
		}
	}

	return accumulated, unspentOutputs
}
