package cli

import (
	"log"
	"net/rpc"
	"net/rpc/jsonrpc"
	"testing"
)

var (
	client  *rpc.Client
	client2 *rpc.Client
	err     error
)

const OPT = 0

func initClient() {
	client, err = rpc.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		log.Fatal(err.Error())
	}
}

func initClient2() {
	client2, err = jsonrpc.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		log.Fatal(err.Error())
	}
}

func Init() *rpc.Client {
	initClient()
	initClient2()

	switch OPT {
	case 0:
		return client
	default:
		return client2
	}
	return nil
}

func TestAddBlock(t *testing.T) {
	c := Init()

	args := Args{Data: "Block2"}
	var reply BlockInfo
	err = c.Call("CLI.AddBlock", args, &reply)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func TestGetBlock(t *testing.T) {
	c := Init()

	args := Args{Height: 2}
	var reply BlockInfo
	err = c.Call("CLI.GetBlock", args, &reply)
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Printf("%+v\n", reply)
}

func TestHeight(t *testing.T) {
	c := Init()

	var args Args
	var reply BlockInfo
	err = c.Call("CLI.Height", args, &reply)
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Println(reply.Height)
}
