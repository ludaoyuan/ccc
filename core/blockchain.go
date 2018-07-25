package core

import (
	"bytes"
	"common"
	"core/types"
	"errors"
	"log"

	"github.com/syndtr/goleveldb/leveldb"
)

type BlockChain struct {
	lastBlock *types.Block
	chainDB   *leveldb.DB
}

func NewBlockChain(chaindb *leveldb.DB) (*BlockChain, error) {
	bc := &BlockChain{
		lastBlock: types.GenesisBlock(),
		chainDB:   chaindb,
	}

	err := bc.lastBlock.Dump(chaindb)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return bc, nil
}

func (bc *BlockChain) GenesisBlock() *types.Block {
	return nil
}

func (bc *BlockChain) Height() uint32 {
	if bc.lastBlock != nil {
		return bc.lastBlock.Height()
	}

	return 0
}

func (bc *BlockChain) GetBlock(hash common.Hash) (*types.Block, error) {
	value, err := bc.chainDB.Get(hash[:], nil)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	var block types.Block
	buf := bytes.NewBuffer(value)
	err = block.Encode(buf)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return &block, nil
}

func (bc *BlockChain) SetLastBlock(block *types.Block) {
	bc.lastBlock = block
}

func (bc *BlockChain) LastBlock() *types.Block {
	return bc.lastBlock
}

func (bc *BlockChain) LastBlockHash() common.Hash {
	return bc.lastBlock.Hash()
}

func (bc *BlockChain) LastBlockHeader() *types.BlockHeader {
	return bc.lastBlock.Header()
}

func (bc *BlockChain) FindParentTransactions(tx *types.Transaction) (map[common.Hash]*types.Transaction, error) {
	parentTxs := make(map[common.Hash]*types.Transaction)

	for _, in := range tx.TxIn {
		parentTx, err := bc.FindTransaction(in.ParentTxHash())
		if err != nil {
			log.Println(err.Error())
			return nil, err
		}
		parentTxs[parentTx.Hash()] = parentTx
	}

	return parentTxs, nil
}

func (bc *BlockChain) FindTransaction(txHash common.Hash) (*types.Transaction, error) {
	iter := NewBlockChainIterator(bc.chainDB, bc.lastBlock.Hash())

	for iter.Next() {
		block := iter.Value()

		for _, tx := range block.Transactions() {
			if bytes.Compare(txHash[:], tx.TxHash[:]) == 0 {
				return tx, nil
			}
		}
	}

	return &types.Transaction{}, errors.New("Transaction is not found")
}

func (bc *BlockChain) TxExist(tx *types.Transaction) bool {
	iter := NewBlockChainIterator(bc.chainDB, bc.lastBlock.Hash())

	for iter.Next() {
		block := iter.Value()

		for _, oldTx := range block.Transactions() {
			if oldTx.Hash() == tx.Hash() {
				return true
			}
		}
	}
	return false
}

func (bc *BlockChain) VerifyBlock(block *types.Block) bool {
	return true
}

func (bc *BlockChain) VerifyTransaction(tx *types.Transaction) bool {
	if tx.IsCoinbase() {
		return true
	}

	parentTxs := make(map[common.Hash]*types.Transaction)

	for _, in := range tx.TxIn {
		parentTx, err := bc.FindTransaction(in.ParentTxHash())
		if err != nil {
			log.Println(err.Error())
			return false
		}

		parentTxs[parentTx.Hash()] = parentTx
	}

	return tx.Verify(parentTxs)
}

func (bc *BlockChain) ChainList() ([]common.Hash, error) {
	iter := NewBlockChainIterator(bc.chainDB, bc.LastBlockHash())

	hashList := make([]common.Hash, 0)

	for iter.Next() {
		hashList = append(hashList, iter.Key())
	}

	err := iter.Error()
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return hashList, nil
}

func (bc *BlockChain) HashExist(hash common.Hash) bool {
	_, err := bc.chainDB.Get(hash[:], nil)
	return err == nil
}

func (bc *BlockChain) FindCommonLastCommonBlock(hashList []common.Hash) ([]*types.Block, error) {
	iter := NewBlockChainIterator(bc.chainDB, bc.LastBlockHash())

	start, mid, end := 0, 0, len(hashList)-1
	for start <= end {
		mid = (start + end) / 2
		if start == end && start == mid+1 {
			break
		}

		isExist := bc.HashExist(hashList[mid])
		if isExist {
			start = mid + 1
			continue
		}

		end = mid - 1
	}

	err := iter.Error()
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return bc.GetMissingBlocks(hashList[end]), nil
}

func (bc *BlockChain) GetMissingBlocks(left common.Hash) []*types.Block {
	blocks := make([]*types.Block, 0)
	iter := NewBlockChainIterator(bc.chainDB, bc.LastBlockHash())

	for iter.Next() {
		if left.IsEqual(iter.Key()) {
			break
		}
		blocks = append(blocks, iter.Value())
	}

	return blocks
}
