package rpc

import (
	"log"
	"net/http"

	"core/types"
)

// 列出某一个地址所有的UTXO
func (c *RPCClient) ListUTXOsByKey(r *http.Request, address *string, reply *types.TxOuts) error {
	var err error
	w := c.wallets.GetWallet(*address)
	pubkeyhash, err := w.PubKeyHash()
	if err != nil {
		log.Println(err.Error())
		return err
	}

	// TODO: 初始化移到程序启动初始化程序中
	utxo, err := c.chain.InitUTXOSet()
	if err != nil {
		log.Println(err.Error())
		return err
	}

	err = c.utxo.ToDB(utxo)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	utxos, err := c.utxo.FindUTXOs(pubkeyhash)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	*reply = *utxos

	return nil
}

// 账户余额
func (c *RPCClient) GetBalance(r *http.Request, args *string, reply *uint32) error {
	var err error

	utxos, err := c.utxo.FindUTXOs([]byte(*args))
	if err != nil {
		log.Println(err.Error())
		return err
	}

	var balance uint32
	for _, out := range utxos.Outs {
		balance += out.Value
	}

	*reply = balance
	return nil
}
