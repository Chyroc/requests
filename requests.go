package requests

import (
	"encoding/json"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

type Request struct {
	err error

	urls   map[string]*url.URL
	client *http.Client
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

func (r *Request) Cookies() []*http.Cookie {
	if r.client.Jar == nil {
		return nil
	}

	var cookies []*http.Cookie
	for _, u := range r.urls {
		c := r.client.Jar.Cookies(u)
		if c != nil {
			cookies = append(cookies, c...)
		}
	}
	return cookies
}

func (r *Request) CookiesSring() (string, error) {
	cookies := r.Cookies()
	if cookies == nil {
		return "", nil
	}

	bs, err := json.Marshal(cookies)
	if err != nil {
		return "", err
	}

	return string(bs), nil
}
