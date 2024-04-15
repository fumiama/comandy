package main

import (
	"context"
	"crypto/tls"
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
		dotv6servers.test(t)
	}
	dotv4servers.test(t)
}

func (ds *dnsservers) test(t *testing.T) {
	ds.RLock()
	defer ds.RUnlock()
	for host, addrs := range ds.m {
		for _, addr := range addrs {
			if !addr.E {
				continue
			}
			conn, err := net.Dial("tcp", addr.A)
			if err != nil {
				continue
			}
			tlsConn := terasu.Use(tls.Client(conn, &tls.Config{ServerName: host}))
			err = tlsConn.Handshake()
			_ = tlsConn.Close()
			if err == nil {
				t.Log("succ:", host, addr.A)
				continue
			}
			t.Fatal("fail:", host, addr.A)
		}
	}
}
