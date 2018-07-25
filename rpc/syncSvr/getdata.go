package syncsvr

import (
	"common"
	"core/types"
	"log"
	"net/rpc/jsonrpc"
)

type getDataArgs struct {
	AddFrom string
	IDs     []common.Hash
}

type Blocks []*types.Block

func (s *SyncServer) GetData(args *getDataArgs, reply *Blocks) error {
	blocks, err := s.chain.FindCommonLastCommonBlock(args.ID)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	for _, b := range blocks {
		*reply = append(*reply, b)
	}

	return nil
}

func (s *SyncServer) sendGetData(to string, hashList []common.Hash) (*types.Block, error) {
	client, err := jsonrpc.Dial("tcp", to)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	args := getDataArgs{AddFrom: myHost, IDs: hashList}
	blocks := make(Blocks, 0)

	err = client.Call("SyncServer.GetData", &args, blocks)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return blocks, nil
}
