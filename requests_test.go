package requests_test

import (
	"github.com/Chyroc/requests"
	"github.com/stretchr/testify/assert"
	"testing"
)

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
