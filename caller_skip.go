package logger

import "sync/atomic"

// CallerSkip is a struct that holds the number of frames to skip
type CallerSkip struct {
	skip *atomic.Int32
}

func NewCallerSkip(skip int) CallerSkip {
	cs := CallerSkip{skip: new(atomic.Int32)}
	cs.Set(skip)
	return cs
}

// Set sets the number of frames to skip
func (c CallerSkip) Set(skip int) {
	c.skip.Store(int32(skip))
}

// Load loads the number of frames to skip
func (c CallerSkip) Load() int {
	return int(c.skip.Load())
}
