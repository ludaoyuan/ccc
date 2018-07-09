package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/rpc/v2"
	"github.com/gorilla/rpc/v2/json2"
	"github.com/gorilla/websocket"
	"github.com/urfave/cli"
)

type Client struct{}

func (c *Client) HandleParam(r *http.Request, args *int, reply *int) error {
	*reply = 8
	return nil
}

type MyStrings struct {
	S []string
}

func (c *Client) HandleParam2(r *http.Request, args *int, reply *MyStrings) error {
	s := []string{"001", "002"}
	*reply = MyStrings{S: s}
	return nil
}

// curl -X POST -H "Content-Type:application/json" -d '{"jsonrpc":"2.0", "method":"Client.HandleParam3","params":{"S":["001", "002"]},"id":3}' http://localhost:1234
func (c *Client) HandleParam3(r *http.Request, args *MyStrings, reply *MyStrings) error {
	*reply = *args
	return nil
}

func (c *Client) HandleParam4(r *http.Request, args []string, reply *string) error {
	*reply = "sanghai"
	return nil
}

func (c *Client) HandleParam5(r *http.Request, args *string, reply *string) error {
	*reply = "sanghai"
	return nil
}

func HandleRPC(ctx *cli.Context) error {
	s := rpc.NewServer()
	s.RegisterCodec(json2.NewCodec(), "application/json")
	s.RegisterService(new(Client), "")

	http.Handle("/", s)
	// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// })

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ws, err := websocket.Upgrade(w, r, nil, 0, 0)
		if err != nil {
			log.Println(err.Error())
			return
		}
		defer ws.Close()

		for {
			mt, message, err := ws.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				break
			}
			log.Printf("recv: %s", message)
			err = ws.WriteMessage(mt, message)
			if err != nil {
				log.Println("write:", err)
				break
			}
		}
	})

	log.Fatal(http.ListenAndServe(":1234", nil))
	return nil
}

func main() {
	app := cli.NewApp()
	app.Name = "TestRPCSvr"
	app.Action = HandleRPC

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
