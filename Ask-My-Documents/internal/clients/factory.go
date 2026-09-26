package clients

type ApiFactory interface {
	Create(url string, method string, queryParams map[string]string, bodyParams map[string]any, headers map[string]string, httpClients *Clients) ApiInterface
}

type DefaultApiFactory struct{}

func (d *DefaultApiFactory) Create(url string, method string, queryParams map[string]string, bodyParams map[string]any, headers map[string]string, httpClients *Clients) ApiInterface {
	return NewOutboundCall(url, method, queryParams, bodyParams, headers, httpClients)
}
