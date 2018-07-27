// 同步区块服务提供
package syncsvr

import (
	"core"
	"core/types"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"sync"
)

const myHost = "127.0.0.1:4501"
const needUpdateBlock = 6

type Addresses []string

type SyncServer struct {
	syncMu sync.RWMutex

	chain *core.BlockChain
	addrs Addresses

	netTxMsg     chan *types.Transaction
	netBlockMsgs chan []*types.Block
	netBlockMsg  chan *types.Block
}

func (s *SyncServer) NotifyNetBlocks() <-chan []*types.Block {
	return s.netBlockMsgs
}

func (s *SyncServer) NotifyNetBlock() <-chan *types.Block {
	return s.netBlockMsg
}

func (s *SyncServer) NotifyNetTx() <-chan *types.Transaction {
	return s.netTxMsg
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

func NewSyncServer(chain *core.BlockChain) *SyncServer {
	return &SyncServer{
		chain:        chain,
		addrs:        make([]string, 0),
		netTxMsg:     make(chan *types.Transaction),
		netBlockMsgs: make(chan []*types.Block),
		netBlockMsg:  make(chan *types.Block),
	}
}

func (s *SyncServer) Stop() {
	close(s.netTxMsg)
	close(s.netBlockMsgs)
}

func (s *SyncServer) Start() {
	rpcs := RPCS(*s)

	rpcs.registerMyself()
	rpc.Register(&rpcs)

	tcpAddr, err := net.ResolveTCPAddr("tcp", myHost)
	if err != nil {
		log.Fatal(err.Error())
	}

	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		log.Fatal(err.Error())
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err.Error())
		}
		jsonrpc.ServeConn(conn)
	}
}

type RPCS SyncServer
