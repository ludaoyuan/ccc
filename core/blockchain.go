package core

import (
	"core/types"
	"crypto/ecdsa"
	"errors"
	"log"
	"sync"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/filter"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

const (
	chainPath           = "./data/chain"
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

func (bc *Blockchain) CreateGenesisBlock() (*types.Block, error) {
	genesisHash := types.GenesisBlock.Hash()

	genesisBlockStream, err := types.GenesisBlock.EncodeToBytes()
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	// log.Println(string(genesisHash[:]))
	// log.Println(string(genesisBlockStream))
	err = bc.chaindb.Put(genesisHash[:], genesisBlockStream, nil)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return types.GenesisBlock, nil
}

func NewBlockchain() *Blockchain {
	opts := opt.Options{
		ErrorIfExist: false,
		Strict:       opt.DefaultStrict,
		Compression:  opt.NoCompression,
		Filter:       filter.NewBloomFilter(10),
	}
	chaindb, err := leveldb.OpenFile(chainPath, &opts)
	if err != nil {
		log.Fatal(err.Error())
	}

	chain := &Blockchain{
		hbLock:          new(sync.RWMutex),
		heightHashCache: make(map[uint32][32]byte),
		chaindb:         chaindb,
		lastBlock:       types.GenesisBlock,
	}

	chain.hbLock.Lock()
	chain.heightHashCache[chain.lastBlock.Height()] = chain.lastBlock.Hash()
	chain.hbLock.Unlock()

	return chain
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
	blockStream, err := bc.chaindb.Get(blockhash, nil)
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
	iter := NewBlockchainIterator(bc.chaindb, bc.lastBlock.Hash())

	for iter.Next() {
		block := iter.Value()
		for _, tx := range block.Transactions {
			txID := string(tx.TxHash[:])

		Outputs:
			for outIdx, out := range tx.TxOut {
				if stxos[txID] != nil {
					for _, stxoindex := range stxos[txID] {
						if int(stxoindex) == outIdx {
							continue Outputs
						}
					}
				}

				// outs := UTXO[txID]
				var outs types.TxOuts
				outs.Outs = append(outs.Outs, out)
				UTXO[txID] = &outs
			}

			if tx.IsCoinbase() == false {
				for _, in := range tx.TxIn {
					inTxID := string(in.ParentTxHash[:])
					stxos[inTxID] = append(stxos[inTxID], in.ParentTxOutIndex)
				}
			}
		}

		if block.IsGenesisBlock() == true {
			break
		}
	}

	return UTXO, nil
}

func (bc *Blockchain) FindTransaction(txHash [32]byte) (*types.Transaction, error) {
	bci := NewBlockchainIterator(bc.chaindb, bc.lastBlock.Hash())

	for bci.Next() {
		block := bci.Value()

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
		parentTxs[string(parentTx.TxHash[:])] = parentTx
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
		ParentTxs[string(ParentTx.TxHash[:])] = ParentTx
	}

	return tx.Verify(ParentTxs)
}

func (bc *Blockchain) MineBlock(minerAddr []byte, txs types.Transactions, utxo *UTXOSet) error {
	// 更新状态信息
	for _, tx := range txs {
		if bc.VerifyTransaction(tx) != true {
			err := errors.New("ERROR: Invalid transaction")
			log.Println(err.Error())
			return err
		}
	}

	newBlock, err := types.NewBlock(txs, bc.lastBlock.Hash(), bc.lastBlock.Height()+1)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	// 更新数据库
	stream, err := newBlock.EncodeToBytes()
	if err != nil {
		log.Println(err.Error())
		return err
	}

	key := newBlock.Hash()
	err = bc.chaindb.Put(key[:], stream, nil)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	err = utxo.UpdateByBlock(newBlock)
	if err != nil {
		log.Println(err.Error())
		key := newBlock.Hash()
		bc.chaindb.Delete(key[:], nil)

		return err
	}

	err = bc.DumpDB(newBlock)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

func (bc *Blockchain) DumpDB(newBlock *types.Block) error {
	bc.hbLock.Lock()
	newhash := newBlock.Hash()
	bc.lastBlock = newBlock
	bc.heightHashCache[newBlock.Height()] = newhash
	bc.hbLock.Unlock()

	value, err := newBlock.EncodeToBytes()
	if err != nil {
		log.Println(err.Error())
		return err
	}

	err = bc.chaindb.Put(newhash[:], value, nil)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}
