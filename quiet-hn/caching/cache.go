package caching

type CacheItem interface {
}

type Cache[T CacheItem] interface {
	Init() error
	IsExpired() bool
	Get() []*T
	Set([]*T) error
	SetupTicker(func())
	ToString() string
}
