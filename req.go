package requests

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

func (r *Request) Get(url string, options ...ReqOption) (Response, error) {
	return r.request(http.MethodGet, url, options...)
}

func (r *Request) Post(url string, options ...ReqOption) (Response, error) {
	return r.request(http.MethodPost, url, options...)
}

func (r *Request) Delete(url string, options ...ReqOption) (Response, error) {
	return r.request(http.MethodDelete, url, options...)
}

func (r *Request) Put(url string, options ...ReqOption) (Response, error) {
	return r.request(http.MethodPut, url, options...)
}

func (r *Request) request(method, uri string, options ...ReqOption) (Response, error) {
	if r.err != nil {
		return Response{}, r.err
	}

	URL, err := url.Parse(uri)
	if err != nil {
		return Response{}, err
	}
	r.lock.Lock()
	r.urls[URL.Host] = URL
	r.lock.Unlock()

	var body io.Reader
	for _, v := range options {
		if bodyRequestOption, ok := v.(*bodyReqOption); ok {
			body = bodyRequestOption.body
			break
		}
	}

	req, err := http.NewRequest(method, uri, body)
	if err != nil {
		return Response{}, err
	}

	// hook option
	for _, v := range r.reqOptions {
		if err := v.Bind(req); err != nil {
			return Response{}, err
		}
	}

	// option
	for _, v := range options {
		if err := v.Bind(req); err != nil {
			return Response{}, err
		}
	}

	resp, err := r.client.Do(req)
	if err != nil {
		return Response{}, err
	}
	defer resp.Body.Close()

	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Response{}, err
	}

	res := &Response{Bytes: bs, Response: resp}

	// hook option
	for _, v := range r.respOptions {
		if err := v(res); err != nil {
			return Response{}, err
		}
	}

	return *res, nil
}
