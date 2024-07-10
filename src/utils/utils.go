package utils

import (
	"fmt"
	"io"
	"log"
	"net/url"
)

func Close(c io.Closer) {
	if err := c.Close(); err != nil {
		log.Fatal(err)
	}
}

func Ptr[T any](d T) *T {
	return &d
}

func HostURL(u *url.URL) *url.URL {
	base, err := url.Parse(fmt.Sprintf("%s://%s", u.Scheme, u.Host))
	if err != nil {
		panic(err)
	}

	return base
}

func SafeDeref[T any](p *T) T {
	if p == nil {
		var v T
		return v
	}
	return *p
}
