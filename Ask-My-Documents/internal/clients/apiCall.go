package clients

import (
	"bytes"
	"context"
	"net/http"
	"strings"

	"github.com/bytedance/sonic"
)

type ApiInterface interface {
	ApiCall(ctx context.Context)
}

type OutboundCall struct {
	url         string
	method      string
	queryParams map[string]string
	bodyParams  map[string]any
	headers     map[string]string
	clients     *Clients
}

func NewOutboundCall(url string, method string, queryParams map[string]string, bodyParams map[string]any, headers map[string]string, httpClients *Clients) *OutboundCall {
	return &OutboundCall{
		url:         url,
		method:      method,
		queryParams: queryParams,
		bodyParams:  bodyParams,
		headers:     headers,
		clients:     httpClients,
	}
}

func buildUrl(url string, queryParams map[string]string) string {
	var sb strings.Builder

	sb.WriteString(url)
	sb.WriteString("?")
	for key, value := range queryParams {
		sb.WriteString(key)
		sb.WriteString("=")
		sb.WriteString(value)
		sb.WriteString("&")
	}

	return sb.String()
}

func (ob *OutboundCall) ApiCall(ctx context.Context) {
	finalUrl := buildUrl(ob.url, ob.queryParams)

	bodyBytes, err := sonic.Marshal(ob.bodyParams)
	if err != nil {

	}

	req, err := http.NewRequestWithContext(ctx, ob.method, finalUrl, bytes.NewBuffer(bodyBytes))
	if err != nil {

	}

	// setting headers
	for key, value := range ob.headers {
		req.Header.Set(key, value)
	}

	resp, err := ob.clients.clients.Do(req)
	if err != nil {

	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {

	}

}
