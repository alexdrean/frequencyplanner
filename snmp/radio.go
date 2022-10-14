package snmp

import (
	"fmt"
	"github.com/gosnmp/gosnmp"
	"log"
	"time"
)

func GetNext(ip string, community string, oids []string, c chan []string) {
	handle := gosnmp.GoSNMP{
		Target:    ip,
		Port:      161,
		Community: community,
		Version:   gosnmp.Version2c,
		Timeout:   2 * time.Second,
		Retries:   5,
	}
	err := handle.Connect()
	if err != nil {
		log.Println(ip, err)
		c <- nil
		return
	}
	defer handle.Conn.Close()
	results, err := handle.GetNext(oids)
	if err != nil {
		log.Println(ip, err)
		c <- nil
		return
	}
	res := make([]string, len(results.Variables))
	for i, result := range results.Variables {
		if r, ok := result.Value.(int); ok {
			res[i] = fmt.Sprintf("%d", r)
		} else if r, ok := result.Value.(string); ok {
			res[i] = r
		} else if r, ok := result.Value.([]byte); ok {
			res[i] = string(r)
		} else {
			res[i] = fmt.Sprintf("%v", result.Value)
		}
	}
	c <- res
	return
}
