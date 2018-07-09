package types

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
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
	Height: 0,
	// TODO:
	// Hash:       common.NewHashFromStr2("00000000000000000000000000000000"),
	// ParentHash: types.ZeroHash,
	Timestamp: uint32(time.Now().Unix()),
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
