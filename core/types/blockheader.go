package types

import (
	"bytes"
	"common"
	"encoding/gob"
	"io"
	"time"
)

const MaxBlockHeaderPayload = 16 + (common.HashLen * 2)

type BlockHeader struct {
	Version    int64
	ParentHash common.Hash
	MerkleRoot common.Hash
	Timestamp  int64
	Height     uint32
	Bits       uint32
	Nonce      uint32
}

func NewBlockHeader(version int64, root common.Hash, lastBlock *Block) *BlockHeader {
	return &BlockHeader{
		Version:    version,
		ParentHash: lastBlock.Hash(),
		MerkleRoot: root,
		Timestamp:  time.Now().Unix(),
		Height:     lastBlock.Height(),
		Bits:       lastBlock.Bits(),
	}
}

func GenesisBlockHeader() *BlockHeader {
	header := &BlockHeader{
		Version:    0,
		ParentHash: common.ZeroHash,
		MerkleRoot: common.ZeroHash,
		Timestamp:  time.Now().Unix(),
		Height:     1,
		Bits:       0x1d00ffff,
		Nonce:      0,
	}
	return header
}

func GenesisBlockHash() common.Hash {
	genesisHeader := GenesisBlockHeader()
	return genesisHeader.BlockHash()
}

// Decode the BlockHeader to Reader
// 创建解码器并接收一些值。
// Create a decoder and receive a value.
func (bh *BlockHeader) Decode(r io.Reader) error {
	// gob.NewDecoder(r)
	// decode(&v)
	gob.Register(BlockHeader{})
	dec := gob.NewDecoder(r)
	return dec.Decode(bh)
}

//  Encode the BlockHeader to writer
// Create an encoder and send a value.
// 创建编码器并发送一些值。
func (bh *BlockHeader) Encode(w io.Writer) error {
	// gob.NewEncoder()
	// encode(v)
	gob.Register(BlockHeader{})
	enc := gob.NewEncoder(w)
	return enc.Encode(*bh)
}

func (bh *BlockHeader) BlockHash() common.Hash {
	buf := bytes.NewBuffer(make([]byte, 0, MaxBlockHeaderPayload))
	_ = bh.Encode(buf)

	return common.DoubleHashH(buf.Bytes())
}

func (bh *BlockHeader) IsGenesisBlock() bool {
	return bh.Height == 1 && bh.ParentHash == common.ZeroHash
}
