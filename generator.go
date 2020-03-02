package idGenerator

type IdGenerator interface {
	GetUID() (id uint64, err error)
}
