package cli

import (
	"blockchain"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"os"
)

type BlockInfo struct {
	Height int64
	Data   string
}

type Args BlockInfo

// CLI responsible for processing command line arguments
type CLI struct {
	BC *blockchain.Blockchain
}

func (c *CLI) printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  addblock -data BLOCK_DATA - add a block to the blockchain")
	fmt.Println("  printchain - print all the blocks of the blockchain")
}

func (c *CLI) validateArgs() {
	if len(os.Args) < 2 {
		c.printUsage()
		os.Exit(1)
	}
}

func (c *CLI) addBlock(data string) {
	c.BC.AddBlock(data)
	fmt.Println("Success!")
}

func (c *CLI) AddBlock(args *Args, reply *BlockInfo) error {
	c.addBlock(args.Data)
	return nil
}

func (c *CLI) GetBlock(args *Args, reply *BlockInfo) error {
	reply.Data = c.BC.GetBlock(args.Height)
	return nil
}

func (c *CLI) Height(args *Args, reply *BlockInfo) error {
	reply.Height = c.BC.Height()
	return nil
}

// Run parses command line arguments and processes commands
func Run() {
	bc := blockchain.NewBlockchain()

	c := CLI{bc}
	rpc.Register(c)

	addr, err := net.ResolveTCPAddr("tcp", ":8080")
	if err != nil {
		log.Fatal(err.Error())
	}

	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Fatal(err.Error())
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		rpc.ServeConn(conn)
	}

}

// jsonrpc
func Run2() {
	bc := blockchain.NewBlockchain()

	c := CLI{bc}
	rpc.Register(c)

	addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1")
	if err != nil {
		log.Fatal(err.Error())
	}

	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Fatal(err.Error())
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		jsonrpc.ServeConn(conn)
	}
}
