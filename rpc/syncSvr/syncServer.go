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

const myHost = ":8080"
const needUpdateBlock = 6

type Addresses []string

type SyncServer struct {
	syncMu sync.RWMutex

	chain *core.BlockChain
	Addrs Addresses

	netTxMsg    chan *types.Transaction
	netBlockMsg chan *types.Block
	netNodeMsg  chan string
}

func (s *SyncServer) NotifyNetBlock() <-chan types.Block {
	return s.netBlockMsg
}

func (s *SyncServer) NotifyNetTx() <-chan types.Transaction {
	return s.netTxMsg
}

func (s *SyncServer) NotifyNetNodeMsg() <-chan string {
	return s.netNodeMsg
}

func NewSyncServer(chain *core.Blockchain) *SyncServer {
	return &SyncServer{
		chain:       chain,
		Addrs:       make([]string, 0),
		netTxMsg:    make(chan *types.Transaction),
		netBlockMsg: make(chan *types.Block),
		netNodeMsg:  make(chan string),
	}
}

func (s *SyncServer) Start() {
	rpc.Register(sync)
	rpc.HandleHTTP()

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
