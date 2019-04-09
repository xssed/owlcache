//关于这个包我不是很满意。。。。因为是Owl集群的工作模式是server to server
//想提升效率又怕用net包自己写client后面维护成本高 以后官方包http2 client完善再说
package httpclient

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	MaxIdleConnections int = 20
	RequestTimeout     int = 2 //初始默认值，在配置文件中已经定义。如果调用OwlHttp.SetTimeout()可以在使用时再次更改
)

//var OwlTransport *http.Transport

func NewOwlTransport() *http.Transport {
	OwlTransport := &http.Transport{
		MaxIdleConnsPerHost: MaxIdleConnections,
	}
	return OwlTransport
}

type OwlHttp struct {
	Request *http.Request
	Client  *http.Client
	Query   url.Values //QueryString。url.Values结构是map[string][]string非并发安全
	Param   url.Values //PostFromParams。url.Values结构是map[string][]string非并发安全
}

//创建HttpClient实体
func NewOwlHttpClient(OwlTransport *http.Transport) *OwlHttp {

	client := &http.Client{
		Transport: OwlTransport,
		Timeout:   time.Duration(RequestTimeout) * time.Second,
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
func (h *OwlHttp) UserAgent(useragent string) {
	h.Request.Header.Set("User-Agent", useragent)
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

//清空查询数据
func (h *OwlHttp) Claer() *OwlHttp {
	h.Query = url.Values{}
	h.Param = url.Values{}
	return h
}

//打印数据，测试用
func (h *OwlHttp) EchoInfo() {
	fmt.Println(*h.Request)
	fmt.Println(*h.Client)
	fmt.Println(h.Query)
	fmt.Println(h.Param)
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
