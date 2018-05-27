package user

import (
	"log"
)

func generate(form map[string][]string) error {
	log.Println(form)
	return nil
}

func get(form map[string][]string, callback func(map[string][]string) error) error {
	err := callback(form)
	if err != nil {
		return err
	}

	return nil
}
