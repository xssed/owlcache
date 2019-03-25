package httpclient

import (
	"fmt"
	"time"
)

//发送GET请求
func (c *OwlClient) Get(address string, value string) {
	owlclient := NewOwlHttpClient(c.OwlTransport)
	fmt.Println(owlclient)
	owlclient.Get("httpbin.org/get")
	owlclient.SetTimeout(*c.HCRequestTimeout * time.Second)
	owlclient.Query.Add("key", "value")
	res, err := owlclient.Do()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res.StatusCode)
	fmt.Println(res.String())

	owlclient.Claer()
}
