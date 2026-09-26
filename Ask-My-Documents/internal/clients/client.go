package clients

import (
	"net/http"
	"time"
)

type Clients struct {
	clients *http.Client
}

func NewClient() *Clients {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.MaxIdleConns = 100       // Idle connections across all hosts
	transport.MaxIdleConnsPerHost = 20 // Idle connections kept for each host
	transport.MaxConnsPerHost = 50     // Total connections allowed per host
	transport.IdleConnTimeout = 90 * time.Second

	client := &http.Client{
		Transport: transport,
		Timeout:   10 * time.Second,
	}

	return &Clients{
		clients: client,
	}
}
