package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestPing(t *testing.T) {
	r := httptest.NewRequest(
		http.MethodGet,
		"/ping",
		strings.NewReader(""))
	w := httptest.NewRecorder()
	pingHandler(w, r)
	if w.Result().StatusCode != 200 {
		t.Fatal("Status code from /ping was not 200")
	}
}
