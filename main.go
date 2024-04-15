package main

import "C"

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"reflect"
	"strings"
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
		if !canUseIPv6.Get() {
			return C.CString("cannot use ipv6")
		}
		dotv6servers.add(m)
		return nil
	}
	dotv4servers.add(m)
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
	r := capsule{}
	err := json.Unmarshal(stringToBytes(C.GoString(para)), &r)
	if err != nil {
		return C.CString(r.printerr(err))
	}
	if r.U == "" || !strings.HasPrefix(r.U, "https://") {
		return C.CString(r.printstrerr("invalid url '" + r.U + "'"))
	}
	if r.M != "GET" && r.M != "POST" && r.M != "DELETE" {
		return C.CString(r.printstrerr("invalid method '" + r.U + "'"))
	}
	var body io.Reader
	if len(r.D) > 0 {
		body = strings.NewReader(r.D)
	}
	req, err := http.NewRequest(r.M, r.U, body)
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
