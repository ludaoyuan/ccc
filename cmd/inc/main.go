package main

import (
	"node"
	"os"
	"path/filepath"

	"github.com/urfave/cli"
)

func init() {
	app := cli.NewApp()
	app.Name = filepath.Base(os.Args[0])
	app.Usage = "command line interface"

	app.Action = inc
	app.Commands = []cli.Command{
		initCommand,
		chainCommand,
		blockCommand,
		walletCommand,
		transactionCommand,
	}
}

func main() {
	node := node.NewNode()
	node.Start()
}
