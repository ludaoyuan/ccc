package syncsvr

import (
	"common"
	"core/types"
	"log"
	"net/rpc/jsonrpc"
)

func (s *SyncServer) Address(args *common.Nil, addr *string) error {
	s.syncMu.Lock()
	s.addrs = append(s.addrs, *addr)
	s.syncMu.Unlock()

	return nil
}

func (s *SyncServer) BroadCastBlock(block *types.Block) {
	for _, addr := range s.addrs {
		client, err := jsonrpc.Dial("tcp", addr)
		if err != nil {
			log.Println(err.Error())
			return
		}

		err = client.Call("SyncServer.Version", block, nil)
		if err != nil {
			log.Println(err.Error())
			return
		}
	}
}

func (s *SyncServer) isUnkonwnNode(newAddr string) bool {
	for _, addr := range s.addrs {
		if addr == newAddr {
			return true
		}
	}

	return false
}
