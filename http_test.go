package main

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"testing"
)

func TestRequest(t *testing.T) {
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
