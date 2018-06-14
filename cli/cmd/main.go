package main

import (
	"cli"
	"flag"
	"log"
	"net/rpc"
)

var (
	host   string
	cmd    string
	height int64
	data   string
	err    error
	client *rpc.Client
)

var mf map[string]string

func init() {
	mf = make(map[string]string)
	mf["add"] = "CLI.AddBlock"
	mf["get"] = "CLI.GetBlock"
	mf["height"] = "CLI.Height"

	flag.StringVar(&host, "host", "127.0.0.1:8080", "Method")
	flag.StringVar(&cmd, "cmd", "Height", "Remote Method:[add get height]")
	flag.StringVar(&data, "data", "", "Data")
	flag.Int64Var(&height, "height", 0, "Add Block: Data")

	client, err = rpc.Dial("tcp", host)
	if err != nil {
		log.Fatal(err.Error())
	}

	flag.Parse()
}

func main() {
	args := cli.Args{Data: data, Height: height}
	var reply cli.BlockInfo

	if cmd == "" {
		log.Fatal("Usage --cmd [] --data [] --height --host")
	}
	err = client.Call(mf[cmd], args, &reply)
	if err != nil {
		log.Fatal(err.Error())
	}
}
