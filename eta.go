package eta

import (
	"fmt"
	"strconv"
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

// New creates a new ETA object, initiated with
// the specified total of items to be expected for
// completion.
func New(total int64) *ETA {

	clc := &ETA{
		total: total,
		arr:   fixedarr.New(int(total)),
	}

	return clc
}

// GetDone returns the number of items that
// have been completed.
func (eta *ETA) GetDone() int64 {
	done := atomic.LoadInt64(&eta.completed)
	if done > eta.GetTotal() {
		return eta.GetTotal()
	}
	return done
}

// GetTotal returns the total with which the ETA
// was initiated with.
func (eta *ETA) GetTotal() int64 {
	return eta.total
}

func getPercent(done int64, all int64) float64 {
	if all == 0 || done == 0 {
		return 0.0
	}
	percent := float64(100) / (float64(all) / float64(done))
	return percent
}
func getFormattedPercent(done int64, all int64) string {
	percentDone := fmt.Sprintf("%s%%", strconv.FormatFloat(getPercent(done, all), 'f', 2, 64))
	return percentDone
}

// GetFormattedPercentDone return the percent of completed work
// in the 12.34% format; example: 99.98%
func (eta *ETA) GetFormattedPercentDone() string {
	return getFormattedPercent(eta.GetDone(), eta.GetTotal())
}

// Done is used to signal that n items have been completed.
// Done returns the new number of completed items.
// If the new number of completed items is greater that the total,
// then the total is returned.
func (eta *ETA) Done(n int64) int64 {
	if n < 0 {
		panic("n is less that 0")
	}
	ts := time.Now().UnixNano()
	old := atomic.SwapInt64(&eta.lastCompletedTimestamp, ts)

	if old > 0 {
		eta.arr.Push(ts - old)
	}
	return atomic.AddInt64(&eta.completed, n)
}

// GetLastDoneTs gets the timestamp of the last done item.
func (eta *ETA) GetLastDoneTs() time.Time {
	nano := atomic.LoadInt64(&eta.lastCompletedTimestamp)

	return time.Unix(0, nano)
}

// GetETA gets the estimated time left to completion.
// NOTE: it might need a few items done to be able to
// compute an estimate.
func (eta *ETA) GetETA() time.Duration {
	vals := eta.arr.Value()
	if len(vals) == 0 {
		return time.Duration(0)
	}

	var totalDuration time.Duration
	for _, v := range vals {
		took := time.Duration(v.(int64))
		totalDuration += took
	}

	last := eta.GetLastDoneTs()
	durPerSingle := totalDuration / time.Duration(len(vals))
	total, completed := eta.GetTotal(), eta.GetDone()
	todo := total - completed

	timeDiff := time.Now().Sub(last)

	timeToGo := durPerSingle * time.Duration(todo)
	timeToGo -= timeDiff

	if timeToGo < 0 {
		return time.Duration(0)
	}
	return timeToGo
}
