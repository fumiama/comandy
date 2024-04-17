package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/fumiama/terasu"
)

func TestResolver(t *testing.T) {
	t.Log("canUseIPv6:", canUseIPv6.Get())
	addrs, err := resolver.LookupHost(context.TODO(), "api.mangacopy.com")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(addrs)
	if len(addrs) == 0 {
		t.Fail()
	}
}

func TestDNS(t *testing.T) {
	if canUseIPv6.Get() {
		dotv6servers.test()
	}
	dotv4servers.test()
	for i := 0; i < 100; i++ {
		addrs, err := resolver.LookupHost(context.TODO(), "api.mangacopy.com")
		if err != nil {
			t.Fatal(err)
		}
		t.Log(addrs)
		if len(addrs) == 0 {
			t.Fail()
		}
		time.Sleep(time.Millisecond * 50)
	}
}

func TestBadDNS(t *testing.T) {
	dotv6serversbak := dotv6servers.m
	dotv4serversbak := dotv4servers.m
	defer func() {
		dotv6servers.m = dotv6serversbak
		dotv4servers.m = dotv4serversbak
	}()
	if canUseIPv6.Get() {
		dotv6servers = dnsservers{
			m: map[string][]*dnsstat{},
		}
		dotv6servers.add(map[string][]string{"test.bad.host": {"169.254.122.111"}})
	} else {
		dotv4servers = dnsservers{
			m: map[string][]*dnsstat{},
		}
		dotv4servers.add(map[string][]string{"test.bad.host": {"169.254.122.111:853"}})
	}
	for i := 0; i < 10; i++ {
		addrs, err := resolver.LookupHost(context.TODO(), "api.mangacopy.com")
		t.Log(err)
		if err == nil && len(addrs) > 0 {
			t.Fatal("unexpected")
		}
		time.Sleep(time.Millisecond * 50)
	}
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
