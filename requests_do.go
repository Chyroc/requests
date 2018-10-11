package requests

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

type RequestOption interface {
	Bind(req *http.Request) error
}

func NewRequestOption(f func(req *http.Request) error) RequestOption {
	return &normalRequestOption{f}
}

type normalRequestOption struct {
	f func(req *http.Request) error
}

func (r *normalRequestOption) Bind(req *http.Request) error {
	return r.f(req)
}

type bodyRequestOption struct {
	body io.Reader
}

func (r *bodyRequestOption) Bind(req *http.Request) error {
	return fmt.Errorf("body request option no nedd cal bind")
}

func ReqOptionBody(body io.Reader) RequestOption {
	return &bodyRequestOption{
		body: body,
	}
}

func ReqOptionAddHeaderKV(k, v string) RequestOption {
	return NewRequestOption(func(req *http.Request) error {
		req.Header.Set(k, v)
		return nil
	})
}

func (r *requests) Get(url string, options ...RequestOption) ([]byte, error) {
	return r.request(http.MethodGet, url, options...)
}

func (r *requests) Post(url string, options ...RequestOption) ([]byte, error) {
	return r.request(http.MethodPost, url, options...)
}

func (r *requests) request(method, url string, options ...RequestOption) ([]byte, error) {
	if r.err != nil {
		return nil, r.err
	}

	var body io.Reader
	for _, v := range options {
		if bodyRequestOption, ok := v.(*bodyRequestOption); ok {
			body = bodyRequestOption.body
			break
		}
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	for _, v := range options {
		if _, ok := v.(*bodyRequestOption); !ok {
			if err := v.Bind(req); err != nil {
				return nil, err
			}
		}
	}

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}
