package httpclient

import (
	"net/http"
)

//定义HTTP客户端结构
type OwlClient struct {
	OwlTransport *http.Transport
}

//创建一个HTTP客户端
func NewOwlClient() *OwlClient {
	owltransport := NewOwlTransport()
	owlhttpclient := &OwlClient{OwlTransport: owltransport}
	return owlhttpclient
}
