package user

import (
	"net/http"
	"log"
)

var funcs [4]map[string]func(form map[string][]string) error

func Init() {
	// 0 get; 1 post; 2 put; 3 delete
	funcs[0] = make(map[string]func(form map[string][]string) error)
	funcs[1] = make(map[string]func(form map[string][]string) error)
	funcs[2] = make(map[string]func(form map[string][]string) error)
	funcs[3] = make(map[string]func(form map[string][]string) error)

	funcs[0]["/user/generate"] = generate
}

func User(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var err error

	switch r.Method {
	case "GET":
		r.ParseForm()

		log.Println(r.URL.Path, funcs[0])
		err = get(r.Form, funcs[0][r.URL.Path])
		if err != nil {
			log.Println(err.Error())
			return
		}
	case "POST":
		err = post(r)
		if err != nil {
			log.Println(err.Error())
			return
		}
	case "PUT":
		err = put(r)
		if err != nil {
			log.Println(err.Error())
			return
		}
	case "DELETE":
		err = d_elete(r)
		if err != nil {
			log.Println(err.Error())
			return
		}
	default:
		w.Write([]byte(`{"code":404,"msg":"Page not found"}`))
	}
}
