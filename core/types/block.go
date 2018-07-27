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
	Header       *BlockHeader
	Transactions Transactions
}

func GenesisBlock() *Block {
	return &Block{
		Header:       GenesisBlockHeader(),
		Transactions: GenesisTransactions(),
	}
}

func NewBlock(header *BlockHeader, txs Transactions) *Block {
	return &Block{
		Header:       header,
		Transactions: txs,
	}
}

func (b *Block) Height() uint32 {
	return b.Header.Height
}

func (b *Block) Bits() uint32 {
	return b.Header.Bits
}

func (b *Block) Hash() common.Hash {
	return b.Header.BlockHash()
}

func (b *Block) ParentBlockHash() common.Hash {
	return b.Header.ParentHash
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
	return b.Header.IsGenesisBlock()
}

func (b *Block) Dump(db *leveldb.DB) error {
	hash := b.Header.BlockHash()
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
