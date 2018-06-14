package main

import (
	"cli"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/rpc"
)

var (
	client *rpc.Client
	mp     map[string]string
)

func init() {
	flag.StringVar(&host, "host", "127.0.0.1:8080", "Method")
	client, err = rpc.Dial("tcp", host)
	if err != nil {
		log.Fatal(err.Error())
	}

	mf = make(map[string]string)
	mf["add"] = "CLI.AddBlock"
	mf["get"] = "CLI.GetBlock"
	mf["height"] = "CLI.Height"
}

func handleFunc(cmd string) {
}

func main() {
	args := cli.Args{Data: data, Height: height}
	var reply cli.BlockInfo

	fmt.Println("vim-go")
	http.HandleFunc("add", handleFunc)
	// http.handle()
}
