package http

import (
	"net/http"
	"html/template"
	"log"
)

func index(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	log.Println(r.URL.Path[1:])
	if r.URL.Path[1:] == "favicon.ico" {
		return
	}

	if r.Method != "GET" {
		log.Println("error: method should be GET")
		return
	}

	t, err := template.ParseFiles(webroot + "index.html")
	if err != nil {
		log.Println(err.Error())
		return
	}

	err = t.Execute(w, nil)
	if err != nil {
		log.Println(err.Error())
		return
	}
}
