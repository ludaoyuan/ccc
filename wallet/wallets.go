package wallet

import (
	"bytes"
	"common"
	"crypto/elliptic"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"log"
)

const walletFile = "./data/wallet.dat"

// Wallets stores a collection of wallets
type Wallets struct {
	Wallets map[string]*Wallet
}

// NewWallets creates Wallets and fills it from a file if it exists
func NewWallets() (*Wallets, error) {
	wallets := Wallets{}
	wallets.Wallets = make(map[string]*Wallet)

	err := wallets.LoadWallet()

	return &wallets, err
}

// CreateWallet adds a Wallet to Wallets
func (ws *Wallets) CreateWallet() (string, error) {
	wallet, err := NewWallet()
	if err != nil {
		log.Println(err.Error())
		return "", err
	}

	hashPubKey, err := wallet.GetAddress()
	if err != nil {
		log.Println(err.Error())
		return "", err
	}
	address := fmt.Sprintf("%s", hashPubKey)

	ws.Wallets[address] = wallet

	return address, nil
}

// GetAddresses returns an array of addresses stored in the wallet file
func (ws *Wallets) GetAddresses() []string {
	var addresses []string

	for address := range ws.Wallets {
		addresses = append(addresses, address)
	}

	return addresses
}

func (ws Wallets) GetWallet(address string) Wallet {
	return *ws.Wallets[address]
}

func (ws *Wallets) LoadWallet() error {
	exists := common.CheckPath(walletFile)
	if !exists {
		return nil
	}

	fileContent, err := ioutil.ReadFile(walletFile)
	if err != nil {
		log.Println(err)
		return err
	}

	var wallets Wallets
	gob.Register(elliptic.P256())
	decoder := gob.NewDecoder(bytes.NewReader(fileContent))
	err = decoder.Decode(&wallets)
	if err != nil {
		log.Println(err)
		return err
	}

	ws.Wallets = wallets.Wallets

	return nil
}

func (ws Wallets) DumpWallet() error {
	var content bytes.Buffer

	gob.Register(elliptic.P256())

	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(ws)
	if err != nil {
		log.Println(err)
		return err
	}

	err = ioutil.WriteFile(walletFile, content.Bytes(), 0644)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
