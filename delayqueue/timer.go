package delayqueue

import (
	"time"
)

// Timer for provide the current ms
type Timer interface {
	Now() int64
}

// the default Timer implement
type timer struct{}

// Now implement Timer.Now
func (t timer) Now() int64 {
	return int64(time.Duration(time.Now().UnixNano()) / time.Millisecond)
}

var (
	defaultTimer = timer{}
)
