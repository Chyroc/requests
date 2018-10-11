package requests_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
)

var (
	baseURL string
	once    = new(sync.Once)
)

func NewServer() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/with-cookie", func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("a")
		if err != nil && err != http.ErrNoCookie {
			fmt.Fprintf(w, err.Error())
			return
		}

		if c == nil {
			w.Header().Set("Set-Cookie", "a=b")
			fmt.Fprintf(w, "with cookie!")
			return
		}

		fmt.Fprintf(w, "cookie: "+c.Value)
	})
	return mux
}

func NewTestURL(p string) string {
	once.Do(func() {
		baseURL = httptest.NewServer(NewServer()).URL
	})
	return baseURL + p
}
