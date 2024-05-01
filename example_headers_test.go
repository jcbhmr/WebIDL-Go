package webidl_test

import "fmt"

type t[A any, B any] struct {
	A A
	B B
}

type Headers interface {
	Append(name string, value string)
	Delete(name string)
	Get(name string) *string
	GetSetCookie() []string
	Has(name string) bool
	Set(name string, value string)
	Iter() func(func(key string, value string) bool)
}

type headersImpl struct {
	headers []t[string, string]
}

func (h *headersImpl) Append(name string, value string) {
	h.headers = append(h.headers, t[string, string]{name, value})
}

func (h *headersImpl) Delete(name string) {
	for i, header := range h.headers {
		if header.A == name {
			h.headers = append(h.headers[:i], h.headers[i+1:]...)
			return
		}
	}
}

func (h *headersImpl) Get(name string) *string {
	for _, header := range h.headers {
		if header.A == name {
			return &header.B
		}
	}
	return nil
}

func (h *headersImpl) GetSetCookie() []string {
	var cookies []string
	for _, header := range h.headers {
		if header.A == "Set-Cookie" {
			cookies = append(cookies, header.B)
		}
	}
	return cookies
}

func (h *headersImpl) Has(name string) bool {
	for _, header := range h.headers {
		if header.A == name {
			return true
		}
	}
	return false
}

func (h *headersImpl) Set(name string, value string) {
	for i, header := range h.headers {
		if header.A == name {
			h.headers[i].B = value
			return
		}
	}
	h.Append(name, value)
}

func (h *headersImpl) Iter() func(func(key string, value string) bool) {
	return func(cb func(key string, value string) bool) {
		for _, header := range h.headers {
			if !cb(header.A, header.B) {
				return
			}
		}
	}
}

func NewHeaders() Headers {
	return &headersImpl{
		headers: []t[string, string]{},
	}
}

func ExampleResponse() {
	headers := NewHeaders()
	headers.Append("Accept", "application/json")
	headers.Append("Accept", "text/html")
	headers.Iter()(func(key string, value string) bool {
		fmt.Println(key, value)
		return true
	})
	// Output:
	// Accept application/json
	// Accept text/html
}
