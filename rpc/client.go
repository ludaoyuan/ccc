package rpc

import (
	"core"
	"log"
	"net/http"
	"wallet"

	"github.com/gorilla/rpc/v2"
	"github.com/gorilla/rpc/v2/json2"
	"github.com/gorilla/websocket"
)

const rpcHost = ":8080"

type RPCClient struct {
	wallets *wallet.Wallets
	// utxos   map[string]types.TxOuts
	utxo   *core.UTXOSet
	chain  *core.Blockchain
	txpool *core.TxPool
}

// curl -X POST -H "Content-Type:application/json" -d '{"jsonrpc":"2.0", "method":"RPCClient.Say","params":{"Who":"sang"},"id":1}' http://localhost:1234
func (c *RPCClient) Start() {
	s := rpc.NewServer()
	s.RegisterCodec(json2.NewCodec(), "application/json")
	s.RegisterService(c, "")

	http.Handle("/", s)
	// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// })

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ws, err := websocket.Upgrade(w, r, nil, 0, 0)
		if err != nil {
			log.Println(err.Error())
			return
		}
		defer ws.Close()

		for {
			mt, message, err := ws.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				break
			}
			log.Printf("recv: %s", message)
			err = ws.WriteMessage(mt, message)
			if err != nil {
				log.Println("write:", err)
				break
			}
		}
	})

	log.Fatal(http.ListenAndServe(rpcHost, nil))
}

const minerAddr = "1AeGtuczZ6aHoZRkWWHBWpUjeY3HxAe5ie"

func NewRPCClient() *RPCClient {
	chain := core.NewBlockchain()

	ws, err := wallet.NewWallets()
	if err != nil {
		log.Fatal(err.Error())
	}

	w := ws.GetWallet(minerAddr)

	utxo, err := core.NewUTXOSet(chain)
	if err != nil {
		log.Fatal(err.Error())
	}

	txpool := core.NewTxPool(w.PublicKey, chain, utxo)

	return &RPCClient{
		wallets: ws,
		utxo:    utxo,
		chain:   chain,
		txpool:  txpool,
	}
}
