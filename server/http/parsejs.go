package http

import (
	"net/http"
	"log"
)

func parseFiles(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if r.Method != "GET" {
		log.Println("error: method should be GET")
		return
	}

	http.ServeFile(w, r, webroot + r.URL.Path[1:])
}
