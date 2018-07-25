package miner

import (
	"bytes"
	"common"
	"core/types"
	"fmt"
	"log"
	"math"
	"math/big"
)

var (
	maxNonce = math.MaxInt64
)

const targetBits = 5

// 挖矿难度调节
const DiffcultThreshold = 180

type Worker struct {
	header *types.BlockHeader
	target *big.Int

	stop <-chan struct{}
}

func NewWorker(header *types.BlockHeader, stop <-chan struct{}) *Worker {
	w := &Worker{
		header: header,
		target: big.NewInt(1),
		stop:   stop,
	}

	w.target.Lsh(w.target, uint(256-targetBits))
	return w
}

func (w *Worker) SetHeader(header *types.BlockHeader) {
	w.header = header
}

func (w *Worker) Stop() {
	close(w.stop)
	w.stop = make(chan struct{})
}

func (w *Worker) tryData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			common.IntToHex(int64(w.header.Version)),
			w.header.ParentHash,
			w.header.MerkleRoot,
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

func (pow *Worker) tryOnce(nonce int64) (int64, common.Hash) {
	var hashInt big.Int
	var hash common.Hash

	data := pow.tryData(nonce)

	hash = ssha256.Sum256(data)
	if math.Remainder(float64(nonce), 100000) == 0 {
		fmt.Printf("\r%x", hash)
	}
	hashInt.SetBytes(hash[:])

	if hashInt.Cmp(pow.target) == -1 {
		break
	} else {
		nonce++
	}
	return nonce, hash
}

func (pow *Worker) Run() {
	var hash common.Hash
	var nonce int64

	log.Printf("Mining the block containing \"%s\"\n", pow.block.Data)
	for nonce < maxNonce {
		select {
		case <-pow.stop:
			return

		default:
			nonce, hash = pow.tryOnce(nonce)
			pow.header.Nonce = nonce
			break
		}
	}
}
