// Copyright © 2019 Mephis Pheies <mephistommm@gmail.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.
package concurrent

import (
	"context"
	"math"
	"strconv"
	"sync"
	"sync/atomic"
)

// TimeoutCond is a sync.Cond  improve for support wait timeout.
type TimeoutCond struct {
	L          sync.Locker
	signal     chan int
	hasWaiters uint64
}

// NewTimeoutCond return a new TimeoutCond
func NewTimeoutCond(l sync.Locker) *TimeoutCond {
	cond := TimeoutCond{L: l, signal: make(chan int, 0)}
	return &cond
}

func (cond *TimeoutCond) addWaiter() {
	v := atomic.AddUint64(&cond.hasWaiters, 1)
	if v == 0 {
		panic("too many waiters; max is " + strconv.FormatUint(math.MaxUint64, 10))
	}
}

func (cond *TimeoutCond) removeWaiter() {
	// Decrement. See notes here: https://godoc.org/sync/atomic#AddUint64
	v := atomic.AddUint64(&cond.hasWaiters, ^uint64(0))

	if v == math.MaxUint64 {
		panic("removeWaiter called more than once after addWaiter")
	}
}

// HasWaiters queries whether any goroutine are waiting on this condition
func (cond *TimeoutCond) HasWaiters() bool {
	return atomic.LoadUint64(&cond.hasWaiters) > 0
}

// Wait waits for a signal, or for the context do be done. Returns true if signaled.
func (cond *TimeoutCond) Wait(ctx context.Context) bool {
	cond.addWaiter()
	//copy signal in lock, avoid data race with Interrupt
	ch := cond.signal
	//wait should unlock mutex,  if not will cause deadlock
	cond.L.Unlock()
	defer cond.removeWaiter()
	defer cond.L.Lock()

	select {
	case _, ok := <-ch:
		return !ok
	case <-ctx.Done():
		return false
	}
}

// Signal wakes one goroutine waiting on c, if there is any.
func (cond *TimeoutCond) Signal() {
	select {
	case cond.signal <- 1:
	default:
	}
}

// Interrupt goroutine wait on this TimeoutCond
func (cond *TimeoutCond) Interrupt() {
	cond.L.Lock()
	defer cond.L.Unlock()
	close(cond.signal)
	cond.signal = make(chan int, 0)
}
