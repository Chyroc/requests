package requests_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Chyroc/requests"
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

func TestNew(t *testing.T) {
	as := assert.New(t)
	url := NewTestURL("/with-cookie")
	t.Log(url)

	noCookie := requests.New()
	withCookie := requests.New(requests.OptionEnableCookie)

	{
		resp, err := noCookie.Get(url)
		as.Nil(err)
		as.Equal("with cookie!", string(resp.Bytes))
		as.Nil(noCookie.Cookies())
	}

	{
		resp, err := withCookie.Get(url)
		as.Nil(err)
		as.Equal("with cookie!", string(resp.Bytes))
		c := withCookie.Cookies()
		as.Len(c, 1)
		as.Equal("a", c[0].Name)
		as.Equal("b", c[0].Value)
	}

	{
		resp, err := withCookie.Get(url)
		as.Nil(err)
		as.Equal("cookie: b", string(resp.Bytes))
		c := withCookie.Cookies()
		as.Len(c, 1)
		as.Equal("a", c[0].Name)
		as.Equal("b", c[0].Value)
	}
}
