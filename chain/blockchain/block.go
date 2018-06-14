package blockchain

import (
	"bytes"
	"crypto/sha256"
	"strconv"
	"time"
)

// Block keeps block headers
type Block struct {
	Timestamp     int64
	Transactions  []byte
	PrevBlockHash []byte
	Hash          []byte
	Index         int64
}

func (b *Block) SetHash() {
	timestamp := []byte(strconv.FormatInt(b.Timestamp, 10))
	headers := bytes.Join([][]byte{b.PrevBlockHash, b.Transactions, timestamp}, []byte{})
	hash := sha256.Sum256(headers)

	b.Hash = hash[:]
}

// NewBlock creates and returns Block
func NewBlock(data string, prevBlockHash []byte) *Block {
	block := &Block{time.Now().Unix(), []byte(data), prevBlockHash, []byte{}, 0}
	block.SetHash()

	return block
}
