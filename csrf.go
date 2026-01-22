package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

func generateCsrf() string {
	csrf := make([]byte, 32)
	if _, err := rand.Read(csrf); err != nil {
		_ = fmt.Errorf("failed to generate csrf: %w", err)
	}
	return base64.StdEncoding.EncodeToString(csrf)
}
