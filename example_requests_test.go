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

	resp, err := r.Get(url, requests.ReqOptionHeaderKV("test", "test-header"))
	if err != nil {
		panic(err)
	}
	fmt.Println(string(resp.Bytes))

	resp2, err := r.Get(url, requests.ReqOptionHeaderMap(map[string]string{"test-2": "test-header"}))
	if err != nil {
		panic(err)
	}
	fmt.Println(string(resp2.Bytes))
}

func Example_withQuery() {
	r := requests.New(requests.OptionEnableCookie)

	resp, err := r.Get(url, requests.ReqOptionQueryKV("test", "test-query"))
	if err != nil {
		panic(err)
	}
	fmt.Println(string(resp.Bytes))

	resp2, err := r.Get(url, requests.ReqOptionQueryMap(map[string]string{"test-2": "test-query"}))
	if err != nil {
		panic(err)
	}
	fmt.Println(string(resp2.Bytes))
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
		requests.NewReqOption(func(req *http.Request) error {
			req.SetBasicAuth("name", "password")
			return nil
		}),
	)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(resp.Bytes))
}

func Example_hook() {
	r := requests.New(requests.OptionEnableCookie)

	// all request have this token-header
	r.ReqHook(requests.ReqOptionHeaderKV("Authorization", "token"))

	// all request well check resp-code
	r.RespHook(func(resp *requests.Response) error {
		var response = make(map[string]interface{})
		if err := resp.BindJson(&response); err != nil {
			return err
		}
		if response["code"].(int) != 0 {
			return fmt.Errorf("request err")
		}

		return nil
	})

	resp, err := r.Get(url)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(resp.Bytes))
}
