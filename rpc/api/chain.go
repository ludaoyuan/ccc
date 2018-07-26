package api

import (
	"common"
	"core/types"
	"log"
	"net/http"
)

func (api *API) GenesisBlock(r *http.Request, args *common.Nil, block *types.Block) error {
	*block = *api.chainSvr.GenesisBlock()
	return nil
}

func (api *API) Height(r *http.Request, args *common.Nil, height *int32) error {
	*height = api.chainSvr.Height()
	return nil
}

func (api *API) GetBlock(r *http.Request, hexHash *string, block *types.Block) error {
	var hash common.Hash
	hash.FromHex(*hexHash)

	b, err := api.chainSvr.GetBlock(hash)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	*block = *b
	return nil
}

func (api *API) LastBlock(r *http.Request, hash *common.Nil, block *types.Block) error {
	*block = *api.chainSvr.LastBlock()
	return nil
}
