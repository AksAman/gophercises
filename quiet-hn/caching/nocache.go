package caching

type NoCache[T CacheItem] struct {
}

func (c *NoCache[T]) ToString() string {
	return "NoCache"
}

func (c *NoCache[T]) Init() error {
	return nil
}

func (c *NoCache[T]) SetupTicker(event func()) {
}

func (c *NoCache[T]) IsExpired() bool {
	return true
}

func (c *NoCache[T]) Set(items []*T) error {
	return nil
}

func (c *NoCache[T]) Get() []*T {
	return nil
}
