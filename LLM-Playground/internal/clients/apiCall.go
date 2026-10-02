package clients

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type ApiInterface interface {
	ApiCall(ctx context.Context) ([]byte, error)
}

type OutboundCall struct {
	url         string
	method      string
	queryParams map[string]string
	bodyParams  []byte
	headers     map[string]string
	clients     *Clients
}

func NewOutboundCall(url string, method string, queryParams map[string]string, bodyParams []byte, headers map[string]string, clients *Clients) *OutboundCall {
	return &OutboundCall{
		url:         url,
		method:      method,
		queryParams: queryParams,
		bodyParams:  bodyParams,
		headers:     headers,
		clients:     clients,
	}
}

func buildUrl(url string, queryParams map[string]string) string {
	if len(queryParams) == 0 {
		return url
	}

	var sb strings.Builder
	sb.WriteString(url)
	if strings.ContainsRune(url, '?') == false {
		sb.WriteString("?")
	}
	for key, value := range queryParams {
		sb.WriteString(key)
		sb.WriteString("=")
		sb.WriteString(value)
		sb.WriteString("&")
	}

	return strings.TrimRight(sb.String(), "&")
}

func (ob *OutboundCall) ApiCall(ctx context.Context) ([]byte, error) {
	if ob == nil || ob.clients == nil || ob.clients.clients == nil {
		return nil, fmt.Errorf("invalid outbound call configuration")
	}

	finalURL := buildUrl(ob.url, ob.queryParams)

	var bodyReader io.Reader = http.NoBody
	if len(ob.bodyParams) > 0 {
		// bodyBytes, err := sonic.Marshal(ob.bodyParams)
		// if err != nil {
		// 	return nil, err
		// }
		bodyReader = bytes.NewBuffer(ob.bodyParams)
	}

	req, err := http.NewRequestWithContext(ctx, ob.method, finalURL, bodyReader)
	if err != nil {
		return nil, err
	}

	for key, value := range ob.headers {
		req.Header.Set(key, value)
	}
	if req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := ob.clients.clients.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		fmt.Println("response body :", string(responseBody))
		return responseBody, fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(responseBody))
	}

	return responseBody, nil
}
