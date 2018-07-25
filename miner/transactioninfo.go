package miner

import (
	"common"
	"core/types"
)

type TransactionInfos []*TransactionInfo

type TransactionInfo struct {
	From   common.Address
	tx     *types.Transaction
	Height int64
	Added  int64
}
