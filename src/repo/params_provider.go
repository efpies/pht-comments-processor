package repo

type ParamsProvider interface {
	PrefetchParams() error
	GetParam(key string) string
	UpdateParam(key string, value string) error
}
