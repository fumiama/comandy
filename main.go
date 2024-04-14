package main

import "C"

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/fumiama/terasu"
)

func main() {}

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

//export request
func request(para *C.char) *C.char {
	r := capsule{}
	err := json.Unmarshal(stringToBytes(C.GoString(para)), &r)
	if err != nil {
		return C.CString(r.printerr(err))
	}
	if r.U == "" || !strings.HasPrefix(r.U, "http") {
		return C.CString(r.printstrerr("invalid url '" + r.U + "'"))
	}
	if r.M != "GET" && r.M != "POST" && r.M != "DELETE" {
		return C.CString(r.printstrerr("invalid method '" + r.U + "'"))
	}
	req, err := http.NewRequest(r.M, r.U, strings.NewReader(r.D))
	if err != nil {
		return C.CString(r.printerr(err))
	}
	for k, vs := range r.H {
		lk := strings.ToLower(k)
		if strings.HasPrefix(lk, "x-") {
			continue
		}
		switch x := vs.(type) {
		case string:
			req.Header.Add(k, x)
		case []string:
			for _, v := range x {
				req.Header.Add(k, v)
			}
		default:
			return C.CString(r.printstrerr("unsupported H type " + reflect.ValueOf(x).Type().Name()))
		}
	}
	resp, err := cli.Do(req)
	if err != nil {
		return C.CString(r.printerr(err))
	}
	defer resp.Body.Close()
	sb := strings.Builder{}
	enc := base64.NewEncoder(base64.StdEncoding, &sb)
	_, err = io.CopyN(enc, resp.Body, resp.ContentLength)
	_ = enc.Close()
	if err != nil {
		return C.CString(r.printerr(err))
	}
	r.C = resp.StatusCode
	r.H = make(map[string]any, len(resp.Header)*2)
	for k, vs := range resp.Header {
		if len(vs) == 1 {
			r.H[k] = vs[0]
			continue
		}
		r.H[k] = vs
	}
	r.D = sb.String()
	outbuf := strings.Builder{}
	err = json.NewEncoder(&outbuf).Encode(&r)
	if err != nil {
		return C.CString(r.printerr(err))
	}
	return C.CString(outbuf.String())
}
