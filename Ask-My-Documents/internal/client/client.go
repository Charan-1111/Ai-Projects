package client

import "net/http"

type Client struct {
	client *http.Client
}

func NewClient() {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.MaxIdleConns = 100          // Idle connections across all hosts
	transport.MaxIdleConnsPerHost = 20     // Idle connections kept for each host
	transport.MaxConnsPerHost = 50         // Total connections allowed per host
	transport.IdleConnTimeout = 90 * time.Second

	client := &http.Client{
		Transport: transport,
		Timeout:   10 * time.Second,
	}

	return Client{
		client: client,
	}
}
