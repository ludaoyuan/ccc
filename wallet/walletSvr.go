package wallet

import (
	"bytes"
	"common"
	"core"
	"core/types"
	"crypto/elliptic"
	"encoding/gob"
	"io/ioutil"
	"log"
	"os"
	"sync"

	"github.com/syndtr/goleveldb/leveldb"
)

const walletPath = "./walletdb"

type WalletSvr struct {
	chain   *core.BlockChain
	chainDB *leveldb.DB

	myWallet *Wallet

	localTx chan *types.Transaction
	utxos   map[string]types.TxOuts

	quit   chan struct{}
	quitMu sync.Mutex
}

func (w *WalletSvr) Start() {
	w.quitMu.Lock()
	w.quitMu.Unlock()
}

func NewWalletSvr(chaindb *leveldb.DB, chain *core.BlockChain) *WalletSvr {
	return &WalletSvr{
		chain:   chain,
		chainDB: chaindb,
		utxos:   make(map[string]types.TxOuts),
	}
}

func (ws *WalletSvr) Coinbase() common.Hash {
	if ws.myWallet == nil {
		log.Println("init first!")
		addr, _ := ws.Init()
		log.Printf("Your Wallet Address is %s\n", addr)
		return common.ZeroHash
	}
	return ws.myWallet.PubKeyHash()
}

func (ws *WalletSvr) Init() (string, error) {
	w, err := NewWallet()
	if err != nil {
		log.Println(err.Error())
		return "", err
	}

	gob.Register(elliptic.P256())

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err = enc.Encode(w)
	if err != nil {
		log.Println(err.Error())
		return "", err
	}

	err = ioutil.WriteFile(walletPath, buf.Bytes(), 0644)
	if err != nil {
		log.Println(err.Error())
		return "", err
	}

	ws.myWallet = w
	return w.Address(), nil
}

func (ws *WalletSvr) ID() (string, error) {
	_, err := os.Stat(walletPath)

	if err != nil && os.IsNotExist(err) {
		log.Println(err.Error())
		return "", err
	}

	fileContent, err := ioutil.ReadFile(walletPath)
	if err != nil {
		log.Println(err.Error())
		return "", err
	}

	var w Wallet
	gob.Register(elliptic.P256())
	dec := gob.NewDecoder(bytes.NewReader(fileContent))
	err = dec.Decode(&w)
	if err != nil {
		log.Println(err.Error())
		return "", err
	}

	ws.myWallet = &w

	return w.Address(), nil
}
