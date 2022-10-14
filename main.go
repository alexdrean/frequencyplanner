package main

import (
	"fmt"
	"net/http"
)

func main() {
	config, err := loadConfig("config.json")
	if err != nil {
		panic(err)
	}
	handler := &RequestHandler{config}
	fmt.Printf("webserver listening on %s\n", config.Listen)
	err = http.ListenAndServe(config.Listen, handler)
	if err != nil {
		panic(fmt.Errorf("could not listen: %err", err))
	}
}
