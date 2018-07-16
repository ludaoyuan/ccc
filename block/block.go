package block

import (
	"time"
	"bytes"
	"crypto/sha256"
	"strconv"
	"fmt"
)

type Block struct {
	Timestamp	int64
	Data		[]byte
	PrevBlockHash	[]byte
	Hash		[]byte
	Nonce		int
}

func (b *Block) SetHash() {
	timestamp := []byte(strconv.FormatInt(b.Timestamp, 10))
	headers := bytes.Join([][]byte{timestamp, b.Data, b.PrevBlockHash}, []byte{})
	hash := sha256.Sum256(headers)

	b.Hash = hash[:]
}

func NewBlock(Data []byte, PrevBlockHash[]byte, blockchain *[]*Block) *Block {
	block := &Block{time.Now().Unix(), Data, PrevBlockHash, []byte{}, 0}
//	block.SetHash()

	/* Mining */
	if blockchain != nil {
		proofOfWork := &ProofOfWork{(*blockchain)[len(*blockchain) - 1], nil}

		proofOfWork.SetTarget()

		fmt.Printf("%b\n", *proofOfWork.Target)
		fmt.Printf("%X\n", *proofOfWork.Target)
		fmt.Printf("%v\n", *proofOfWork.Target)

		proofOfWork.Mining()
	}
	/* ************************/

	return block
}

/* -----------------------------------------------------------------*/

func NewGenesisBlock() *Block {
	return NewBlock([]byte("Genesis Block"), []byte{}, nil)
}
