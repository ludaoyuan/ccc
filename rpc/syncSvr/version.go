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
	if s.chain.Height() >= args.BestHeight+PreBlocks {
		s.sendVersion(args.AddFrom)
		return nil
	}

	if s.chain.Height() < args.BestHeight() {
		hashList, err := s.sendGetInv(args.AddFrom)
		if err != nil {
			log.Println(err.Error())
			return err
		}
	}

	for _, hash := range hashList {
		blocks, err := s.sendGetData(to, hash)
		if err != nil {
			log.Println(err.Error())
			return err
		}
		s.minerSvr.Update(blocks)
	}

	if s.isUnkonwnNode(args.AddFrom) {
		s.peerSvr.AddNewNode(args.AddFrom)
	}
	return nil
}

func (s *SyncServer) sendGetVersion(to string) {
	client, err := jsonrpc.Dial("tcp", vinfo.AddFrom)
	if err != nil {
		log.Println(err.Error())
		return
	}

	args := &versionArgs{Version: 0, BestHeight: s.chain.Height(), AddFrom: myHost}
	err = client.Call("SyncServer.Version", args, common.Nil)
	if err != nil {
		log.Println(err.Error())
		return
	}
}
