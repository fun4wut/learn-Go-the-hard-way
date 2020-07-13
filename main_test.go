package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

var s = NewServer()
var recorder = httptest.NewRecorder()

func TestTiny(t *testing.T) {
	recorder.Flush()
	req, _ := http.NewRequest("GET", "http://localhost:3000/hello?name=foo", nil)
	s.Get("/hello", func(ctx *Context) string {
		name := ctx.URL.Query().Get("name")
		return name
	})
	s.ServeHTTP(recorder, req)
	name, err := ioutil.ReadAll(recorder.Body)
	if err != nil {
		t.Fatal(err)
	}
	if string(name) != "foo" {
		println(string(name))
		t.Fail()
		return
	}
}

func TestTiny2(t *testing.T) {
	recorder.Flush()
	req, _ := http.NewRequest("GET", "http://localhost:3000/rua", nil)
	s.Get("/rua", func() string {
		return "2333"
	})
	s.ServeHTTP(recorder, req)
	val, err := ioutil.ReadAll(recorder.Body)
	if err != nil {
		t.Fatal(err)
	}
	if string(val) != "2333" {
		println(string(val))
		t.Fail()
		return
	}
}
