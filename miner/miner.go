package miner

import "github.com/syndtr/goleveldb/leveldb"

type Miner struct {
	blockchain *types.Blockchain
	chainDB    *leveldb.DB
}
