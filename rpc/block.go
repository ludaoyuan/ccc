package rpc

import (
	"core/types"
	"net/http"
)

func (c *RPCClient) GenesisBlock(r *http.Request, args *types.Nil, reply *types.Block) error {
	*reply = *types.GenesisBlock
	err := c.chain.DumpDB(types.GenesisBlock)
	if err != nil {
		return err
	}
	// hash, err := types.GenesisBlock.GenerateHash()
	// if err != nil {
	// 	return err
	// }

	// log.Println(hex.EncodeToString(hash[:]))
	// log.Println(hash)
	return nil
}
