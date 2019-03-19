package main

import (
	"bytes"
	"encoding/base64"
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

	"strings"
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

	return &OwlHttp{Client: client, Query: url.Values{}, Param: url.Values{}}

}

//设置Request的Body
func (h *OwlHttp) Body(body io.Reader) {
	rc, ok := body.(io.ReadCloser)
	if !ok && body != nil {
		rc = ioutil.NopCloser(body)
	}
	if body != nil {
		switch v := body.(type) {
		case *bytes.Buffer:
			h.Request.ContentLength = int64(v.Len())
		case *bytes.Reader:
			h.Request.ContentLength = int64(v.Len())
		case *strings.Reader:
			h.Request.ContentLength = int64(v.Len())
		}
	}
	h.Request.Body = rc
}

//设置Request的Cookie
func (h *OwlHttp) AddCookie(key, value string) {
	h.Request.AddCookie(&http.Cookie{Name: key, Value: value})
}

//设置Request的User-Agent
func (h *OwlHttp) UserAgent(UA string) {
	h.Request.Header.Set("User-Agent", UA)
}

//设置Request的header的Host
func (h *OwlHttp) Host(hostname string) {
	h.Request.Host = hostname
}

//返回Header
func (h *OwlHttp) Header() http.Header {
	return h.Request.Header
}

//设置BasicAuth
func (h *OwlHttp) BasicAuth(username, password string) {
	auth := username + ":" + password
	h.Request.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(auth)))
}

//设置Request的timeout
func (h *OwlHttp) SetTimeout(t time.Duration) {
	h.Client.Timeout = t
}

//创建Request实体
func newRequest(method, Url string) *http.Request {
	if !strings.HasPrefix(Url, "//") {
		if !strings.HasPrefix(Url, "http://") && !strings.HasPrefix(Url, "https://") {
			Url = "http://" + Url
		}
	}
	u, err := url.Parse(Url)

	if err != nil {
		//panic(err.Error())
		log.Println(err.Error())
		fmt.Println(err.Error())
	}
	if u.Scheme == "" {
		u.Scheme = "http"
	}
	req := &http.Request{
		Method:     method,
		URL:        u,
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Header:     make(http.Header),
		Host:       u.Host,
	}
	return req
}

//刷新旧查询和旧参数
func (h *OwlHttp) New() *OwlHttp {
	h.Query = url.Values{}
	h.Param = url.Values{}
	return h
}

//设置GET请求
func (h *OwlHttp) Get(Url string) {
	h.Request = newRequest(http.MethodGet, Url)
}

//设置POST请求
func (h *OwlHttp) Post(Url, bodyType string, body io.Reader) {
	r := newRequest(http.MethodPost, Url)
	r.Header.Set("Content-Type", bodyType)
	h.Request = r
	h.Body(body)
}

//设置POST请求(表单形式)
func (h *OwlHttp) PostForm(Url string) {
	h.Request = newRequest(http.MethodPost, Url)
	h.Request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if len(h.Param) != 0 {
		h.Body(strings.NewReader(h.Param.Encode()))
	}
}

//Do  return Response and err
func (h *OwlHttp) Do() (*Response, error) {
	rawquery := h.Query.Encode()
	if rawquery != "" && h.Request.URL.RawQuery != "" {
		rawquery = "&" + rawquery
	}
	h.Request.URL.RawQuery += rawquery
	if len(h.Param) > 0 {
		h.Body(strings.NewReader(h.Param.Encode()))
	}
	res, err := h.Client.Do(h.Request)
	if err != nil {
		return nil, err
	}
	return &Response{res}, nil
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
