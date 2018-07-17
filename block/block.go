package block

import (
	"time"
	"bytes"
	"crypto/sha256"
	"strconv"
	"fmt"

	"util"
)

type Block struct {
	Timestamp	int64
	Data		[]byte
	PrevBlockHash	[]byte
	Hash		[]byte
	Nonce		uint64
}

func (b *Block) SetHash() {
	timestamp := []byte(strconv.FormatInt(b.Timestamp, 10))
	headers := bytes.Join([][]byte{timestamp, b.Data, b.PrevBlockHash, util.IntToHex(int64(b.Nonce))}, []byte{})
	hash := sha256.Sum256(headers)

	b.Hash = hash[:]
}

func (b *Block) Serialize() ([]byte, error) {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(b)
	if err != nil {
		return err
	}

	return result.Bytes()
}

func DeserializeBlock(d []byte) (*Block, error) {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(d))

	err := decoder.Decode(&block)
	if err != nil {
		return err
	}

	return &block
}

func NewBlock(Data []byte, PrevBlockHash[]byte, blockchain *[]*Block) *Block {
	block := &Block{time.Now().Unix(), Data, PrevBlockHash, []byte{}, 0}
//	block.SetHash()

	/* Mining */
	if blockchain != nil /* && len(*blockchain) > 0 */ {
		fmt.Println(len(*blockchain))
		proofOfWork := &ProofOfWork{(*blockchain)[len(*blockchain) - 1], nil}

		proofOfWork.SetTarget()

		fmt.Printf("%b\n", *proofOfWork.Target)
		fmt.Printf("%X\n", *proofOfWork.Target)
		fmt.Printf("%v\n", *proofOfWork.Target)

		var isHard bool
		block.Nonce, block.Hash, isHard = proofOfWork.Mining()

		if isHard {
			targetBits--
		} else {
			targetBits++
		}

	} else {
		block.SetHash()
	}
	/* ************************/

	return block
}

/* -----------------------------------------------------------------*/

func NewGenesisBlock() *Block {
	return NewBlock([]byte("Genesis Block"), []byte{}, nil)
}
