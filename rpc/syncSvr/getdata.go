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
	blocks, err := s.chain.FindCommonLastCommonBlock(args.IDs)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	for _, b := range blocks {
		*reply = append(*reply, b)
	}

	return nil
}

func (s *SyncServer) sendGetData(hashList []common.Hash) ([]*types.Block, error) {
	blocks := make([]*types.Block, 0)
	for _, addr := range s.addrs {
		client, err := jsonrpc.Dial("tcp", addr)
		if err != nil {
			log.Println(err.Error())
			return nil, err
		}

		args := getDataArgs{AddFrom: myHost, IDs: hashList}

		newBlocks := make([]*types.Block, 0)
		err = client.Call("SyncServer.GetData", &args, newBlocks)
		if err != nil {
			log.Println(err.Error())
			return nil, err
		}

		blocks = append(blocks, newBlocks...)
	}

	return blocks, nil
}
