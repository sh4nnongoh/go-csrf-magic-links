package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	Host               = "http://127.0.0.1:8080"
	LoginRoute         = "/login"
	GenerateMagicRoute = "/magic/generate"
)

func BenchmarkGin(b *testing.B) {
	router := NewRouter(true)
	req, err := http.NewRequest(http.MethodPost, Host+GenerateMagicRoute, nil)
	req.Header.Set("X-CSRF-Token", generateCsrf())
	if err != nil {
		b.Error(err)
		return
	}
	b.ResetTimer()
	for b.Loop() {
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
	}
}
