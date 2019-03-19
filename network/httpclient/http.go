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
