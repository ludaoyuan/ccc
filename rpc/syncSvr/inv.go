package syncsvr

import (
	"common"
	"log"
	"net/rpc/jsonrpc"
)

type invCMD struct {
	AddFrom       string
	chainHashList []*common.Hash
}

type HashList []common.Hash

func (s *SyncServer) Inv(inv *invCMD, hashList *HashList) error {
	newBlocks, err := s.chain.FindCommonLastCommonBlock(*hashList)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	for _, b := range newBlocks {
		*blocks = append(*blocks, b)
	}

	for _, hash := range newBlocks {
		s.sendGetData(to, hash)
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
