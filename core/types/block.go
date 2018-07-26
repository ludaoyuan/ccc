package types

import (
	"bytes"
	"common"
	"encoding/gob"
	"io"
	"log"

	"github.com/syndtr/goleveldb/leveldb"
)

type Block struct {
	header       *BlockHeader
	transactions Transactions
}

func GenesisBlock() *Block {
	return &Block{
		header:       GenesisBlockHeader(),
		transactions: GenesisTransactions(),
	}
}

func NewBlock(header *BlockHeader, txs Transactions) *Block {
	return &Block{
		header:       header,
		transactions: txs,
	}
}

func (b *Block) Height() uint32 {
	return b.header.Height
}

func (b *Block) Bits() uint32 {
	return b.header.Bits
}

func (b *Block) Hash() common.Hash {
	return b.header.BlockHash()
}

func (b *Block) ParentBlockHash() common.Hash {
	return b.header.ParentHash
}

func (b *Block) Header() *BlockHeader {
	return b.header
}

func (b *Block) Encode(w io.Writer) error {
	gob.Register(Block{})
	enc := gob.NewEncoder(w)
	return enc.Encode(*b)
}

func (b *Block) Decode(r io.Reader) error {
	gob.Register(Block{})
	dec := gob.NewDecoder(r)
	return dec.Decode(b)
}

func (b *Block) IsGenesisBlock() bool {
	return b.header.IsGenesisBlock()
}

func (b *Block) Transactions() Transactions {
	return b.transactions
}

func (b *Block) Dump(db *leveldb.DB) error {
	hash := b.header.BlockHash()
	var buf bytes.Buffer
	err := b.Encode(&buf)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	err = db.Put(hash.Bytes(), buf.Bytes(), nil)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}
