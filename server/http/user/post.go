package user

import (
	"net/http"
	"log"
)

func post(r *http.Request) error {
	r.ParseForm()

	log.Println(r.Form)
	return nil
}
