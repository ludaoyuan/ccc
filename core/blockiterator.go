package core

import (
	"core/types"

	"github.com/syndtr/goleveldb/leveldb"
)

var ZeroHash = [32]byte{}

type BlockchainIterator struct {
	pivotHash [32]byte
	db        *leveldb.DB
	value     *types.Block
	err       error
}

func (bci *BlockchainIterator) Value() *types.Block {
	return bci.value
}

func (bci *BlockchainIterator) Next() bool {
	// 结束迭代
	if bci.pivotHash == ZeroHash {
		return false
	}

	blockStream, err := bci.db.Get(bci.pivotHash[:], nil)
	if err != nil {
		bci.err = err
		return false
	}

	bci.value, err = types.DecodeToBlock(blockStream)
	if err != nil {
		bci.err = err
		return false
	}

	bci.pivotHash = bci.value.ParentHash()

	return true
}

func (bci *BlockchainIterator) Error() error {
	return bci.err
}

func NewBlockchainIterator(db *leveldb.DB, blockHash [32]byte) *BlockchainIterator {
	return &BlockchainIterator{db: db, pivotHash: blockHash}
}
