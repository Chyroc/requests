package requests

import "net/http"

type Response struct {
	*http.Response
	Bytes []byte
}
