package main

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/fumiama/terasu"
)

var dialer = net.Dialer{
	Timeout: time.Minute,
}

var cli = http.Client{
	Transport: &http.Transport{
		DialTLSContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			conn, err := dialer.DialContext(ctx, "tcp", addr)
			if err != nil {
				return nil, err
			}
			host, _, err := net.SplitHostPort(addr)
			if err != nil {
				return nil, err
			}
			return terasu.Use(tls.Client(conn, &tls.Config{
				ServerName:         host,
				InsecureSkipVerify: true,
			})), nil
		},
	},
}

type capsule struct {
	C int            `json:"code,omitempty"`
	M string         `json:"method,omitempty"`
	U string         `json:"url,omitempty"`
	H map[string]any `json:"headers,omitempty"`
	D string         `json:"data,omitempty"`
}

func (r *capsule) printerr(err error) string {
	buf := strings.Builder{}
	r.C = http.StatusInternalServerError
	r.D = base64.StdEncoding.EncodeToString(stringToBytes(err.Error()))
	_ = json.NewEncoder(&buf).Encode(r)
	return buf.String()
}

func (r *capsule) printstrerr(err string) string {
	buf := strings.Builder{}
	r.C = http.StatusInternalServerError
	r.D = base64.StdEncoding.EncodeToString(stringToBytes(err))
	_ = json.NewEncoder(&buf).Encode(r)
	return buf.String()
}
