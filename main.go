package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

func main() {
	start := time.Now()
	const concurrency = 10
	var wg sync.WaitGroup
	wg.Add(concurrency)

	for range concurrency {
		go func() {
			defer wg.Done()
			resp, err := http.Get("http://www.google.com")
			if err != nil {
				panic(err)
			}
			defer resp.Body.Close()
			fmt.Println(resp.Status, time.Since(start))
		}()
	}

	wg.Wait()
	fmt.Println("Done:", time.Since(start))
}
