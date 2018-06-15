package main

import (
	"cli/httpclient/block"

	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/block/add", httpblock.HandleAddBlock)
	http.HandleFunc("/block/get", httpblock.HandleGetBlock)
	http.HandleFunc("/chain/height", httpblock.HandleHeight)

	err := http.ListenAndServe("127.0.0.1:8080", nil)
	if err != nil {
		log.Println(err.Error())
	}
}
