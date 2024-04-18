package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"

	"github.com/fumiama/terasu/http2"
)

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

func gorequest(para string) (ret string) {
	r := capsule{}
	defer func() {
		err := recover()
		if err != nil {
			ret = r.printstrerr(fmt.Sprint())
		}
	}()
	err := json.Unmarshal(stringToBytes(para), &r)
	if err != nil {
		return r.printerr(err)
	}
	if r.U == "" || !strings.HasPrefix(r.U, "https://") {
		return r.printstrerr("invalid url '" + r.U + "'")
	}
	if r.M != "GET" && r.M != "POST" && r.M != "DELETE" {
		return r.printstrerr("invalid method '" + r.U + "'")
	}
	var body io.Reader
	if len(r.D) > 0 {
		body = strings.NewReader(r.D)
	}
	req, err := http.NewRequest(r.M, r.U, body)
	if err != nil {
		return r.printerr(err)
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
			return r.printstrerr("unsupported H type " + reflect.ValueOf(x).Type().Name())
		}
	}
	resp, err := http2.DefaultClient.Do(req)
	if err != nil {
		return r.printerr(err)
	}
	defer resp.Body.Close()
	sb := strings.Builder{}
	enc := base64.NewEncoder(base64.StdEncoding, &sb)
	_, err = io.Copy(enc, resp.Body)
	_ = enc.Close()
	if err != nil {
		return r.printerr(err)
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
		return r.printerr(err)
	}
	return outbuf.String()
}
