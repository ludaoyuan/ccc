package miner

import (
	"bytes"
	"common"
	"core/types"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
)

var (
	maxNonce = int64(math.MaxInt64)
)

const targetBits = 5

// 挖矿难度调节
const DiffcultThreshold = 180

type Worker struct {
	header *types.BlockHeader
	target *big.Int

	stop chan struct{}
}

func NewWorker(header *types.BlockHeader) *Worker {
	w := &Worker{
		header: header,
		target: big.NewInt(1),
		stop:   make(chan struct{}),
	}

	w.target.Lsh(w.target, uint(256-targetBits))
	return w
}

func (w *Worker) SetHeader(header *types.BlockHeader) {
	w.header = header
}

func (w *Worker) Stop() {
}

func (w *Worker) tryData(nonce int64) []byte {
	data := bytes.Join(
		[][]byte{
			common.IntToHex(int64(w.header.Version)),
			w.header.ParentHash[:],
			w.header.MerkleRoot[:],
			common.IntToHex(int64(w.header.Timestamp)),
			common.IntToHex(int64(nonce)),
			// TODO:
			common.IntToHex(int64(targetBits)),
			common.IntToHex(int64(w.header.Height)),
		},
		[]byte{},
	)

	return data
}

func (pow *Worker) tryOnce(nonce int64) {
	var hashInt big.Int
	var hash common.Hash

	data := pow.tryData(nonce)

	hash = sha256.Sum256(data)
	if math.Remainder(float64(nonce), 100000) == 0 {
		fmt.Printf("\r%x", hash)
	}
	hashInt.SetBytes(hash[:])

	if hashInt.Cmp(pow.target) == -1 {
		pow.header.Nonce = uint32(nonce)
	} else {
		nonce++
	}
	return
}

func (pow *Worker) Run() {
	var nonce int64

	for nonce < maxNonce {
		select {
		case <-pow.stop:
			return

		default:
			pow.tryOnce(nonce)
			pow.header.Nonce = uint32(nonce)
			break
		}
	}
}
