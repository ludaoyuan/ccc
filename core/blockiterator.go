package core

import (
	"core/types"
	"errors"

	"github.com/syndtr/goleveldb/leveldb"
)

var ZeroHash = [32]byte{}

type BlockchainIterator struct {
	pivot [32]byte
	db    *leveldb.DB
	value *types.Block
	err   error
}

func (bci *BlockchainIterator) Value() *types.Block {
	return bci.value
}

func (bci *BlockchainIterator) Next() bool {
	if bci.pivot == ZeroHash {
		bci.err = errors.New("Hash NIL")
		return false
	}

	blockStream, err := bci.db.Get(bci.pivot[:], nil)
	if err != nil {
		bci.err = err
		return false
	}

	bci.value, err = types.DecodeToBlock(blockStream)
	if err != nil {
		bci.err = err
		return false
	}

	bci.pivot = bci.value.ParentHash()

	return true
}

func (bci *BlockchainIterator) Error() error {
	return bci.err
}

func NewBlockchainIterator(db *leveldb.DB, blockHash [32]byte) *BlockchainIterator {
	return &BlockchainIterator{db: db, pivot: blockHash}
}
