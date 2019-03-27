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
	"sync/atomic"
	"testing"
)

func TestWait(t *testing.T) {
	swg := New(10)
	var c uint32

	for i := 0; i < 10000; i++ {
		swg.Add()
		go func(c *uint32) {
			defer swg.Done()
			atomic.AddUint32(c, 1)
		}(&c)
	}

	swg.Wait()

	if c != 10000 {
		t.Fatalf("%d, not all routines have been executed.", c)
	}
}

func TestThrottling(t *testing.T) {
	var c uint32

	swg := New(4)

	if len(swg.current) != 0 {
		t.Fatalf("the SizedWaitGroup should start with zero.")
	}

	for i := 0; i < 10000; i++ {
		swg.Add()
		go func(c *uint32) {
			defer swg.Done()
			atomic.AddUint32(c, 1)
			if len(swg.current) > 4 {
				t.Fatalf("not the good amount of routines spawned.")
				return
			}
		}(&c)
	}

	swg.Wait()
}

func TestNoThrottling(t *testing.T) {
	var c uint32
	swg := New(0)
	if len(swg.current) != 0 {
		t.Fatalf("the SizedWaitGroup should start with zero.")
	}
	for i := 0; i < 10000; i++ {
		swg.Add()
		go func(c *uint32) {
			defer swg.Done()
			atomic.AddUint32(c, 1)
		}(&c)
	}
	swg.Wait()
	if c != 10000 {
		t.Fatalf("%d, not all routines have been executed.", c)
	}
}

func TestAddWithContext(t *testing.T) {
	ctx, cancelFunc := context.WithCancel(context.TODO())

	swg := New(1)

	if err := swg.AddWithContext(ctx); err != nil {
		t.Fatalf("AddContext returned error: %v", err)
	}

	cancelFunc()
	if err := swg.AddWithContext(ctx); err != context.Canceled {
		t.Fatalf("AddContext returned non-context.Canceled error: %v", err)
	}

}
