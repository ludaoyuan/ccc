package types

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"log"
	"time"
)

var blockHeaderLen = 80

type BlockHeader struct {
	// 区块高度
	Height uint32
	// 当前区块头Hash
	Hash [32]byte
	// 前一个区块头hash
	ParentHash [32]byte
	// 用于简单支付验证
	// MerkleRoot []byte
	// 打包时间
	Timestamp uint32
}

var genesisBlockHeader = &BlockHeader{
	Height:    1,
	Hash:      [32]byte{222, 108, 139, 74, 52, 198, 125, 44, 26, 194, 244, 163, 88, 238, 111, 198, 98, 89, 11, 42, 108, 62, 243, 159, 155, 200, 177, 226, 112, 6, 105, 210},
	Timestamp: 1531072238,
}

func (bh *BlockHeader) String() string {
	return fmt.Sprintf("Height: %d, Hash: %s, ParentHash: %s, Timestamp: %d\n", bh.Height, bh.Hash, bh.ParentHash, bh.Timestamp)
}

func (bh *BlockHeader) IsGenesisBlock() bool {
	return bh.Height == 1 && bh.ParentHash == ZeroHash
}

// 所有的这些类似函数,有没有必要做成统一接口(反射太慢)
func (bh *BlockHeader) EncodeToBytes() ([]byte, error) {
	var buffer bytes.Buffer

	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(bh)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (bh *BlockHeader) GenerateHash() ([32]byte, error) {
	bhCopy := *bh
	bhCopy.Hash = [32]byte{}
	stream, err := bh.EncodeToBytes()
	if err != nil {
		log.Println(err.Error())
		return [32]byte{}, err
	}
	return sha256.Sum256(stream), nil
}

func DecodeToBlockHeader(stream []byte) (*BlockHeader, error) {
	var bh BlockHeader

	decoder := gob.NewDecoder(bytes.NewReader(stream))
	err := decoder.Decode(&bh)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return &bh, nil
}

func NewBlockHeader(parentHeight uint32, parentHash [32]byte) (*BlockHeader, error) {
	header := &BlockHeader{
		Height:     parentHeight + 1,
		ParentHash: parentHash,
		Timestamp:  uint32(time.Now().Unix()),
	}

	headerBytes, err := header.EncodeToBytes()
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	header.Hash = sha256.Sum256(headerBytes)

	return header, nil
}
