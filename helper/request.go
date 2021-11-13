package helper

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"nfly_client/common"
	"time"
)

type Request struct {
	client  http.Client
}

func NewRequest() *Request{
	jar,_ := cookiejar.New(nil)
	return &Request{
		client:  http.Client{
			Timeout: 5 * time.Second,
			Jar: jar,
		},
	}
}

func (x *Request) readResp(outer io.ReadCloser) (*common.Reply,error) {
	var reply common.Reply
	err := json.NewDecoder(outer).Decode(&reply)
	if err == nil {
		return &reply,nil
	}
	log.Println("parse server reply failed")
	return nil, err
}

// Get with session
func (x *Request) Get(url string) (*common.Reply,error) {
	resp,err := x.client.Get(url)
	if err == nil {
		return x.readResp(resp.Body)
	}
	return nil, err
}

func (x *Request) PostForm(url string, data url.Values) (*common.Reply,error) {
	resp,err := x.client.PostForm(url,data)
	if err == nil {
		return x.readResp(resp.Body)
	}
	return nil, err
}

func (x *Request) Post(url string,ct string,data io.Reader) (*common.Reply,error) {
	switch ct {
	case "json":
		ct = "application/json"
	case "text":
		ct = "text/plain"
	default:
		panic("not implemented")
	}
	resp,err := x.client.Post(url,ct,data)
	if err == nil {
		return x.readResp(resp.Body)
	}
	return nil, err
}

func (x *Request) Put(url string, data []byte) (*common.Reply,error) {
	req,_ := http.NewRequest(http.MethodPut,url,bytes.NewReader(data))
	resp,err := x.client.Do(req)
	if err == nil {
		return x.readResp(resp.Body)
	}
	return nil, err
}

func (x *Request) Delete(url string) (*common.Reply,error) {
	req,_ := http.NewRequest(http.MethodDelete,url,nil)
	resp,err := x.client.Do(req)
	if err == nil {
		return x.readResp(resp.Body)
	}
	return nil, err
}

func (x *Request) GetCookies(u *url.URL) []*http.Cookie{
	return x.client.Jar.Cookies(u)
}