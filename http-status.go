package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"
)

// GOOS=darwin GOARCH=arm64 go build -o http-status-macos http-status.go
// GOOS=windows GOARCH=amd64 go build -o http-status-win.exe http-status.go
func main() {
	jsonFile := flag.String("file", "", "Path to the JSON file containing the https URL endpoints (required)")
	ip := flag.String("ip", "", "IP address to resolve the domain to (optional)")
	flag.Parse()

	if *jsonFile == "" {
		flag.Usage()
		os.Exit(1)
	}

	data, err := ioutil.ReadFile(*jsonFile)
	if err != nil {
		fmt.Println("Error reading JSON file:", err)
		os.Exit(1)
	}

	var endpoints []string
	if err := json.Unmarshal(data, &endpoints); err != nil {
		fmt.Println("Error parsing JSON:", err)
		os.Exit(1)
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // Allow insecure requests
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				if *ip != "" {
					addr = *ip + ":443"
				}

				dialer := &net.Dialer{
					Timeout: 5 * time.Second,
				}

				return dialer.DialContext(ctx, network, addr)
			},
		},
	}

	var wg sync.WaitGroup
	mu := &sync.Mutex{}
	for _, endpoint := range endpoints {
		wg.Add(1)
		go func(endpoint string) {
			defer wg.Done()

			u, err := url.Parse(endpoint)
			if err != nil {
				fmt.Println("Error parsing URL:", err)
				return
			}

			resp, err := client.Get(u.String())
			if err != nil {
				fmt.Println("Error making request:", err)
				return
			}
			defer resp.Body.Close()

			mu.Lock()
			fmt.Printf("%d %s\n", resp.StatusCode, u.Path)
			mu.Unlock()
		}(endpoint)
	}

	wg.Wait()
}
