package http

import (
	"net/http"
	"html/template"
	"log"
)

func require(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if r.Method != "GET" {
		log.Println("error: method should be GET")
		return
	}

	log.Println(r.URL.Path[1:])
	t, err := template.ParseFiles(webroot + r.URL.Path[1:])
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
