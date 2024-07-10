package utils

import (
	"io"
	"log"
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
