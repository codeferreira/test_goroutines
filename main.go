package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"time"
)

func main() {
	start := time.Now()
	const concurrency = 10
	var wg sync.WaitGroup
	wg.Add(concurrency)

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(10 * time.Second)
		fmt.Fprintln(w, "Hello, client")
	}))

	for range concurrency {
		go func(ctx context.Context) {
			defer wg.Done()

			req, err := http.NewRequestWithContext(ctx, "GET", server.URL, nil)
			if err != nil {
				panic(err)
			}

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				if errors.Is(err, context.DeadlineExceeded) {
					fmt.Println("Deadline exceeded")
					return
				}
				panic(err)
			}
			defer resp.Body.Close()
		}(ctx)
	}

	wg.Wait()
	fmt.Println("Done:", time.Since(start))
}
