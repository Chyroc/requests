package requests

import (
	"golang.org/x/net/publicsuffix"
	"net/http"
	"net/http/cookiejar"
)

type requests struct {
	err error

	client *http.Client
}

var Default = &requests{}

type Option func(r *requests) error

// 开启 cookie option
func OptionEnableCookie(r *requests) error {
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		return err
	}
	r.client = new(http.Client)
	r.client.Jar = jar
	return nil
}

func New(options ...Option) *requests {
	r := &requests{}
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
