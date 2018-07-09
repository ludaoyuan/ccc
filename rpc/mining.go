package rpc

import (
	"net/http"
)

type MiningCmd struct {
}

// 测试使用, 产生number 个blocks 并返回他们的hash
func (c *RPCClient) Mining(r *http.Request, address *string, reply []string) error {
	// reply = c.NewGenerateCmd(number)
	return nil
}
