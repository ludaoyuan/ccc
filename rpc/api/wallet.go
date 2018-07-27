package api

import (
	"common"
	"net/http"
)

func (api *API) Init(r *http.Request, args *common.Nil, hexHash *string) error {
	var err error
	*hexHash, err = api.walletSvr.Init()
	return err
}

func (api *API) ID(r *http.Request, args *common.Nil, hexHash *string) error {
	var err error
	*hexHash, err = api.walletSvr.ID()
	return err
}

func (api *API) GetBalance(r *http.Request, address *string, amount *int64) error {
	*amount = api.walletSvr.GetBalance(*address)
	return nil
}

// P2PHK: pay to public key hash
type CreateTxCmd struct {
	From  string
	To    string
	Value int64
}

func (api *API) CreateTx(r *http.Request, args *CreateTxCmd, hexHash *common.Nil) error {
	api.walletSvr.CreateTx(args.From, args.To, args.Value)
	return nil
}
