package miner

import (
	"bytes"
	"common"
	"core"
	"core/types"
	"log"

	"github.com/syndtr/goleveldb/leveldb"
)

type Miner struct {
	version int64
	chain   *core.BlockChain
	chainDB *leveldb.DB

	worker   *Worker
	txPool   *TxPool
	coinbase common.Hash // PHK

	// Message
	txMsg         chan *types.Transaction
	localBlockMsg chan *types.Block

	stop chan struct{}
}

func NewMiner(chaindb *leveldb.DB, chain *core.BlockChain, coinbase common.Hash) *Miner {
	m := &Miner{
		chainDB:  chaindb,
		chain:    chain,
		coinbase: coinbase,
		txPool:   NewTxPool(),
		worker:   NewWorker(chain.LastBlockHeader()),
		stop:     make(chan struct{}),

		txMsg:         make(chan *types.Transaction),
		localBlockMsg: make(chan *types.Block),
	}
	return m
}

func (m *Miner) Start() {
	for {
		select {
		case <-m.stop:
			m.stopWork()
			return
		case tx, ok := <-m.txMsg:
			if ok {
				m.txPool.AddTx(tx)
			}
		case <-m.txPool.Notify():
			go m.Mining()
		}
	}
}

func (m *Miner) stopWork() {
	m.worker.Stop()
}

func (m *Miner) Stop() {
	close(m.stop)
	m.stop = make(chan struct{})
}

func (m *Miner) Close() {
}

func (m *Miner) ReceiveTx(newTx *types.Transaction) {
	if !m.chain.VerifyTransaction(newTx) {
		return
	}

	m.txMsg <- newTx
}

func (m *Miner) Mining() {
	txs := m.txPool.Txs()
	tx := m.CreateCoinbaseTx()
	txs = append(txs, tx)
	header := types.NewBlockHeader(m.version, common.ZeroHash, m.chain.LastBlock())

	m.worker.SetHeader(header)
	m.worker.Run()

	block := types.NewBlock(header, txs)
	m.localBlockMsg <- block
}

func (m *Miner) NotifyNewLocalBlock() <-chan *types.Block {
	return m.localBlockMsg
}

func (m *Miner) CreateCoinbaseTx() *types.Transaction {
	return nil
}

// func (m *Miner) MergeBlock(newBlocks []*types.Block) {
// 	for _, block := range newBlocks {
// 		if !m.chain.VerifyBlock(block) {
// 			return
// 		}
//
// 		if m.chain.Height() > block.Height() {
// 			return
// 		}
// 	}
// }

func (m *Miner) Update(blocks []*types.Block) error {
	m.stopWork()

	var buf bytes.Buffer
	for _, block := range blocks {
		err := block.Encode(&buf)
		if err != nil {
			log.Println(err.Error())
			return err
		}
		hash := block.Hash()
		err = m.chainDB.Put(hash[:], buf.Bytes(), nil)
		if err != nil {
			log.Println(err.Error())
			return err
		}
	}

	m.chain.SetLastBlock(blocks[len(blocks)-1])
	m.worker.SetHeader(m.chain.LastBlockHeader())

	m.updateCache()
	return nil
}

func (m *Miner) updateCache() {
	txs := m.txPool.Txs()
	for _, tx := range txs {
		if m.chain.TxExist(tx) {
			m.txPool.removeTx(tx)
		}
	}
}
