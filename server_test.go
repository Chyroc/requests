package requests_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"

	"github.com/gavv/httpexpect"
)

var (
	baseURL string
	once    = new(sync.Once)
)

func NewServer() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/no-cookie", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "no cookie!")
	})
	mux.HandleFunc("/with-cookie", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Set-Cookie", "cookie")
		fmt.Fprintf(w, "with cookie!")
	})
	return mux
}

type test struct {
	server http.Handler
	expect *httpexpect.Expect
}

func NewTestURL(p string) string {
	once.Do(func() {
		baseURL = httptest.NewServer(NewServer()).URL
	})
	return baseURL + p
}
