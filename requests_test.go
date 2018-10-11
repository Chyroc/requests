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
		bs, err := noCookie.Get(url)
		as.Nil(err)
		t.Log(string(bs))
		as.Equal("with cookie!", string(bs))
	}

	{
		bs, err := withCookie.Get(url)
		as.Nil(err)
		t.Log(string(bs))
		as.Equal("with cookie!", string(bs))
	}

}
