package api

import (
	"wallet"

	"core"
	"log"
	"miner"
	"net/http"

	"github.com/gorilla/rpc/v2"
	"github.com/gorilla/rpc/v2/json2"
	"github.com/gorilla/websocket"
)

const apiHost = "127.0.0.1:4500"

type API struct {
	minerSvr  *miner.Miner
	chainSvr  *core.BlockChain
	walletSvr *wallet.WalletSvr
}

func NewAPI(miner *miner.Miner, chain *core.BlockChain, wallets *wallet.WalletSvr) *API {
	return &API{miner, chain, wallets}
}

func (api *API) Start() {
	s := rpc.NewServer()
	s.RegisterCodec(json2.NewCodec(), "application/json")
	s.RegisterService(api, "")

	http.Handle("/", s)
	// http.HandleFunc("ws", handleWs)

	log.Println(http.ListenAndServe(apiHost, nil))
}

func handleWs(w http.ResponseWriter, r *http.Request) {
	ws, err := websocket.Upgrade(w, r, nil, 0, 0)
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer ws.Close()
}
