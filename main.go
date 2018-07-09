package main

import (
	"log"
	"node"
)

// func init() {
// 	app := cli.NewApp()
// 	app.Name = filepath.Base(os.Args[0])
// 	app.Usage = "command line interface"
//
// 	app.Action = inc
// }

// func inc(ctx *cli.Context) error {
// 	node := node.NewNode(ctx)
// 	node.Start()
// 	return nil
// }

func main() {
	// if err := app.Run(os.Args); err != nil {
	// 	fmt.Fprintln(os.Stderr, err)
	// 	os.Exit(1)
	// }
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	node := node.NewNode()
	node.Start()
}
