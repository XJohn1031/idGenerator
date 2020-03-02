package idGenerator

import (
	"log"
	"sync/atomic"
	"time"
)

const (
	Padding    = 1
	NotPadding = 0
)

type PaddingExecutor struct {
	whetherPadding int32  // 是否正在载入
	currentSecond  uint64 // 当前秒数

	workId uint64

	provider UidProvider
}

type UidProvider func(currentSecond uint64, workId uint64) (ids []uint64, err error)

func (p *PaddingExecutor) execute() []uint64 {

	if !atomic.CompareAndSwapInt32(&p.whetherPadding, NotPadding, Padding) {
		log.Println("executor is running")
		return []uint64{}
	}

	nextSecond := uint64(time.Now().Unix())
	for p.currentSecond >= nextSecond {
		time.Sleep(1e7)
		nextSecond = uint64(time.Now().Unix())
	}
	ids, e := p.provider(nextSecond, p.workId)
	if e != nil {
		return []uint64{}
	}

	p.currentSecond = nextSecond
	atomic.CompareAndSwapInt32(&p.whetherPadding, Padding, NotPadding)
	return ids
}
