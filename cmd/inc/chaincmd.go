package main

import "github.com/urfave/cli"

var (
	initCommand = cli.Command{
		Action: kslfj,
		Name:   "init",
	}
)
