package main

import (
	"fmt"
	"time"

	"github.com/xssed/owlcache/network/httpclient"
)

var owlclient *httpclient.OwlHttp

func main() {

	owlclient := httpclient.NewOwlHttpClient()
	fmt.Println(owlclient)

	//		urladdress := "https://httpbin.org/get"

	//		v := url.Values{}
	//		v.Set("cmd", "get")
	//		v.Set("key", "hello")

	//		fmt.Println(v.Encode())

	//		req, err := http.NewRequest("GET", urladdress, nil) //bytes.NewBuffer([]byte("?"+v.Encode()))
	//		if err != nil {
	//			log.Fatalf("Error Occured. %+v", err)
	//		}
	//		//req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	//		response, err := httpclient.Client.Do(req)
	//		if err != nil && response == nil {
	//			log.Fatalf("Error sending request to API endpoint. %+v", err)
	//		}

	//		// Close the connection to reuse it
	//		defer response.Body.Close()

	//		//	for {
	//		body, err := ioutil.ReadAll(response.Body)
	//		if err != nil {
	//			log.Fatalf("Couldn't parse response body. %+v", err)
	//		}

	//		log.Println(string(body))
	//		}

	//time.Sleep(time.Second * 2000)

	//c := New()

	go func() {
		owlclient.Get("httpbin.org/get")
		owlclient.SetTimeout(3 * time.Second)
		owlclient.Query.Add("key", "value")
		res, err := owlclient.Do()
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(res.StatusCode)
		fmt.Println(res.String())

		owlclient.Claer()
	}()

	go func() {
		owlclient.Get("httpbin.org/get")
		owlclient.SetTimeout(3 * time.Second)
		owlclient.Query.Add("key22", "value22")
		res2, err2 := owlclient.Do()
		if err2 != nil {
			fmt.Println(err2)
			return
		}
		fmt.Println(res2.StatusCode)
		fmt.Println(res2.String())

		owlclient.Claer()
	}()

	go func() {
		owlclient.PostForm("httpbin.org/post")
		owlclient.SetTimeout(3 * time.Second)
		owlclient.Query.Add("key22", "value22")
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
