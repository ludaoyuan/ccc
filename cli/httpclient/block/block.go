package httpblock

import (
	"cli"
	"cli/conf"
	"fmt"
	"log"
	"net/http"
	"net/rpc"
	"strconv"
)

const (
	CfgPath = "./cnf.conf"
)

var (
	client *rpc.Client
	cfg    conf.Config
	err    error
	host   string
)

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	cfg = make(conf.Config)
	cfg.Init(CfgPath)

	client, err = rpc.Dial("tcp", cfg["HTTP"]["host"])
	if err != nil {
		log.Fatal(err.Error())
	}
}

func HandleAddBlock(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if r.URL.Path[1:] == "favicon.ico" {
		return
	}

	// switch r.Method {
	// case "get":
	// case "post":
	// }

	args := cli.Args{Data: r.FormValue("data")}
	var reply cli.BlockInfo
	err = client.Call("CLI.AddBlock", args, &reply)
	if err != nil {
		log.Println(err.Error())
		w.Write([]byte(err.Error()))
		return
	}
	w.Write([]byte("Success"))
}

func HandleGetBlock(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if r.URL.Path[1:] == "favicon.ico" {
		return
	}

	h, err := strconv.ParseInt(r.FormValue("height"), 10, 64)
	log.Println(r.FormValue("height"))
	if err != nil {
		log.Println(err.Error())
		w.Write([]byte(err.Error()))
		return
	}
	args := cli.Args{Height: h}
	var reply cli.BlockInfo
	err = client.Call("CLI.GetBlock", args, &reply)
	if err != nil {
		log.Println(err.Error())
		w.Write([]byte(err.Error()))
		return
	}
	w.Write([]byte(reply.Data))
}

func HandleHeight(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if r.URL.Path[1:] == "favicon.ico" {
		return
	}

	var reply cli.BlockInfo
	err = client.Call("CLI.Height", cli.Args{}, &reply)
	if err != nil {
		log.Println(err.Error())
		w.Write([]byte(err.Error()))
		return
	}

	w.Write([]byte(fmt.Sprintf("%d", reply.Height)))
}
