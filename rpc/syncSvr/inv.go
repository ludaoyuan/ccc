package syncsvr

import (
	"common"
	"log"
	"net/rpc/jsonrpc"
)

type invCMD struct {
	AddFrom       string
	chainHashList []common.Hash
}

type HashList []common.Hash

func (s *SyncServer) Inv(args *invCMD, hashList *HashList) error {
	missingList, err := s.chain.GetMissingBlocksHash(args.chainHashList)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	for _, hash := range missingList {
		*hashList = append(*hashList, hash)
	}

	return nil
}

func (s *SyncServer) sendGetInv(to string) (HashList, error) {
	client, err := jsonrpc.Dial("tcp", to)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	chainList, err := s.chain.ChainList()
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	args := invCMD{AddFrom: myHost, chainHashList: chainList}
	hashList := make(HashList, 0)
	err = client.Call("SyncServer.Inv", &args, &hashList)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return hashList, nil
}
