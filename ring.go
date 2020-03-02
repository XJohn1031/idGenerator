package idGenerator

import (
	"log"
	"sync"
	"sync/atomic"
	"time"

	"github.com/pkg/errors"
)

type Ring struct {
	consumePosition int64 // 表示当前消费的位置
	productPosition int64 // 表示当前生产的位置

	slot  []uint64    // 暂时存储
	flags []*FlagSlot // true表示当前位置可以取出, false表示当前位置应该放置

	paddingSize int // 填充大小
	bufferSize  int // slot的大小

	paddingExecutor PaddingExecutor // 执行填充任务

	productFunc ProductFunc

	idStruct *IDStruct
}

type ProductFunc func(currentTimestamp uint64)

type FlagSlot struct {
	flag bool
	*sync.Mutex
}

func NewRing(bufferSize int, workId uint64) *Ring {
	now := uint64(time.Now().Unix())
	ring := &Ring{
		consumePosition: -1,
		productPosition: -1,
		slot:            make([]uint64, bufferSize),
		flags:           InitFlag(bufferSize),
		paddingSize:     bufferSize / 2,
		bufferSize:      bufferSize,
		paddingExecutor: PaddingExecutor{
			whetherPadding: NotPadding,
			currentSecond:  now,
			workId:         workId,
			provider:       nil,
		},
		idStruct: NewIDStruct(28, 18, 17),
	}
	ring.productFunc = ring.Product
	ring.paddingExecutor.provider = ring.Provider

	ring.productFunc(uint64(time.Now().Unix()))

	go ring.SchedulePut()
	return ring
}

func (r *Ring) SchedulePut() {
	schedule := time.NewTicker(time.Second)
	for {
		select {
		case <-schedule.C:
			r.productFunc(uint64(time.Now().Unix()))
		}
	}
}

func (r *Ring) Provider(timestamp uint64, workId uint64) (ids []uint64, err error) {
	firstId, err := r.idStruct.GenerateID(timestamp, workId, 0)
	if err != nil {
		log.Printf("err is %v", err)
		return []uint64{}, errors.Wrap(err, "raw generate err")
	}

	result := make([]uint64, r.bufferSize)
	for i := 0; i < r.bufferSize; i++ {
		result[i] = firstId + uint64(i)
	}
	return result, nil
}

func (r *Ring) Put(uid uint64) bool {
	if (r.productPosition - r.consumePosition) == int64(r.bufferSize-1) {
		log.Println("the ring is full")
		return false
	}

	nextIndex := r.calculateNextPosition(r.productPosition + 1)

	r.flags[nextIndex].Lock()
	defer r.flags[nextIndex].Unlock()
	if r.flags[nextIndex].flag {
		log.Println("the position is not able to put")
		return false
	}

	r.slot[nextIndex] = uid
	r.flags[nextIndex].flag = true

	atomic.AddInt64(&r.productPosition, 1)
	return true
}

func (r *Ring) calculateNextPosition(input int64) int {
	return int(int64(r.bufferSize-1) & input)
}

func (r *Ring) Take() (id uint64, err error) {
	next := atomic.AddInt64(&r.consumePosition, 1)
	//if r.productPosition-next < uint64(r.paddingSize) {
	//	log.Printf("wait to padding, %v, %v", r.productPosition, next)
	//	timeStamp := time.Now().Unix()
	//	go r.productFunc(uint64(timeStamp))
	//}

	if next == r.productPosition {
		return 0, errors.New("nothing to take")
	}

	nextPosition := r.calculateNextPosition(r.consumePosition)
	r.flags[nextPosition].Lock()
	defer r.flags[nextPosition].Unlock()
	if !r.flags[nextPosition].flag {
		return 0, errors.New("no available uid")
	}

	id = r.slot[nextPosition]
	r.flags[nextPosition].flag = false
	return id, nil
}

// 当productPosition - consumePosition < paddingSize时, 需要补充id
func (r *Ring) Product(timestamp uint64) {
	ids := r.paddingExecutor.execute()

	for i := range ids {
		put := r.Put(ids[i])
		if !put {
			return
		}
	}
}

// 初始化flag数组
func InitFlag(bufferSize int) []*FlagSlot {
	log.Println("begin to init flag")
	result := make([]*FlagSlot, bufferSize)
	for i := 0; i < bufferSize; i++ {
		result[i] = &FlagSlot{
			flag:  false,
			Mutex: &sync.Mutex{},
		}
	}

	log.Println("end to init flag")

	return result
}
