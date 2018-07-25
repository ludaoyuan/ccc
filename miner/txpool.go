package miner

import (
	"common"
	"core/types"
	"sync"
)

const defaultTxsLength = 10
const maxTxLimits = 10

type TxPool struct {
	poolMu  sync.Mutex
	pending map[common.Address]*types.Transaction
	next    types.Transactions

	full chan struct{}
}

func NewTxPool() *TxPool {
	return &TxPool{
		pending: make(map[common.Address]*types.Transaction),
		next:    make(types.Transactions, 0),
		full:    make(chan struct{}),
	}
}

func (pool *TxPool) Notify() <-chan struct{} {
	pool.full <- struct{}{}
	return pool.full
}

func (pool *TxPool) Len() int64 {
	return int64(len(pool.pending))
}

func (pool *TxPool) AddTx(tx *types.Transaction) {
	if pool.IsAreadyExits(tx) {
		return
	}

	addr := tx.FromAddr()
	_, ok := pool.pending[addr]
	if !ok {
		pool.addNext(addr, tx)
		return
	}
	pool.pending[addr] = tx
	if pool.Len() >= maxTxLimits {
		pool.Notify()
	}
}

func (pool *TxPool) IsAreadyExits(tx *types.Transaction) bool {
	return true
}

func (pool *TxPool) addNext(addr common.Address, tx *types.Transaction) {
	pool.next = append(pool.next, tx)
}

func (pool *TxPool) removeTx(tx *types.Transaction) {
	addr := tx.FromAddr()
	_, ok := pool.pending[addr]
	if ok {
		delete(pool.pending, addr)
	}

	hash := tx.Hash()
	for i, oldTx := range pool.next {
		if hash.IsEqual(oldTx.Hash()) {
			pool.next = append(pool.next[:i], pool.next[i+1:]...)
			break
		}
	}
}

func (pool *TxPool) Txs() types.Transactions {
	pool.poolMu.Lock()
	defer pool.poolMu.Unlock()

	txs := make(types.Transactions, 0, len(pool.pending)+len(pool.next))
	for addr, tx := range pool.pending {
		txs = append(txs, tx)
		delete(pool.pending, addr)
	}

	for _, tx := range pool.next {
		txs = append(txs, tx)
	}
	return txs
}
