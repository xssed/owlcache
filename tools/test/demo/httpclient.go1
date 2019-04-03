package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/xssed/owlcache/network/httpclient"
)

var OwlTransport *http.Transport

func main() {

	OwlTransport = httpclient.NewOwlTransport()

	//	go func() { //*http.Transport

	//		owlclient := httpclient.NewOwlHttpClient(OwlTransport)
	//		//fmt.Println(owlclient)
	//		owlclient.Get("httpbin.org/get")
	//		owlclient.SetTimeout(3 * time.Second)
	//		owlclient.Query.Add("key", "value")
	//		res, err := owlclient.Do()
	//		if err != nil {
	//			fmt.Println(err)
	//			return
	//		}
	//		fmt.Println(res.StatusCode)
	//		fmt.Println(res.String())

	//		owlclient.Claer()
	//	}() //OwlTransport

	//	go func() {
	//		owlclient := httpclient.NewOwlHttpClient(OwlTransport)
	//		//fmt.Println(owlclient)
	//		owlclient.Get("httpbin1111.org/get")
	//		owlclient.SetTimeout(3 * time.Second)
	//		owlclient.Query.Add("key22", "value22")
	//		res2, err2 := owlclient.Do()
	//		if err2 != nil {
	//			fmt.Println(err2)
	//			return
	//		}
	//		fmt.Println(res2.StatusCode)
	//		fmt.Println(res2.String())

	//		owlclient.Claer()
	//	}()

	//	go func() {
	//		owlclient := httpclient.NewOwlHttpClient(OwlTransport)
	//		//fmt.Println(owlclient)
	//		owlclient.PostForm("httpbin.org/post")
	//		owlclient.SetTimeout(3 * time.Second)
	//		owlclient.Query.Add("key22", "value22")
	//		res2, err2 := owlclient.Do()
	//		if err2 != nil {
	//			fmt.Println(err2)
	//			return
	//		}
	//		fmt.Println(res2.StatusCode)
	//		fmt.Println(res2.String())

	//		owlclient.Claer()
	//	}()

	go func() {
		owlclient := httpclient.NewOwlHttpClient(OwlTransport)
		//fmt.Println(owlclient)
		owlclient.PostForm("127.0.0.1:7721" + "/data")
		owlclient.SetTimeout(3 * time.Second)
		owlclient.Query.Add("cmd", "get")
		owlclient.Query.Add("key", "hello")

		owlclient.EchoInfo()

		res2, err2 := owlclient.Do()
		if err2 != nil {
			fmt.Println(err2)
			return
		}
		fmt.Println(res2.StatusCode)
		fmt.Println(res2.String())

		owlclient.Claer()
	}()

	time.Sleep(time.Second * 100)
}
