// Just a load test script
package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync/atomic"
	"time"

	"golang.org/x/sync/errgroup"
)

const (
	ConcurrentRequestsCount  = 500
	BenchmarkDurationSeconds = 60
	Host                     = "http://127.0.0.1:8080"
	LoginRoute               = "/login"
	GenerateMagicRoute       = "/magic/generate"
)

func generateCsrf() string {
	csrf := make([]byte, 32)
	if _, err := rand.Read(csrf); err != nil {
		_ = fmt.Errorf("failed to generate csrf: %w", err)
	}
	return base64.StdEncoding.EncodeToString(csrf)
}

func main() {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	ctx, cancel := context.WithTimeout(context.Background(), BenchmarkDurationSeconds*time.Second)
	defer cancel()
	group, ctx := errgroup.WithContext(ctx)
	semaphore := make(chan struct{}, ConcurrentRequestsCount)
	var counter uint64

	go func() {
		ticker := time.NewTicker(1 * time.Second)
		for {
			<-ticker.C
			fmt.Println("Concurrent Requests Count: ", counter)
		}
	}()

loop:
	for {
		select {
		case <-ctx.Done():
			break loop
		default:
			semaphore <- struct{}{}
			atomic.AddUint64(&counter, 1)
			group.Go(func() error {
				defer func() {
					atomic.AddUint64(&counter, ^uint64(0))
					<-semaphore
				}()
				req, err := http.NewRequestWithContext(ctx, http.MethodPost, Host+GenerateMagicRoute, nil)
				if err != nil {
					return err
				}
				req.Header.Set("X-CSRF-Token", generateCsrf())
				res, err := client.Do(req)
				if err != nil {
					return err
				}
				if _, err := io.ReadAll(res.Body); err != nil {
					log.Println("error reading:", err)
				}
				defer func() {
					if err := res.Body.Close(); err != nil {
						log.Println("error closing body:", err)
					}
				}()
				return nil
			})
		}
	}
	if err := group.Wait(); err != nil {
		log.Println("request failed:", err)
	}
}
