package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	config, err := loadConfig("config.json")
	if err != nil {
		panic(err)
	}
	htmlTemplate, err := ioutil.ReadFile("index.html")
	if err != nil {
		panic(fmt.Errorf("could not read index.html: %s", err))
	}
	handler := &RequestHandler{string(htmlTemplate), config}
	fmt.Printf("webserver listening on %s\n", config.Listen)
	err = http.ListenAndServe(config.Listen, handler)
	if err != nil {
		panic(fmt.Errorf("could not listen: %err", err))
	}
}
