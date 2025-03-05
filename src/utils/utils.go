package utils

import (
	"fmt"
	"io"
	"log"
	"net/url"
	"strconv"
)

func Close(c io.Closer) {
	if err := c.Close(); err != nil {
		log.Fatal(err)
	}
}

func Ptr[T any](d T) *T {
	return &d
}

func SafeDeref[T any](p *T) T {
	if p == nil {
		var v T
		return v
	}
	return *p
}

func HostURL(u *url.URL) (*url.URL, error) {
	urlStr := fmt.Sprintf("%s://%s", u.Scheme, u.Host)
	if u.Port() != "" {
		urlStr = fmt.Sprintf("%s:%s", urlStr, u.Port())
	}

	base, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	return base, nil
}

func AtoiSafe(s string) (int, error) {
	if s == "" {
		return 0, nil
	}

	i, err := strconv.Atoi(s)
	if err != nil {

		return 0, err
	}

	return i, nil
}
