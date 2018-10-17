package requests

import (
	"io"
	"net/http"
	"net/url"
	"strings"
)

type ReqOption interface {
	Bind(req *http.Request) error
}

func NewReqOption(f func(req *http.Request) error) ReqOption {
	return &normalReqOption{f}
}

type normalReqOption struct {
	f func(req *http.Request) error
}

func (r *normalReqOption) Bind(req *http.Request) error {
	return r.f(req)
}

type bodyReqOption struct {
	body io.Reader
	fs   []func(req *http.Request) error
}

func (r *bodyReqOption) Bind(req *http.Request) error {
	for _, v := range r.fs {
		if err := v(req); err != nil {
			return err
		}
	}
	return nil
}

func ReqOptionBody(body io.Reader, f ...func(req *http.Request) error) ReqOption {
	return &bodyReqOption{
		body: body,
	}
}

func ReqOptionHeaderKV(k, v string) ReqOption {
	return NewReqOption(func(req *http.Request) error {
		req.Header.Set(k, v)
		return nil
	})
}

func ReqOptionHeaderMap(header map[string]string) ReqOption {
	return NewReqOption(func(req *http.Request) error {
		for k, v := range header {
			req.Header.Set(k, v)
		}
		return nil
	})
}

func ReqOptionQueryKV(k, v string) ReqOption {
	return NewReqOption(func(req *http.Request) error {
		q := req.URL.Query()
		q.Add(k, v)
		req.URL.RawQuery = q.Encode()
		return nil
	})
}

func ReqOptionQueryMap(query map[string]string) ReqOption {
	return NewReqOption(func(req *http.Request) error {
		if query != nil {
			q := req.URL.Query()
			for k, v := range query {
				q.Add(k, v)
			}
			req.URL.RawQuery = q.Encode()
		}
		return nil
	})
}

func ReqOptionFrom(form map[string]string) ReqOption {
	values := make(url.Values)
	for k, v := range form {
		values.Set(k, v)
	}

	return ReqOptionBody(strings.NewReader(values.Encode()), func(req *http.Request) error {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		return nil
	})
}
