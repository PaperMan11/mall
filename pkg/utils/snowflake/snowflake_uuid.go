package snowflake

import (
	"fmt"
	"sync"
	"time"
)

const (
	workerBits     uint8 = 10
	maxWorker      int64 = -1 ^ (-1 << workerBits)
	sequenceBits   uint8 = 12
	sequenceMask   int64 = -1 ^ (-1 << sequenceBits)
	workerShift    uint8 = sequenceBits
	timestampShift uint8 = sequenceBits + workerBits
)

type IDWorker struct {
	sequence      int64
	lastTimestamp int64
	workerId      int64
	mutex         sync.Mutex
}

func NewIDWorker(workerId int64) (*IDWorker, error) {
	if workerId < 0 || workerId > maxWorker {
		return nil, fmt.Errorf("worker ID can't be greater than %d or less than 0", maxWorker)
	}

	return &IDWorker{
		workerId:      workerId,
		lastTimestamp: -1,
		sequence:      0,
	}, nil
}

func (w *IDWorker) NextID() (int64, error) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	timestamp := time.Now().UnixNano() / 1000000

	if timestamp < w.lastTimestamp {
		return 0, fmt.Errorf("clock moved backwards")
	}

	if timestamp == w.lastTimestamp {
		w.sequence = (w.sequence + 1) & sequenceMask
		if w.sequence == 0 {
			timestamp = w.nextMillisecond(timestamp)
		}
	} else {
		w.sequence = 0
	}

	w.lastTimestamp = timestamp

	return (timestamp << timestampShift) | (w.workerId << workerShift) | w.sequence, nil
}

func (w *IDWorker) nextMillisecond(currentTimestamp int64) int64 {
	for currentTimestamp <= w.lastTimestamp {
		currentTimestamp = time.Now().UnixNano() / 1000000
	}
	return currentTimestamp
}
