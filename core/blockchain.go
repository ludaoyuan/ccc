package core

import (
	"core/types"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"log"
	"sync"

	"github.com/syndtr/goleveldb/leveldb"
)

const (
	chainPath           = "./data/chaindb"
	genesisCoinbaseData = "Inc Block Chain Start at 2018/7/3"
)

// db 中需要保留最新一个区块hash
// db 中需要保留创世区块hash
// 全局一个 "l" --> *Block
type Blockchain struct {
	hbLock    *sync.RWMutex
	lastBlock *types.Block
	// 缓存的区块高度与hash映射关系
	heightHashCache map[uint32][32]byte
	chaindb         *leveldb.DB
}

func CreateGenesisBlock() (*types.Block, error) {
	chain := NewBlockchain()

	genesisHash := types.GenesisBlock.Hash()

	genesisBlockStream, err := types.GenesisBlock.EncodeToBytes()
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	err = chain.chaindb.Put(genesisHash[:], genesisBlockStream, nil)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return types.GenesisBlock, nil
}

func NewBlockchain() *Blockchain {
	chaindb, err := leveldb.OpenFile(chainPath, nil)
	if err != nil {
		log.Fatal(err.Error())
	}

	chain := &Blockchain{
		hbLock:          new(sync.RWMutex),
		heightHashCache: make(map[uint32][32]byte),
		chaindb:         chaindb,
		lastBlock:       types.GenesisBlock,
	}

	chain.Init()

	return chain
}

func (bc *Blockchain) Iterator() *BlockchainIterator {
	bci := &BlockchainIterator{bc.lastBlock.Hash(), bc.chaindb, nil, nil}

	return bci
}

func (bc *Blockchain) Init() {
	bc.hbLock.Lock()
	bc.heightHashCache[bc.lastBlock.Height()] = bc.lastBlock.Hash()
	bc.hbLock.Unlock()
}

// 获取高度
func (bc *Blockchain) GetBlockCount() uint32 {
	bc.hbLock.RLock()
	defer bc.hbLock.RUnlock()

	if bc.lastBlock != nil {
		return bc.lastBlock.Height()
	}
	return 0
}

func (bc *Blockchain) GetBlockByNumber(height uint32) (*types.Block, error) {
	bc.hbLock.RLock()
	blockhash := bc.heightHashCache[height]
	bc.hbLock.RUnlock()

	return bc.GetBlockByHash(blockhash[:])
}

func (bc *Blockchain) GetBlockByHash(blockhash []byte) (*types.Block, error) {
	blockStream, err := bc.chaindb.Get(blockhash[:], nil)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	block, err := types.DecodeToBlock(blockStream)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return block, nil
}

func (bc *Blockchain) LastBlock() *types.Block {
	bc.hbLock.RLock()
	defer bc.hbLock.RUnlock()

	return bc.lastBlock
}

func (bc *Blockchain) UpdateState(block *types.Block) {
	bc.hbLock.Lock()
	defer bc.hbLock.Unlock()

	bc.lastBlock = block
	bc.heightHashCache[block.Height()] = block.Hash()
}

// TODO: 验证过后存储数据库 并且更新相关状态信息
func (bc *Blockchain) AddBlock(tx *types.Transaction) error {
	return nil
}

// txid --> []*TxOut
func (bc *Blockchain) InitUTXOSet() (map[string]*types.TxOuts, error) {
	UTXO := make(map[string]*types.TxOuts)
	stxos := make(map[string][]int64)
	iter := bc.chaindb.NewIterator(nil, nil)
	// bci := bc.Iterator()

	for iter.Next() {
		block, err := types.DecodeToBlock(iter.Value())
		if err != nil {
			log.Println(err.Error())
			return nil, err
		}
		for _, tx := range block.Transactions {
			txID := hex.EncodeToString(tx.TxHash[:])

		Outputs:
			for outIdx, out := range tx.TxOut {
				if stxos[txID] != nil {
					for _, stxoindex := range stxos[txID] {
						if int(stxoindex) == outIdx {
							continue Outputs
						}
					}
				}

				outs := UTXO[txID]
				outs.Outs = append(outs.Outs, out)
				UTXO[txID] = outs
			}

			if tx.IsCoinbase() == false {
				for _, in := range tx.TxIn {
					inTxID := hex.EncodeToString(in.ParentTxHash[:])
					stxos[inTxID] = append(stxos[inTxID], in.ParentTxOutIndex)
				}
			}
		}

		if len(block.ParentHash()) == 0 {
			break
		}
	}

	return UTXO, nil
}

func (bc *Blockchain) FindTransaction(txHash [32]byte) (*types.Transaction, error) {
	bci := bc.Iterator()

	for bci.Next() {
		block := bci.Value

		tx := block.FindTransaction(txHash)
		if tx != nil {
			return tx, nil
		}

		if len(block.ParentHash()) == 0 {
			break
		}
	}

	return &types.Transaction{}, errors.New("Transaction is not found")
}

func (bc *Blockchain) SignTransaction(tx *types.Transaction, privKey ecdsa.PrivateKey) error {
	parentTxs := make(map[string]*types.Transaction)

	for _, in := range tx.TxIn {
		parentTx, err := bc.FindTransaction(in.ParentTxHash)
		if err != nil {
			log.Println(err.Error())
			return err
		}
		parentTxs[hex.EncodeToString(parentTx.TxHash[:])] = parentTx
	}

	tx.Sign(privKey, parentTxs)
	return nil
}

func (bc *Blockchain) VerifyTransaction(tx *types.Transaction) bool {
	if tx.IsCoinbase() {
		return true
	}

	ParentTxs := make(map[string]*types.Transaction)

	for _, in := range tx.TxIn {
		ParentTx, err := bc.FindTransaction(in.ParentTxHash)
		if err != nil {
			log.Println(err.Error())
			return false
		}
		ParentTxs[hex.EncodeToString(ParentTx.TxHash[:])] = ParentTx
	}

	return tx.Verify(ParentTxs)
}

func (bc *Blockchain) MineBlock(minerAddr []byte, txs types.Transactions, utxo *UTXOSet) (*types.Block, error) {
	// 更新状态信息
	bc.hbLock.Lock()

	for _, tx := range txs {
		if bc.VerifyTransaction(tx) != true {
			err := errors.New("ERROR: Invalid transaction")
			log.Println(err.Error())
			bc.hbLock.Unlock()
			return nil, err
		}
	}

	newBlock, err := types.NewBlock(txs, bc.lastBlock.Hash(), bc.lastBlock.Height()+1)
	if err != nil {
		log.Println(err.Error())
		bc.hbLock.Unlock()
		return nil, err
	}

	// 更新数据库
	stream, err := newBlock.EncodeToBytes()
	if err != nil {
		log.Println(err.Error())
		bc.hbLock.Unlock()
		return nil, err
	}

	key := newBlock.Hash()
	err = bc.chaindb.Put(key[:], stream, nil)
	if err != nil {
		log.Println(err.Error())
		bc.hbLock.Unlock()
		return nil, err
	}

	err = utxo.UpdateByBlock(newBlock)
	if err != nil {
		log.Println(err.Error())
		key := newBlock.Hash()
		bc.chaindb.Delete(key[:], nil)

		bc.hbLock.Unlock()
		return nil, err
	}

	bc.lastBlock = newBlock
	bc.heightHashCache[newBlock.Height()] = newBlock.Hash()
	bc.hbLock.Unlock()

	return newBlock, nil
}
