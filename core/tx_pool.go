package core

import (
	"core/types"
	"errors"
	"log"
	"sync"

	"common"
)

const CachePoolSize = 500

// 接收交易验证后加入交易池
// 需要排序根据时间
type TxPool struct {
	mu          sync.RWMutex
	utxo        *UTXOSet
	minerPubKey []byte
	blockchain  *Blockchain
	// 待打包的交易:根据优先级 交易费 时间等因素优先打包的交易
	pending map[common.Address]types.Transaction
	// 优先级别低一些的交易
	queue map[common.Address]types.Transactions
}

// TODO: 需要新加入交易排序, 根据交易时间 交易权重 交易手续费排序
func NewTxPool(minerPubKey []byte, chain *Blockchain, utxo *UTXOSet) *TxPool {
	pool := &TxPool{
		blockchain:  chain,
		minerPubKey: minerPubKey,
		pending:     make(map[common.Address]types.Transaction),
		queue:       make(map[common.Address]types.Transactions),
	}

	return pool
}

func (pool *TxPool) validateTx(tx *types.Transaction) bool {
	return pool.blockchain.VerifyTransaction(tx)
}

// 添加交易
func (pool *TxPool) AddTx(tx *types.Transaction) error {
	if pool.validateTx(tx) {
		err := errors.New("Validate Transaction Failed")
		log.Println(err.Error())
		return err
	}

	from, err := tx.FromAddr()
	if err != nil {
		log.Println(err.Error())
		return err
	}
	_, ok := pool.pending[from]
	if ok == true {
		if _, ok = pool.queue[from]; !ok {
			pool.queue[from] = make(types.Transactions, 0)
		}
		pool.queue[from] = append(pool.queue[from], tx)
	} else {
		pool.pending[from] = *tx
	}

	if len(pool.pending) >= CachePoolSize {
		err := pool.Mine()
		if err != nil {
			log.Println(err.Error())
			return err
		}
	}
	return nil
}

func (pool *TxPool) Mine() error {
	txs := make(types.Transactions, 0, len(pool.pending))
	for addr, tx := range pool.pending {
		txs = append(txs, &tx)
		delete(pool.pending, addr)
	}

	// TODO: Pubkey
	return pool.blockchain.MineBlock(pool.minerPubKey, txs, pool.utxo)
}
