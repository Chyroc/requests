package requests

import "net/http"

func (r *Request) Cookies() []*http.Cookie {
	if r.client.Jar == nil {
		return nil
	}

	var cookies []*http.Cookie
	for _, u := range r.urls {
		c := r.client.Jar.Cookies(u)
		if c != nil {
			cookies = append(cookies, c...)
		}
	}
	return cookies
}
