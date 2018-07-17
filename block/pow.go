package block

import (
	"math"
	"bytes"
	"crypto/sha256"
	"time"
	"fmt"

	"util"
)

var targetBits = 1

type ProofOfWork struct {
	Block	*Block
	Target	*uint64
}

func (pow *ProofOfWork) SetTarget() {
//	target := big.NewInt(1)
	target := uint64(1)
//	target.Lsh(target, uint(256 - targetBits))

	target = target << uint(64 - targetBits)

	pow.Target = &target
}

func (pow *ProofOfWork) prepareData(nonce uint64) []byte {
	data := bytes.Join([][]byte{pow.Block.PrevBlockHash, pow.Block.Data, util.IntToHex(pow.Block.Timestamp), util.IntToHex(int64(targetBits)), util.IntToHex(int64(nonce))}, []byte{})

	return data
}

func (pow *ProofOfWork) Mining() (uint64, []byte, bool) {
	var hashInt uint64
	var hash [32]byte
	var nonce uint64

	fmt.Printf("Mining the block containing \"%s\"\n", pow.Block.Data)

	t0 := time.Now().Unix()
	for nonce < math.MaxUint64 {
		for i := 0; i < 600000; i++ {
			data := pow.prepareData(nonce)
			hash = sha256.Sum256(data)
			hashInt = util.HexToUint(hash[:])

			if hashInt < *pow.Target {
				fmt.Printf("\b\n")
				fmt.Println("success!")
				fmt.Printf("%b\n", hashInt)
				fmt.Printf("%x\n", hashInt)
				fmt.Printf("%v\n\n", hashInt)

				t1 := time.Now().Unix()
				fmt.Println("\r", t1 - t0, "s\n")
				return nonce, hash[:], (t1 - t0) > 180
			}

			/*
			fmt.Printf("%b\n", *pow.target)
			fmt.Printf("%b\n", hashInt)
			fmt.Printf("%x\n", hashInt)
			fmt.Printf("%v\n\n", hashInt)
			*/

			if i % 100000 == 0 {
				fmt.Printf(".")
			}
			nonce++

	//		time.Sleep(time.Microsecond)
		}
		fmt.Printf("\b\b\b\b\b\b")
		fmt.Printf("      ")
		fmt.Printf("\b\b\b\b\b\b")
	}
	t1 := time.Now().Unix()

	return nonce, hash[:], (t1 - t0) > 180
}
