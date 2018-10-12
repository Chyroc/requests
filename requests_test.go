package requests_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Chyroc/requests"
)

var (
	baseURL string
	once    = new(sync.Once)
)

var serverN = 50
var appN = 50

func NewServer() http.Handler {
	mux := http.NewServeMux()
	for i := 1; i <= appN; i++ {
		mux.HandleFunc("/"+strconv.Itoa(i), func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, strconv.Itoa(i))
		})
	}
	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "hello")
	})
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

func NewNServerURL() func(i int) string {
	var s []string

	for i := 1; i <= serverN; i++ {
		b := httptest.NewServer(NewServer()).URL
		for i := 1; i <= appN; i++ {
			s = append(s, b+"/"+strconv.Itoa(i))
		}
	}

	return func(i int) string {
		index := i % (appN * serverN)
		return s[index]
	}
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

func BenchmarkRequest(b *testing.B) {
	as := assert.New(b)
	url := NewNServerURL()

	b.Run("requests n-client n-times", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			r := requests.New()
			_, err := r.Get(url(i))
			as.Nil(err)
		}
	})

	b.Run("requests one-client n-times", func(b *testing.B) {
		r := requests.New()
		for i := 0; i < b.N; i++ {
			_, err := r.Get(url(i))
			as.Nil(err)
		}
	})

	b.Run("http n-client n-times", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			r := new(http.Client)
			resp, err := r.Get(url(i))
			as.Nil(err)
			_, err = ioutil.ReadAll(resp.Body)
			resp.Body.Close()
		}
	})

	b.Run("http one-client n-times", func(b *testing.B) {
		r := new(http.Client)
		for i := 0; i < b.N; i++ {
			resp, err := r.Get(url(i))
			as.Nil(err)
			_, err = ioutil.ReadAll(resp.Body)
			resp.Body.Close()
		}
	})

}
