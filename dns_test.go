package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"testing"

	"github.com/fumiama/terasu"
)

func TestResolver(t *testing.T) {
	t.Log("canUseIPv6:", canUseIPv6.Get())
	addrs, err := resolver.LookupHost(context.TODO(), "dns.google")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(addrs)
}

func TestDNS(t *testing.T) {
	if canUseIPv6.Get() {
		dotv6servers.test()
	}
	dotv4servers.test()
}

func (ds *dnsservers) test() {
	ds.RLock()
	defer ds.RUnlock()
	for host, addrs := range ds.m {
		for _, addr := range addrs {
			if !addr.E {
				continue
			}
			fmt.Println("dial:", host, addr.A)
			conn, err := net.Dial("tcp", addr.A)
			if err != nil {
				continue
			}
			tlsConn := tls.Client(conn, &tls.Config{ServerName: host})
			err = terasu.Use(tlsConn).Handshake()
			_ = tlsConn.Close()
			if err == nil {
				fmt.Println("succ:", host, addr.A)
				continue
			}
			fmt.Println("fail:", host, addr.A)
		}
	}
}
