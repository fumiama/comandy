package main

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"sync"
	"testing"
	"unsafe"

	"github.com/fumiama/terasu/dns"
)

type testlist struct {
	sync.RWMutex
	m map[string][]*uintptr
}

func TestRequest(t *testing.T) {
	(*testlist)(unsafe.Pointer(&dns.IPv4Servers)).m = make(map[string][]*uintptr)
	(*testlist)(unsafe.Pointer(&dns.IPv6Servers)).m = make(map[string][]*uintptr)
	dns.IPv4Servers.Add(&dns.DNSConfig{
		Servers: map[string][]string{
			"dot.360.cn": {
				"101.198.192.33:853",
				"112.65.69.15:853",
				"101.226.4.6:853",
				"218.30.118.6:853",
				"123.125.81.6:853",
				"140.207.198.6:853",
			},
		},
	})
	dns.IPv6Servers.Add(&dns.DNSConfig{
		Servers: map[string][]string{
			"dot.360.cn": {
				"101.198.192.33:853",
				"112.65.69.15:853",
				"101.226.4.6:853",
				"218.30.118.6:853",
				"123.125.81.6:853",
				"140.207.198.6:853",
			},
		},
	})
	r := gorequest(`{"code":0,"headers":{"authorization":"Token ","host":"api.mangacopy.com","source":"copyApp","webp":"1","region":"1","version":"2.1.7","platform":"3","user-agent":"COPY/2.1.7"},"method":"GET","url":"https://api.mangacopy.com/api/v3/h5/homeIndex?platform\u003d3"}`)
	t.Log(r)
	c := capsule{}
	err := json.Unmarshal(stringToBytes(r), &c)
	if err != nil {
		t.Fatal(err)
	}
	if c.C != http.StatusOK {
		s, err := base64.StdEncoding.DecodeString(c.D)
		if err != nil {
			t.Fatal("status code", c.C, "msg:", c.D)
		} else {
			t.Fatal("status code", c.C, "msg:", s)
		}
	}
	if len(c.D) == 0 {
		t.Fatal("empty data")
	}
}
