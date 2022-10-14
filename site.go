package main

import (
	"frequencyplanner/snmp"
	"html"
	"strings"
)

func (s *Site) getFrequencyData() (map[string][]string, []string, error) {
	channels := make([]chan []string, len(s.Radios))
	var errs []string
	for i, radio := range s.Radios {
		channels[i] = make(chan []string)
		platform := s.config.Platforms[radio.Platform]
		if platform == nil {
			errs = append(errs, "platform '"+radio.Platform+"' not defined in config.json")
			continue
		}
		go snmp.GetNext(radio.Ip, s.Community, platform.oids, channels[i])
	}
	results := map[string][]string{}
	for i, channel := range channels {
		res := <-channel
		if res == nil {
			errs = append(errs, "could not get snmp data for "+s.Radios[i].Ip+" (see logs for details)")
			continue
		}
		platform := s.config.Platforms[s.Radios[i].Platform]
		header := platform.Header
		content := platform.Content
		for i, oid := range platform.oids {
			header = strings.ReplaceAll(header, "oid:"+oid, res[i])
			content = strings.ReplaceAll(content, "oid:"+oid, res[i])
		}
		if results[header] == nil {
			results[header] = []string{}
		}
		results[header] = append(results[header], content)
	}
	return results, errs, nil
}

func (s *Site) getFrequencyTable() (string, error) {
	results, errs, err := s.getFrequencyData()
	if err != nil {
		return "", err
	}
	var headers []string
	mostRadios := 0
	for k, v := range results {
		headers = append(headers, k)
		if len(v) > mostRadios {
			mostRadios = len(v)
		}
	}
	sortSlice(headers)
	body := ""
	if len(errs) > 0 {
		body += "<b>Error occured:</b><br>\n"
	}
	for _, err := range errs {
		body += err + "<br>\n"
	}
	body += "<table>\n<tr>\n"
	for _, header := range headers {
		body += "\t<th>" + html.EscapeString(header) + "</th>\n"
	}
	body += "</tr>\n"
	for i := 0; i < mostRadios; i++ {
		row := "<tr>"
		for _, header := range headers {
			contents := results[header]
			row += "<td>"
			if len(contents) > i {
				row += contents[i]
			}
			row += "</td>"
		}
		body += row + "\n"
	}
	body += "</table>"
	return body, nil
}
