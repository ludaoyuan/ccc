// balancing

package http

import (
	"net/http"
	"log"

	"server/http/block"
	"server/http/user"
)

var webroot string

func Run(port, wr string) {
	webroot = wr

	http.HandleFunc("/js/", parseFiles)
	http.HandleFunc("/css/", parseFiles)
	http.HandleFunc("/img/", parseFiles)
	http.HandleFunc("/html/", require)

	http.HandleFunc("/", index)
	http.HandleFunc("/block/list", block.Block)
	http.HandleFunc("/user/generate", user.User)

	block.Init()
	user.Init()

	log.Println("http listening...", port)

//	log.Fatalln(http.ListenAndServe("127.0.0.1" + port, nil)) // 正式
	log.Fatalln(http.ListenAndServe(port, nil)) // 测试
}
