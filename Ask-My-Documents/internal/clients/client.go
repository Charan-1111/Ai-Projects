package clients

import (
	"net/http"

	"ask-my-documents/internal/constants"
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
