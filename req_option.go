package requests

import (
	"fmt"
	"io"
	"net/http"
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
}

func (r *bodyReqOption) Bind(req *http.Request) error {
	return fmt.Errorf("body request option no nedd cal bind")
}

func ReqOptionBody(body io.Reader) ReqOption {
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
