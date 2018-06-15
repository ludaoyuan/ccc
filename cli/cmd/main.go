package main

import (
	"cli"
	"flag"
	"log"
	"net/rpc"
	"os"
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
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	checkinghealthy()

	mf = make(map[string]string)
	mf["add"] = "CLI.AddBlock"
	mf["get"] = "CLI.GetBlock"
	mf["height"] = "CLI.Height"

	flag.StringVar(&host, "host", "127.0.0.1:8081", "Method")
	flag.StringVar(&cmd, "cmd", "Height", "Remote Method:[add get height]")
	flag.StringVar(&data, "data", "", "Data")
	flag.Int64Var(&height, "height", 0, "Add Block: Data")

	client, err = rpc.Dial("tcp", host)
	if err != nil {
		log.Fatal(err.Error())
	}

	flag.Parse()
}

func showUsage() {
	log.Println("Usage: ")
	log.Println("AddBlock: --host {127.0.0.1:8081} --cmd add --data {something}")
	log.Println("GetBlock: --host {127.0.0.1:8081} --cmd get --height {9}")
	log.Println("Height: --host {127.0.0.1:8081} --cmd height")
}

func checkinghealthy() {
	if len(os.Args) < 2 {
		showUsage()
		os.Exit(1)
	}
}

func main() {
	args := cli.Args{Data: data, Height: height}
	var reply cli.BlockInfo

	err = client.Call(mf[cmd], args, &reply)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Printf("+%v", reply)
}
