package node

import (
	"log"
	"miner"
	"wallet"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/filter"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

const chainPath = "./chaindb"

type Node struct {
	chainDB *leveldb.DB

	chainSvr  *core.BlcokChain
	walletSvr *wallet.WalletSvr
	minerSvr  *miner.Miner
	syncSvr   *syncsvr.SyncSvr
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
	syncSvr := syncSvr.NewSyncServer(chainSvr)

	minerSvr := miner.newminer(chaindb, chainSvr, walletSvr.Coinbase())

	return Node{
		chainDB:   chaindb,
		chainSvr:  chainSvr,
		walletSvr: walletSvr,
		minerSvr:  minerSvr,
		syncSvr:   syncSvr,
	}
}

func (n *Node) Start() error {
	go n.minerSvr.Start()
	go n.walletSvr.Start()
	go n.syncSvr.Start()
	go n.syncSvr.Start()

	for {
		select {
		case block := <-n.minerSvr.NotifyNewLocalBlock:
			go n.syncSvr.BroadCastBlock(block)
		case block := <-n.syncSvr.NotifyNetBlock():
			n.minerSvr.MergeBlock()
		case tx := <-n.syncSvr.NotifyNetTx():
			n.minerSvr.ReceiveTx(tx)
		}
	}
}
