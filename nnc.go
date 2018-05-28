package main

import (
	"conf"
	"server/tcp"
	"server/http"
)

func main() {
	config := make(conf.Config)
	config.Init("/etc/nnc/nnc.conf")

	go tcp.Run(config["TCP"]["domain"], config["TCP"]["port"])
	http.Run(config["HTTP"]["port"], config["HTTP"]["webroot"])
}
