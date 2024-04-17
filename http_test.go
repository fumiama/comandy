package main

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

func TestClientGet(t *testing.T) {
	_ = canUseIPv6.Get()
	req, err := http.NewRequest("GET", "https://api.mangacopy.com/api/v3/h5/homeIndex?platform=3", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("user-agent", "COPY/2.1.7")
	resp, err := (*http.Client)(&cli).Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	t.Log("[T] response code", resp.StatusCode)
	for k, vs := range resp.Header {
		for _, v := range vs {
			t.Log("[T] response header", k+":", v)
		}
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if len(data) == 0 {
		t.Fail()
	}
	t.Log(bytesToString(data))
}

func TestRequest(t *testing.T) {
	r := cli.request(`{"code":0,"headers":{"authorization":"Token ","host":"api.mangacopy.com","source":"copyApp","webp":"1","region":"1","version":"2.1.7","platform":"3","user-agent":"COPY/2.1.7"},"method":"GET","url":"https://api.mangacopy.com/api/v3/h5/homeIndex?platform\u003d3"}`)
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
