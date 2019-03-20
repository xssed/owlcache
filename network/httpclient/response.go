package httpclient

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
)

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
