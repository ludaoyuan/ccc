package block

import (
	"net/http"
	"log"
)

func put(r *http.Request) error {
	r.ParseForm()

	log.Println(r.Form)
	return nil
}
