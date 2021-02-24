package idGenerator

import (
	"log"
	"testing"
	"time"
)

func TestRing_Take(t *testing.T) {
	ring := NewRing(1<<17, 1)
	result := make(map[uint64]bool, 0)
	ticker := time.NewTicker(1e9)

	time.Sleep(1e10)
	for {
		select {
		case <-ticker.C:
			for i := 0; i < 1<<16; i++ {
				id, err := ring.Take()
				if err != nil {
					log.Println("take id err", err)
				}
				result[id] = true
			}
			log.Println("id result size", len(result))
		}
	}

}
