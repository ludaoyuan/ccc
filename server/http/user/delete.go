package user

import (
	"net/http"
	"log"
)

func d_elete(r *http.Request) error {
	r.ParseForm()

	log.Println(r.Form)
	return nil
}
