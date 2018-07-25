package syncsvr

import (
	"common"
)

func (s *SyncServer) Address(args *common.Nil, addr *string) error {
	s.syncMu.Lock()
	s.Addrs = append(s.addrs, *addr)
	s.syncMu.Unlock()

	return nil
}
