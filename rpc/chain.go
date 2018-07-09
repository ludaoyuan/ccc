package rpc

import (
	"core"
	"core/types"
	"log"
	"net/http"
)

func (c *RPCClient) CreateBlockchain(r *http.Request, args *types.Nil, reply *types.Block) error {
	block, err := core.CreateGenesisBlock()
	if err != nil {
		log.Println(err.Error())
		return err
	}

	*reply = *block
	return nil
}

func (c *RPCClient) GetBlockCount(r *http.Request, args *types.Nil, reply *uint32) error {
	*reply = c.chain.GetBlockCount()
	return nil
}

func (c *RPCClient) GetBlockByNumber(r *http.Request, height *uint32, reply *types.Block) error {
	block, err := c.chain.GetBlockByNumber(*height)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	*reply = *block

	return nil
}

func (c *RPCClient) GetBlockByHash(r *http.Request, blockHashStr *string, reply *types.Block) error {
	block, err := c.chain.GetBlockByHash([]byte(*blockHashStr))
	if err != nil {
		log.Println(err.Error())
		return err
	}
	*reply = *block

	return nil
}

func (c *RPCClient) LastBlock(r *http.Request, args *types.Nil, reply *types.Block) error {
	*reply = *c.chain.LastBlock()
	return nil
}
