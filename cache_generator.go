package idGenerator

type CacheGenerator struct {
	ring *Ring
}

func (c *CacheGenerator) GetUID() (id uint64, err error) {
	return c.ring.Take()
}
