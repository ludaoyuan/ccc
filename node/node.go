package node

import (
	"miner"
	"rpc"
)

// 节点分类: 全节点, SPV节点
type Node struct {
	rpc   *rpc.RPCClient
	miner *miner.Miner
}

// 服务启动入口
func (n *Node) Start() {
	// 启动监听服务
	n.rpc.Start()
}

// 启动监听服务
// func NewNode(ctx *cli.Context) *Node {
// 	return &Node{
// 		rpc: rpc.NewRPCClient(),
// 	}
// }

// 钱包相关初始化
func init() {
}

func NewNode() *Node {
	return &Node{
		rpc: rpc.NewRPCClient(),
	}
}
