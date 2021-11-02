package shared

import (
	"math"
	"sync/atomic"
)

type SafeTimestamp struct {
	value int32
}

func (s *SafeTimestamp) MaxInc(otherTime int32) {
	timestamp := math.Max(float64(s.value), float64(otherTime)) + 1
	s.value = atomic.SwapInt32(&s.value, int32(timestamp))
}

func (s *SafeTimestamp) Increment() {
	s.value = atomic.AddInt32(&s.value, 1)
}

func (s *SafeTimestamp) Value() int32 {
	return s.value
}
