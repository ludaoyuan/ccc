package rpc

import (
	"core/types"
	"log"

	"net/http"
)

type WalletSAddressReply struct {
	Addresses []string
}

// 列出本地所有钱包
func (c *RPCClient) ListWallets(r *http.Request, args *types.Nil, reply *WalletSAddressReply) error {
	address := c.wallets.GetAddresses()
	reply.Addresses = append(reply.Addresses, address[:]...)

	return nil
}

// 新建钱包
func (c *RPCClient) CreateWallet(r *http.Request, args *types.Nil, reply *string) error {
	address, err := c.wallets.CreateWallet()
	if err != nil {
		return err
	}

	*reply = address

	w := c.wallets.GetWallet(address)
	key, err := w.PubKeyHash()
	if err != nil {
		return err
	}
	log.Println(key)

	err = c.wallets.DumpWallet()
	if err != nil {
		return err
	}

	return nil
}
