package types

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
)

type Block struct {
	// 区块头
	header *BlockHeader
	// 交易
	Transactions Transactions
	// 区块大小
	// Size uint32
}

var GenesisBlock = &Block{
	header: genesisBlockHeader,
	// Transactions: make(Transactions, 1),
	Transactions: Transactions{genesisTransaction},
}

// TODO: 序列化，反序列化需要修改为小端序 binary包
func (b *Block) EncodeToBytes() ([]byte, error) {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)

	err := encoder.Encode(b)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func DecodeToBlock(blockBytes []byte) (*Block, error) {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(blockBytes))
	err := decoder.Decode(&block)
	if err != nil {
		return nil, err
	}
	return &block, nil
}

func (b *Block) ParentHash() [32]byte {
	return b.header.ParentHash
}

func (b *Block) Hash() [32]byte {
	return b.header.Hash
}

func (b *Block) Height() uint32 {
	return b.header.Height
}

// 创世区块
// func NewGenesisBlock(coinbase *types.Transaction) *Block {
//	// 创世区块,从1开始
//	return NewBlock([]*types.Transaction{coinbase}, []byte{}, 0)
// }

func NewBlock(txs []*Transaction, parentHash [32]byte, parentHeight uint32) (*Block, error) {
	header, err := NewBlockHeader(parentHeight, parentHash)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return &Block{
		header:       header,
		Transactions: txs,
	}, nil
}

func DoubleHash(stream []byte) [32]byte {
	hash := sha256.Sum256(stream[:])
	return sha256.Sum256(hash[:])
}

func (b *Block) HashTransactions() ([32]byte, error) {
	var txHashes [][]byte
	var txHash [32]byte

	for _, tx := range b.Transactions {
		hash, err := tx.Hash()
		if err != nil {
			log.Println(err.Error())
			return [32]byte{}, err
		}
		txHashes = append(txHashes, hash[:])
	}
	txHash = DoubleHash(bytes.Join(txHashes, []byte{}))

	return txHash, nil
}

func (b *Block) FindTransaction(txhash [32]byte) *Transaction {
	for _, tx := range b.Transactions {
		if bytes.Compare(tx.TxHash[:], txhash[:]) == 0 {
			return tx
		}
	}
	return nil
}
