package main

import (
	"common"
	"log"
	"net/http"
	"net/rpc"
	"net/rpc/jsonrpc"
	"sync"
)

const seedHost = "127.0.0.1:8080"

type Peers map[string]struct{}

type ANil struct {
}

type Address struct {
	aMu   sync.RWMutex
	peers Peers
}

func (a *Address) Info(args *ANil, reply *string) error {
	log.Println("Debug Into Info")
	*reply = "sanghaifa"
	return nil
}

func (a *Address) GetAll(args *ANil, reply *Peers) error {
	log.Println("IN")
	*reply = a.peers
	return nil
}

func (a *Address) Register(args *string, reply *ANil) error {
	a.peers[*args] = struct{}{}
	a.Broadcast(*args)
	return nil
}

func callRPC(addr, to string) {
	client, err := jsonrpc.Dial("tcp", to)
	if err != nil {
		log.Println(err.Error())
		return
	}

	err = client.Call("SyncServer.Version", addr, &common.Nil)
	if err != nil {
		log.Println(err.Error())
		return
	}
}

func (a *Address) Broadcast(newAddr string) {
	a.aMu.RLock()
	for k, v := range a.peers {
		if k != newAddr {
			callRPC(newAddr, k)
		}
	}
	a.aMu.RUnlock()
}

func main() {
	log.Println("start seed server")
	addr := &Address{
		peers: make(map[string]struct{}),
	}
	rpc.Register(addr)
	rpc.HandleHTTP()

	log.Println(http.ListenAndServe(seedHost, nil))
}
