package core

import (
	"core/types"

	"github.com/syndtr/goleveldb/leveldb"
)

type BlockchainIterator struct {
	pivot [32]byte
	db    *leveldb.DB
	Value *types.Block
	err   error
}

func (i *BlockchainIterator) Next() bool {
	var block *types.Block

	blockStream, err := i.db.Get(i.pivot[:], nil)
	if err != nil {
		i.err = err
		return false
	}

	i.Value, err = types.DecodeToBlock(blockStream)
	if err != nil {
		i.err = err
		return false
	}

	i.pivot = block.ParentHash()
	return true
}

func (i *BlockchainIterator) Error() error {
	return i.err
}
