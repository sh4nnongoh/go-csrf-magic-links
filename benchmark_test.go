package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"go.uber.org/zap"
)

const (
	Host               = "http://127.0.0.1:8080"
	GenerateMagicRoute = "/magic/generate"
)

func BenchmarkGin(b *testing.B) {
	logger, err := zap.NewProduction()
	if err != nil {
		b.Error(err)
		return
	}
	defer logger.Sync()
	router := NewRouter(logger, true)
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
