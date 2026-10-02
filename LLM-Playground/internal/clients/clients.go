package clients

import (
	"llm-playground/internal/constants"
	"net/http"
)

type Clients struct {
	clients *http.Client
}

func NewClient() *Clients {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.MaxIdleConns = constants.ClientMaxIdleConns
	transport.MaxIdleConnsPerHost = constants.ClientMaxIdleConnsPerHost
	transport.MaxConnsPerHost = constants.ClientMaxConnsPerHost
	transport.IdleConnTimeout = constants.ClientIdleConnTimeout

	client := &http.Client{
		Transport: transport,
		Timeout:   constants.ClientTimeout,
	}

	return &Clients{
		clients: client,
	}
}
