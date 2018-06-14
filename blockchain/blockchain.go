package blockchain

var genesisBlock = &Block{
	Index:         0,
	Timestamp:     1528884611,
	PrevBlockHash: []byte("QmYdbMbpwWUa5rD5S3Uekuk7qTWjWUwRugG66iRBDbqi7w"),
	Transactions:  []byte("Genesis Block"),
}

// Blockchain keeps a sequence of Blocks
type Blockchain struct {
	Blocks    []*Block
	lastIndex int64
}

// AddBlock saves provided data as a block in the blockchain
func (bc *Blockchain) AddBlock(data string) {
	prevBlock := bc.Blocks[len(bc.Blocks)-1]
	newBlock := NewBlock(data, prevBlock.Hash)
	newBlock.Index = prevBlock.Index + 1
	bc.lastIndex = newBlock.Index
	bc.Blocks = append(bc.Blocks, newBlock)
}

func (bc *Blockchain) GetBlock(index int64) string {
	if index >= int64(0) && index < int64(len(bc.Blocks)) {
		return string(bc.Blocks[index].Transactions)
	}
	return "NIL"
}

func (bc *Blockchain) Height() int64 {
	return bc.lastIndex
}

// NewBlockchain creates a new Blockchain with genesis Block
func NewBlockchain() *Blockchain {
	return &Blockchain{[]*Block{genesisBlock}, 0}
}
