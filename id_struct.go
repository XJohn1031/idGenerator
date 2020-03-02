// 定义了
package idGenerator

import (
	"time"
)

var from = uint64(time.Date(2020, time.February, 29, 0, 0, 0, 0, time.Local).Unix())

type IDStruct struct {
	signBit      int // 符号位
	timeStampBit int // 时间戳标志
	workerBit    int // 生成器标志, 每次重启重新分配, 防止时钟回拨
	circleBit    int // 循环序列标志
}

func NewIDStruct(timeStampBit int, workerBit int, circleBit int) *IDStruct {
	return &IDStruct{signBit: 1, timeStampBit: timeStampBit, workerBit: workerBit, circleBit: circleBit}
}

func (i *IDStruct) GenerateID(timeStamp uint64, workerId uint64, sequence uint64) (id uint64, err error) {

	return (timeStamp-from)<<(i.workerBit+i.circleBit) |
		(workerId)<<(i.circleBit) |
		sequence, nil
}
