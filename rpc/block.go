package rpc

import (
	"core/types"
	"net/http"
)

func (c *RPCClient) GenesisBlock(r *http.Request, args *types.Nil, reply *types.Block) error {
	hash, err := types.GenesisBlock.GenerateHash()
	if err != nil {
		return err
	}

	types.GenesisBlock.SetHash(hash)

	err = c.chain.DumpDB(types.GenesisBlock)
	if err != nil {
		return err
	}

	// TODO: 不可导出
	// c.chain.Foreach()
	*reply = *types.GenesisBlock
	return nil
}
