package main

import (
	"log"
	"net"
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
	log.Printf("new peer:%s\n", *args)
	a.peers[*args] = struct{}{}
	a.broadcast(*args)
	return nil
}

func callRPC(addr, to string) {
	client, err := jsonrpc.Dial("tcp", to)
	if err != nil {
		log.Println(err.Error())
		return
	}

	err = client.Call("SyncServer.Address", addr, nil)
	if err != nil {
		log.Println(err.Error())
		return
	}
}

func (a *Address) broadcast(newAddr string) {
	a.aMu.RLock()
	for k, _ := range a.peers {
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
	// rpc.Register(addr)
	// rpc.HandleHTTP()
	//
	// log.Println(http.ListenAndServe(seedHost, nil))

	rpc.Register(addr)

	tcpAddr, err := net.ResolveTCPAddr("tcp", seedHost)
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
			continue
		}
		jsonrpc.ServeConn(conn)
	}

}
