package syncsvr

import (
	"common"
	"log"
	"net/rpc/jsonrpc"
)

type InvCMD struct {
	AddFrom       string
	ChainHashList []common.Hash
}

type HashList []common.Hash

func (s *RPCS) Inv(args *InvCMD, hashList *HashList) error {
	missingList, err := s.chain.GetMissingBlocksHash(args.ChainHashList)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	for _, hash := range missingList {
		*hashList = append(*hashList, hash)
	}

	return nil
}

func (s *RPCS) sendGetInv(to string) (HashList, error) {
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

	args := InvCMD{AddFrom: myHost, ChainHashList: chainList}
	hashList := make(HashList, 0)
	err = client.Call("SyncServer.Inv", &args, &hashList)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return hashList, nil
}
