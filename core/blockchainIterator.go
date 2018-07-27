package core

import (
	"bytes"
	"common"
	"core/types"

	"github.com/syndtr/goleveldb/leveldb"
)

type BlockChainIterator struct {
	currentHash common.Hash
	key         common.Hash
	chainDB     *leveldb.DB
	value       *types.Block
	err         error
}

func NewBlockChainIterator(db *leveldb.DB, blockHash common.Hash) *BlockChainIterator {
	return &BlockChainIterator{
		currentHash: blockHash,
		chainDB:     db,
	}
}

func (bci *BlockChainIterator) Value() *types.Block {
	return bci.value
}

func (bci *BlockChainIterator) Key() common.Hash {
	return bci.key
}

func (bci *BlockChainIterator) Error() error {
	return bci.err
}

func (bci *BlockChainIterator) Next() bool {
	if bci.currentHash.IsNil() {
		return false
	}

	// TODO:
	bci.key.SetBytes(bci.currentHash[:])

	b, err := bci.chainDB.Get(bci.currentHash.Bytes(), nil)
	if err != nil {
		bci.err = err
		return false
	}

	err = bci.value.Decode(bytes.NewBuffer(b))
	if err != nil {
		bci.err = err
		return false
	}

	bci.currentHash = bci.value.Hash()

	return true
}
