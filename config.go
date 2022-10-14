package main

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
)

type Site struct {
	config    *Config
	Radios    []Radio `json:"radios"`
	Community string  `json:"community"`
}

type Radio struct {
	Ip       string `json:"ip"`
	Platform string `json:"platform"`
}

type Config struct {
	Listen    string               `json:"listen"`
	Platforms map[string]*Platform `json:"platforms"`
}
type Platform struct {
	Header  string `json:"header"`
	Content string `json:"content"`
	oids    []string
}

func (p *Platform) parseOIDs() ([]string, error) {
	r, err := regexp.Compile("oid:[0-9.]+")
	if err != nil {
		return nil, fmt.Errorf("regexp returned: %s", err)
	}
	oids := r.FindAllString(p.Header+"\n"+p.Content, -1)
	for i, oid := range oids {
		oids[i] = oid[4:]
	}
	return oids, nil
}

func loadConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("could not open config: %s", err)
	}
	config := &Config{}
	err = json.NewDecoder(file).Decode(config)
	if err != nil {
		return nil, fmt.Errorf("could not parse config: %s", err)
	}
	for name, platform := range config.Platforms {
		platform.oids, err = platform.parseOIDs()
		if err != nil {
			return nil, fmt.Errorf("could not parse platform %s: %s", name, err)
		}
	}
	return config, nil
}
