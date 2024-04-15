package main

import (
	"net/http"

	"github.com/RomiChan/syncx"
)

var canUseIPv6 = syncx.Lazy[bool]{Init: func() bool {
	resp, err := http.Get("http://v6.ipv6-test.com/json/widgetdata.php?callback=?")
	if err != nil {
		return false
	}
	_ = resp.Body.Close()
	return true
}}
