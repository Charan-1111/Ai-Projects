package clients

type ApiFactory interface {
	Create(url string, method string, queryParams map[string]string, bodyParams []byte, headers map[string]string, clients *Clients) ApiInterface
}

type DefaultApiFactory struct{}

func (df *DefaultApiFactory) Create(url string, method string, queryParams map[string]string, bodyParams []byte, headers map[string]string, clients *Clients) ApiInterface {
	return NewOutboundCall(url, method, queryParams, bodyParams, headers, clients)
}
