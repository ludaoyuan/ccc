package block

import (
	"math"
	"bytes"
	"crypto/sha256"
	"fmt"
	"time"

	"util"
)

const targetBits = 28

type ProofOfWork struct {
	Block	*Block
	Target	*uint64
}

func (pow *ProofOfWork) SetTarget() {
//	target := big.NewInt(1)
	target := uint64(1)
//	target.Lsh(target, uint(256 - targetBits))

	target = target << (64 - targetBits)

	pow.Target = &target
}

func (pow *ProofOfWork) prepareData(nonce uint64) []byte {
	data := bytes.Join([][]byte{pow.Block.Hash, util.IntToHex(int64(nonce))}, []byte{})

	return data
}

func (pow *ProofOfWork) Mining() (uint64, []byte) {
	var hashInt uint64
	var hash [32]byte
	var nonce uint64

	fmt.Printf("Mining the block containing \"%s\"", pow.Block.Data)

	go func() {
		p := time.Millisecond * 600
		fmt.Println()
		for {
			fmt.Printf(".")
			time.Sleep(p)
			fmt.Printf(".")
			time.Sleep(p)
			fmt.Printf(".")
			time.Sleep(p)
			fmt.Printf(".")
			time.Sleep(p)
			fmt.Printf(".")
			time.Sleep(p)
			fmt.Printf(".")
			time.Sleep(p)
			fmt.Printf("\b\b\b\b\b\b")
			fmt.Printf("      ")
			fmt.Printf("\b\b\b\b\b\b")
		}
	}()

	for nonce < math.MaxUint64 {
		data := pow.prepareData(nonce)
		hash = sha256.Sum256(data)
		hashInt = util.HexToUint(hash[:])

		if hashInt < *pow.Target {
			fmt.Printf("\b\n")
			fmt.Println("success!")
			fmt.Printf("%b\n", hashInt)
			fmt.Printf("%x\n", hashInt)
			fmt.Printf("%v\n\n", hashInt)

			return nonce, hash[:]
		}

		/*
		fmt.Printf("%b\n", *pow.target)
		fmt.Printf("%b\n", hashInt)
		fmt.Printf("%x\n", hashInt)
		fmt.Printf("%v\n\n", hashInt)
		*/

		nonce++

//		time.Sleep(time.Microsecond)
	}

	return nonce, hash[:]
}
