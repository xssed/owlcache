package httpclient

var DefaultTransport *http.Transport = &http.Transport{
	Dial: (&net.Dialer{
		Timeout: 5 * time.Second,
		//KeepAlive: 30 * time.Second,
	}).Dial,
	TLSHandshakeTimeout:   10 * time.Second,
	ExpectContinueTimeout: 1 * time.Second,
}

////Get方法
//func Get(Url string) (*Response, error) {
//	c := New()
//	c.Request = newRequest(http.MethodGet, Url)
//	return c.Do()
//}

////Post方法
//func Post(Url, bodyType string, body io.Reader) (*Response, error) {
//	c := New()
//	c.Request = newRequest(http.MethodPost, Url)
//	c.Request.Header.Set("Content-Type", bodyType)
//	c.Body(body)
//	return c.Do()
//}

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
func (h *Http) New() *Http {
	h.Query = url.Values{}
	h.Param = url.Values{}
	return h
}

//设置GET请求
func (h *Http) Get(Url string) {
	h.Request = newRequest(http.MethodGet, Url)
}

//设置POST请求
func (h *Http) Post(Url, bodyType string, body io.Reader) {
	r := newRequest(http.MethodPost, Url)
	r.Header.Set("Content-Type", bodyType)
	h.Request = r
	h.Body(body)
}

//设置POST请求(表单形式)
func (h *Http) PostForm(Url string) {
	h.Request = newRequest(http.MethodPost, Url)
	h.Request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if len(h.Param) != 0 {
		h.Body(strings.NewReader(h.Param.Encode()))
	}
}

//设置Request的Body
func (h *Http) Body(body io.Reader) {
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
func (h *Http) AddCookie(key, value string) {
	h.Request.AddCookie(&http.Cookie{Name: key, Value: value})
}

//设置Request的User-Agent
func (h *Http) UserAgent(UA string) {
	h.Request.Header.Set("User-Agent", UA)
}

//设置Request的header的Host
func (h *Http) Host(hostname string) {
	h.Request.Host = hostname
}

//返回Header
func (h *Http) Header() http.Header {
	return h.Request.Header
}

//设置BasicAuth
func (h *Http) BasicAuth(username, password string) {
	auth := username + ":" + password
	h.Request.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(auth)))
}

//设置Request的timeout
func (h *Http) SetTimeout(t time.Duration) {
	h.Client.Timeout = t
}

//Do  return Response and err
func (h *Http) Do() (*Response, error) {
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

func main() {
	//get
	//c := httpclient.New()
	c := New()
	c.Get("httpbin.org/get") //if url Scheme=="" default:http://www.google.com
	c.SetTimeout(3 * time.Second)
	c.Query.Add("key", "value")
	res, err := c.Do()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res.StatusCode)
	fmt.Println(res.String())
	//post
	c2 := New()
	c2.PostForm("www.google.com/example/api")
	c2.Param.Add("key", "value")
	res1, err1 := c2.Do()
	if err1 != nil {
		fmt.Println(err1)
		return
	}
	fmt.Println(res1.String())
}
