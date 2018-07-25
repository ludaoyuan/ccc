package syncsvr

import (
	"common"
	"log"
	"net/rpc/jsonrpc"
)

const PreBlocks = 6

type versionArgs struct {
	Version    int
	BestHeight int
	AddFrom    string
}

func (s *SyncServer) Version(args *versionArgs, reply *common.Nil) error {
	if s.chain.Height() >= uint32(args.BestHeight+PreBlocks) {
		s.sendVersion(args.AddFrom)
		return nil
	}

	var hashList []common.Hash
	var err error
	if s.chain.Height() < uint32(args.BestHeight) {
		hashList, err = s.sendGetInv(args.AddFrom)
		if err != nil {
			log.Println(err.Error())
			return err
		}
	}

	blocks, err := s.sendGetData(hashList)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	s.netBlockMsgs <- blocks

	if s.isUnkonwnNode(args.AddFrom) {
		s.addrs = append(s.addrs, args.AddFrom)
	}
	return nil
}

func (s *SyncServer) sendVersion(to string) {
	client, err := jsonrpc.Dial("tcp", to)
	if err != nil {
		log.Println(err.Error())
		return
	}

	args := &versionArgs{Version: 0, BestHeight: int(s.chain.Height()), AddFrom: myHost}
	err = client.Call("SyncServer.Version", args, nil)
	if err != nil {
		log.Println(err.Error())
		return
	}
}
