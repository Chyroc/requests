package requests_test

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Chyroc/requests"
)

var (
	url = "http://examples.com"
)

func Example_defaultClient() {
	// get
	resp, err := requests.Default.Get(url)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(resp.Bytes))

	// post
	resp, err = requests.Default.Post(url)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(resp.Bytes))
}

func Example_withCookie() {
	r := requests.New(requests.OptionEnableCookie)

	resp, err := r.Get(url)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(resp.Bytes))
	fmt.Println(r.Cookies())
}

func Example_withHeader() {
	r := requests.New(requests.OptionEnableCookie)

	resp, err := r.Get(url, requests.ReqOptionAddHeaderKV("test", "test-header"))
	if err != nil {
		panic(err)
	}

	fmt.Println(string(resp.Bytes))
}

func Example_withBody() {
	r := requests.New(requests.OptionEnableCookie)

	resp, err := r.Get(url, requests.ReqOptionBody(strings.NewReader("a=b&c=d")))
	if err != nil {
		panic(err)
	}

	fmt.Println(string(resp.Bytes))
}

func Example_customized() {
	r := requests.New(requests.OptionEnableCookie)

	resp, err := r.Get(url,
		requests.ReqOptionBody(strings.NewReader("a=b&c=d")),
		requests.NewRequestOption(func(req *http.Request) error {
			req.SetBasicAuth("name", "password")
			return nil
		}),
	)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(resp.Bytes))
}
