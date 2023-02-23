package snmp

import (
	"fmt"
	"github.com/gosnmp/gosnmp"
	"log"
	"time"
)

func Get(ip string, community string, oids []Oid, c chan []OidResult) {
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
	res, err := get(handle, oids)
	c <- res
	return
}

func get(handle gosnmp.GoSNMP, oids []Oid) ([]OidResult, error) {
	next, err := getNext(handle, filter(oids, OidGet))
	if err != nil {
		return nil, err
	}
	count, err := getCount(handle, filter(oids, OidCount))
	if err != nil {
		return nil, err
	}
	res := make([]OidResult, len(oids))
	iNext := 0
	iCount := 0
	for i, oid := range oids {
		var result string
		if oid.Kind == OidGet {
			result = next[iNext]
			iNext++
		} else {
			result = fmt.Sprintf("%d", count[iCount])
			iCount++
		}
		res[i] = OidResult{
			Oid:    oid,
			Result: result,
		}
	}
	return res, nil
}

func getCount(handle gosnmp.GoSNMP, oids []string) ([]int, error) {
	res := make([]int, len(oids))
	for i, oid := range oids {
		results, err := handle.BulkWalkAll(oid)
		if err != nil {
			return nil, err
		}
		res[i] = len(results)
	}
	return res, nil
}

func filter(oids []Oid, count OidKind) []string {
	res := make([]string, 0)
	for _, oid := range oids {
		if oid.Kind == count {
			res = append(res, oid.Oid)
		}
	}
	return res
}

func getNext(handle gosnmp.GoSNMP, oids []string) ([]string, error) {
	results, err := handle.GetNext(oids)
	if err != nil {
		return nil, err
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
	return res, nil
}
