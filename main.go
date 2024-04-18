package main

import "C"

import (
	"encoding/json"

	"github.com/fumiama/terasu/dns"
	"github.com/fumiama/terasu/ip"
)

func main() {}

// para: json of map[host string][]addr:port string
//
//export add_dns
func add_dns(para *C.char, is_ipv6 C.int) *C.char {
	m := map[string][]string{}
	err := json.Unmarshal(stringToBytes(C.GoString(para)), &m)
	if err != nil {
		return C.CString(err.Error())
	}
	if is_ipv6 != 0 {
		if !ip.IsIPv6Available.Get() {
			return C.CString("cannot use ipv6")
		}
		dns.IPv6Servers.Add(m)
		return nil
	}
	dns.IPv4Servers.Add(m)
	return nil
}

// para:
//
//	request("{\"method\":\"GET\","
//		"\"url\":\"https://i.pximg.net/img-master/img/2012/04/04/21/24/46/26339586_p0_master1200.jpg\","
//		"\"headers\":{"
//			"\"Referer\":\"https://www.pixiv.net/\","
//			"\"User-Agent\":\"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36 Edg/123.0.0.0\""
//		"}"
//	"}");
//
//export request
func request(para *C.char) *C.char {
	return C.CString(gorequest(C.GoString(para)))
}
