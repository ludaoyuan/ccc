package syncsvr

import (
	"common"
	"core/types"
)

func (s *SyncServer) NewBlock(block *types.Block, reply *common.Nil) error {
	if s.chain.VerifyBlock(block) {
		s.netBlockMsg <- block
	}
	return nil
}
