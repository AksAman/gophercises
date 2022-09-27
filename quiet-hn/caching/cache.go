package caching

type CacheItem interface {
}

type Cache[T CacheItem] interface {
	Init()
	IsExpired() bool
	Get() []*T
	Set([]*T)
}
