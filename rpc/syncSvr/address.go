package syncsvr

import (
	"common"
	"log"
	"net/rpc/jsonrpc"
)

const seedHost = "127.0.0.1:8080"

func (s *RPCS) Address(args *string, reply *common.Nil) error {
	s.syncMu.Lock()
	s.addrs = append(s.addrs, *args)
	s.syncMu.Unlock()

	return nil
}

func (s *RPCS) registerMyself() {
	client, err := jsonrpc.Dial("tcp", seedHost)
	if err != nil {
		log.Println(err.Error())
		return
	}

	err = client.Call("Address.Register", myHost, nil)
	if err != nil {
		log.Println(err.Error())
		return
	}
}
