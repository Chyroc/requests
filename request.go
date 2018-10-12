package requests

import (
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

type Request struct {
	err error

	urls   map[string]*url.URL
	client *http.Client

	reqOptions  []ReqOption
	respOptions []RespOption
}

var Default = &Request{}

type Option func(r *Request) error

// 开启 cookie option
func OptionEnableCookie(r *Request) error {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return err
	}
	r.client = new(http.Client)
	r.client.Jar = jar
	return nil
}

func New(options ...Option) *Request {
	r := &Request{urls: make(map[string]*url.URL)}
	for _, option := range options {
		if r.err = option(r); r.err != nil {
			return r
		}
	}
	if r.client == nil {
		r.client = new(http.Client)
	}

	return r
}
