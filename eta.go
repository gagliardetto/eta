package eta

import (
	"sync/atomic"
	"time"

	"github.com/gagliardetto/fixedarr"
)

type ETA struct {
	total     int64
	completed int64

	arr                    *fixedarr.Array
	lastCompletedTimestamp int64
}

func New(total int64, frames ...time.Duration) *ETA {
	if total < 1 {
		panic("total is less that 1")
	}
	if len(frames) == 0 {
		panic("not time frames provided")
	}

	clc := &ETA{
		total: total,
		arr:   fixedarr.New(int(total)),
	}

	return clc
}

func (eta *ETA) GetDone() int64 {
	done := atomic.LoadInt64(&eta.completed)
	if done > eta.GetTotal() {
		return eta.GetTotal()
	}
	return done
}
func (eta *ETA) GetTotal() int64 {
	return eta.total
}
func (eta *ETA) Done(n int64) int64 {
	ts := time.Now().UnixNano()
	old := atomic.SwapInt64(&eta.lastCompletedTimestamp, ts)

	if old > 0 {
		eta.arr.Push(ts - old)
	}
	return atomic.AddInt64(&eta.completed, n)
}
func (eta *ETA) GetLastDoneTs() time.Time {
	nano := atomic.LoadInt64(&eta.lastCompletedTimestamp)

	return time.Unix(0, nano)
}

func (etac *ETA) GetETA() time.Duration {
	vals := etac.arr.Value()
	if len(vals) == 0 {
		return time.Duration(0)
	}

	var totalDuration time.Duration
	for _, v := range vals {
		took := time.Duration(v.(int64))
		totalDuration += took
	}

	last := etac.GetLastDoneTs()
	durPerSingle := totalDuration / time.Duration(len(vals))
	total, completed := etac.GetTotal(), etac.GetDone()
	todo := total - completed

	timeDiff := time.Now().Sub(last)

	timeToGo := durPerSingle * time.Duration(todo)
	timeToGo -= timeDiff

	if timeToGo < 0 {
		return time.Duration(0)
	}
	return timeToGo
}
