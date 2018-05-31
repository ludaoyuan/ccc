package main

import (
	"conf"
	"log"
	"server/tcp"
	"server/http"
)

func main() {
	config := make(conf.Config)
	config.Init("/etc/inc/inc.conf")

	log.SetFlags(log.Lshortfile | log.LstdFlags)

	go tcp.Run(config["TCP"]["domain"], config["TCP"]["port"])
	http.Run(config["HTTP"]["port"], config["HTTP"]["webroot"])
}
