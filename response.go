package requests

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	*http.Response
	Bytes []byte
}

func (r *Response) BindJson(v interface{}) error {
	return json.Unmarshal(r.Bytes, &v)
}
