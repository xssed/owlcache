package main

import (
	"bytes"
	//"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"

	//"net"
	"net/http"
	"net/url"
	"os"
	"path"

	//"strings"
	"time"
	//"errors"
)

const (
	MaxIdleConnections int = 20
	RequestTimeout     int = 5
)

type OwlHttp struct {
	Request *http.Request
	Client  *http.Client
	Query   url.Values //QueryString
	Param   url.Values //PostFromParams
}

var httpclient *OwlHttp

//创建HttpClient实体
func newOwlHttpClient() *OwlHttp {

	client := &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: MaxIdleConnections,
		},
		Timeout: time.Duration(RequestTimeout) * time.Second,
	}

	return &OwlHttp{
		Client: client,
		Query:  url.Values{},
		Param:  url.Values{}}

}

type Response struct {
	*http.Response
}

//获取的数据结果转[]byte
func (r *Response) Byte() []byte {
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		return nil
	}
	return b
}

//获取的数据结果转string
func (r *Response) String() string {
	return string(r.Byte())
}

//Response body save as a file
func (r *Response) DownLoadFile(filepath string) error {
	dir, _ := path.Split(filepath)
	if dir != "" {
		if err := os.MkdirAll(dir, 0666); err != nil {
			return err
		}
	}
	f, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer f.Close()
	io.Copy(f, bytes.NewReader(r.Byte()))
	return nil
}

//Json.Unmarshal ResponseBody
func (r *Response) JsonUnmarshal(v interface{}) error {
	return json.Unmarshal(r.Byte(), v)
}

func main() {

	httpclient := newOwlHttpClient()

	urladdress := "https://httpbin.org/get"

	v := url.Values{}
	v.Set("cmd", "get")
	v.Set("key", "hello")

	fmt.Println(v.Encode())

	req, err := http.NewRequest("GET", urladdress, nil) //bytes.NewBuffer([]byte("?"+v.Encode()))
	if err != nil {
		log.Fatalf("Error Occured. %+v", err)
	}
	//req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	response, err := httpclient.Client.Do(req)
	if err != nil && response == nil {
		log.Fatalf("Error sending request to API endpoint. %+v", err)
	}

	// Close the connection to reuse it
	defer response.Body.Close()

	//	for {
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatalf("Couldn't parse response body. %+v", err)
	}

	log.Println(string(body))
	//	}

	//time.Sleep(time.Second * 2000)

}
