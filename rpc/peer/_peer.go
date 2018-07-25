// 维护相邻节点列表服务
package peer

import (
	"log"
	"net/rpc"
	"sync"
	"time"
)

const seedHost = "127.0.0.1:8080"

type Peers map[string]struct{}

// TODO: 锁的粒度比较大
type PeerStore struct {
	mu sync.Mutex

	Peers    Peers
	NewPeers chan string
}

func NewPeerStore(seedHost string) (*PeerStore, error) {
	ps := &PeerStore{
		Peers: make(map[string]struct{}),
	}

	return ps, nil
}

// 返回节点列表
func (ps *PeerStore) Peers() []string {
	ps.Lock()
	defer ps.Unlock()

	addrs := make([]string, 0, len(ps.Peers))
	for _, addr := range ps.Peers {
		addrs = append(addrs, addr)
	}

	return addrs
}

func (ps *PeerStore) Get() {
	client, err := rpc.Dial("tcp", seedHost)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	peers := make(map[string]struct{})
	err := client.Call("Address.GetAll", nil, &peers)
	if err != nil {
		log.Println(err.Error())
		return
	}

	ps.Lock()
	ps.Peers = peers
	ps.Unlock()

	client.Close()
}

func (ps *PeerStore) Init() {
	ps.Get()
}

func (ps *PeerStore) Start() {
	for {
		select {
		case <-time.After(gap * time.Second):
			ps.Get()
		}
	}
}
