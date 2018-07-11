package types

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
)

type Block struct {
	// 区块头
	Header *BlockHeader
	// 交易
	Transactions Transactions
	// 区块大小
	// Size uint32
}

var GenesisBlock = &Block{
	Header:       genesisBlockHeader,
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

func (b *Block) IsGenesisBlock() bool {
	return b.Header.IsGenesisBlock()
}

func (b *Block) ParentHash() [32]byte {
	return b.Header.ParentHash
}

func (b *Block) GenerateHash() ([32]byte, error) {
	return b.Header.GenerateHash()
}

func (b *Block) SetHash(hash [32]byte) {
	b.Header.Hash = hash
}

func (b *Block) Hash() [32]byte {
	return b.Header.Hash
}

func (b *Block) Height() uint32 {
	return b.Header.Height
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
		Header:       header,
		Transactions: txs,
	}, nil
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
	txHash = sha256.Sum256(bytes.Join(txHashes, []byte{}))

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
