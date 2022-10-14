package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
)

type RequestHandler struct {
	config *Config
}

func (r *RequestHandler) Respond(w http.ResponseWriter, nav string, frequencyTable string) {
	htmlTemplate, err := ioutil.ReadFile("index.html")
	if err != nil {
		panic(fmt.Errorf("could not read index.html: %s", err))
	}
	body := strings.ReplaceAll(string(htmlTemplate), "<nav/>", nav)
	body = strings.ReplaceAll(body, "<frequencytable/>", frequencyTable)
	write(w, body)
}

func (r *RequestHandler) ServeHTTP(w http.ResponseWriter, request *http.Request) {
	fmt.Printf("%s %s\n", request.RemoteAddr, request.URL.String())
	nav, err := getNav()
	if err != nil {
		writeError(w, "could not get nav", err)
		return
	}
	if request.URL.Path == "/" {
		r.Respond(w, nav, "")
	} else {
		siteFilename := "sites/" + path.Base(request.URL.Path)
		siteFile, err := os.Open(siteFilename + ".json")
		if err != nil {
			writeError(w, "could not open site file", err)
			return
		}
		site := &Site{
			config: r.config,
		}
		err = json.NewDecoder(siteFile).Decode(site)
		if err != nil {
			writeError(w, "could not parse site file", err)
			return
		}
		frequencyTable, err := site.getFrequencyTable()
		if err != nil {
			writeError(w, "could not get frequency table", err)
			return
		}
		r.Respond(w, nav, frequencyTable)
	}
}

func getNav() (string, error) {
	files, err := ioutil.ReadDir("sites")
	if err != nil {
		return "", err
	}
	res := ""
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".json") {
			siteName := strings.TrimSuffix(file.Name(), ".json")
			res += fmt.Sprintf("<a href='%s'>%s</a><br/>\n", siteName, siteName)
		}
	}
	return res, nil
}

func write(w http.ResponseWriter, body string) {
	_, err := w.Write([]byte(body))
	if err != nil {
		log.Printf("could not write to client: %s\n", err)
		return
	}
}

func writeError(w http.ResponseWriter, message string, err error) {
	write(w, fmt.Sprintf("%s: %s", message, err))
}
