package rpc

import (
	"core/types"
	"net/http"
)

func (c *RPCClient) GenesisBlock(r *http.Request, args *types.Nil, reply *types.Block) error {
	*reply = *types.GenesisBlock
	return nil
}
