package main

import (
	"encoding/json"
	"fmt"
	"frequencyplanner/snmp"
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
	oids    []snmp.Oid
}

func (p *Platform) parseOIDs() ([]snmp.Oid, error) {
	r, err := regexp.Compile("(oid|oidcount):([0-9.]+)")
	if err != nil {
		return nil, fmt.Errorf("regexp returned: %s", err)
	}
	oids := r.FindAllStringSubmatch(p.Header+"\n"+p.Content, -1)
	res := make([]snmp.Oid, len(oids))
	for i, o := range oids {
		var kind snmp.OidKind
		if o[1] == "oid" {
			kind = snmp.OidGet
		} else {
			kind = snmp.OidCount
		}
		res[i] = snmp.Oid{Oid: o[2], Kind: kind}
	}
	return res, nil
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
