package backend

import (
	"strconv"
	"fmt"

	"block"
)

//var blockchain []*block.Block

func AddBlock(blockchain *[]*block.Block, data []byte) {
	prevBlock := (*blockchain)[len(*blockchain) - 1]
	block := block.NewBlock(data, prevBlock.Hash, blockchain)
	// 以下改成入库
	*blockchain = append(*blockchain, block)
}

func NewBlockchain() []*block.Block {
	var tip []byte
	blockchain := []*block.Block{block.NewGenesisBlock()}

	return blockchain
}

func Run() {
	blockchain := NewBlockchain()

	for i := 0; i < 64; i++ {
		AddBlock(&blockchain, []byte("Send 1 BTC to Ivan" + strconv.Itoa(i)))
	}
//	AddBlock(&blockchain, []byte("Send 2 more BTC to Ivan"))

	for k, v := range blockchain {
		fmt.Printf("INDEX:\t\t%d\nHASH:\t\t%X\nTIMESTAMP:\t%d\nDATA:\t\t%s\nPREVHASH:\t%X\n\n", k, v.Hash, v.Timestamp, string(v.Data), v.PrevBlockHash)
	}
}
