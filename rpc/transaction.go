package rpc

import (
	"common"
	"core/types"
	"crypto/ecdsa"
	"log"
	"net/http"
)

// 获取交易
func (c *RPCClient) GetTransaction(r *http.Request, txhash *string, reply *types.Transaction) error {
	hash := common.ToHash([]byte(*txhash))
	tx, err := c.chain.FindTransaction(hash)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	*reply = *tx
	return nil
}

type SignTransactionCmd struct {
	Tx      *types.Transaction
	privKey ecdsa.PrivateKey
}

func (c *RPCClient) SignTransaction(r *http.Request, args *SignTransactionCmd, reply *types.Nil) error {
	err := c.chain.SignTransaction(args.Tx, args.privKey)

	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

type SendTransactionCmd struct {
	From  string
	To    string
	Value uint32
}

func (c *RPCClient) SendTransaction(r *http.Request, args *SendTransactionCmd, reply *types.Nil) error {
	// if !wallet.ValidateAddress(args.From) {
	// 	err := errors.New("ERROR: Sender address is not valid")
	// 	log.Println(err.Error())
	// 	return err
	// }
	// if !wallet.ValidateAddress(args.To) {
	// 	err := errors.New("ERROR: Reciever address is not valid")
	// 	log.Println(err.Error())
	// 	return err
	// }
	//
	// w := c.wallets.GetWallet(args.From)
	//
	// tx, err := w.CreateTx(c.chain, []byte(args.To), args.Value, c.utxo)
	// if err != nil {
	// 	log.Println(err.Error())
	// 	return err
	// }
	//
	// // TODO:
	// // tx.TxIn :=
	return nil
}
