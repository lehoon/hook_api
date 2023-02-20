package utils

import (
	"math"
	"sync"
)

var (
	messageSequeceNo     = uint16(0x0001)
	messageSequeceNoLock sync.RWMutex
)

func GetNextSequenceNo() uint16 {
	messageSequeceNoLock.Lock()
	if messageSequeceNo == math.MaxUint16 - 1 {
		messageSequeceNo = 0
	}
	messageSequeceNo += 1
	messageSequeceNoLock.Unlock()
	return messageSequeceNo
}