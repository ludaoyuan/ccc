package chain

import (
	"bytes"
	"io/ioutil"
//	"crypto/sha256"
//	"encoding/binary"
	"encoding/gob"
	"time"
	"log"
)

var (
	block Block
	err error
	buf bytes.Buffer
	bs []byte
	b []byte
)

type Block struct {
	Timestamp	int64
	Data	map[string]interface{}
	Hash	[]byte
	PrevBlockHash	string
}

func GenesisCreate() error {
	gob.Register(map[string]interface{} {})
	encoder := gob.NewEncoder(&buf)

	block.Timestamp = time.Now().UnixNano()

	err = encoder.Encode(block.Timestamp)
	if err != nil {
		return err
	}

	bs, err = ioutil.ReadAll(&buf)

	block.Data = make(map[string]interface{})

	block.Data["x"] = "xxxx"
	block.Data["y"] = "yyyy"
	block.Data["z"] = "zzzz"

	err = encoder.Encode(block.Data)
	if err != nil {
		return err
	}

	b, err = ioutil.ReadAll(&buf)
	if err != nil {
		return err
	}

	bs = append(bs[:], b[:]...)

	err = encoder.Encode(block.PrevBlockHash)
	if err != nil {
		return err
	}

	b, err = ioutil.ReadAll(&buf)
	if err != nil {
		return err
	}

	bs = append(bs[:], b[:]...)

	log.Println(string(bs))

	return nil
}

/*
func (block Block)Get() error {
	Decoder := gob.NewDecoder(&buf)
	err := Decoder.Decode(block)
}
*/

/*
func (block * Block)validation(blocks []Block) {
	var block Block
	var err error

	bs := bytes.Buffer{}

	encoder := gob.NewEncoder()

	range
	err = encoder.Encode(blocks)
	blocks[len(blocks)]
}
*/
