package node

import (
	"core"
	"core/types"
	"log"
	"miner"
	"rpc/api"
	syncsvr "rpc/syncSvr"
	"wallet"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/filter"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

const chainPath = "./chaindb"

type Node struct {
	chainDB *leveldb.DB

	apiSvr    *api.API
	chainSvr  *core.BlockChain
	walletSvr *wallet.WalletSvr
	minerSvr  *miner.Miner
	syncSvr   *syncsvr.SyncServer
}

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
}

func NewNode() *Node {
	opts := opt.Options{
		ErrorIfExist: false,
		Strict:       opt.DefaultStrict,
		Compression:  opt.NoCompression,
		Filter:       filter.NewBloomFilter(10),
	}
	chaindb, err := leveldb.OpenFile(chainPath, &opts)
	if err != nil {
		log.Fatal(err.Error())
	}

	chainSvr, err := core.NewBlockChain(chaindb)
	if err != nil {
		log.Fatal(err.Error())
	}

	walletSvr := wallet.NewWalletSvr(chaindb, chainSvr)
	syncSvr := syncsvr.NewSyncServer(chainSvr)

	minerSvr := miner.NewMiner(chaindb, chainSvr, walletSvr.Coinbase())

	apiSvr := api.NewAPI(minerSvr, chainSvr, walletSvr)

	return &Node{
		chainDB:   chaindb,
		chainSvr:  chainSvr,
		walletSvr: walletSvr,
		minerSvr:  minerSvr,
		syncSvr:   syncSvr,
		apiSvr:    apiSvr,
	}
}

func (n *Node) Start() {
	go n.minerSvr.Start()
	go n.walletSvr.Start()
	go n.syncSvr.Start()
	go n.syncSvr.Start()
	go n.apiSvr.Start()

	for {
		select {
		case block, ok := <-n.minerSvr.NotifyNewLocalBlock():
			if ok {
				go n.syncSvr.BroadCastBlock(block)
			}
		case blocks, ok := <-n.syncSvr.NotifyNetBlocks():
			if ok {
				n.minerSvr.Update(blocks)
			}
		case block, ok := <-n.syncSvr.NotifyNetBlock():
			if ok {
				n.minerSvr.Update([]*types.Block{block})
			}
		case tx := <-n.syncSvr.NotifyNetTx():
			n.minerSvr.ReceiveTx(tx)
		}
	}
}
