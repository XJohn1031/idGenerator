package idGenerator

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func prepare() *IDStruct {
	return NewIDStruct(28, 22, 13)
}

func TestID_GenerateID(t *testing.T) {
	id_struct := prepare()

	timeStamp := time.Now().Unix()
	result := []int{0,
		8193,
		16386,
		24579,
		32772,
		40965,
		49158,
		57351,
		65544,
		73737}
	for i := 0; i < 10; i++ {
		id, err := id_struct.GenerateID(uint64(timeStamp), uint64(i), uint64(i))
		fmt.Println(id)
		assert.NoError(t, err)
		fmt.Println(id & (1<<35 - 1))
		assert.Equal(t, int(id&(1<<35-1)), result[i])
	}

}
