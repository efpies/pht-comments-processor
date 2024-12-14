package model

type Platform string

const (
	Pht Platform = "pht"
)

var PlatformEnum = struct {
	Pht Platform
}{
	Pht: Pht,
}
