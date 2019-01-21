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
package lib

// Status is the status of a process node
type Status uint8

const (
	READY Status = 1 + iota
	RUNNING
	STOPPING
	STOPPED
	FINISHED
)

// Signal is send to signal chan of a context
type Signal uint8

const (
	STOP_SIGNAL = 0
)

// Context is the context for a process node
type Context interface {
	Status() Status
	SignalChan() chan<- Signal
}

// Object describe data structs and could be compared by result of Id()
type Object interface {
	Id() string
}

// ProcessorNode is the basic interface of nodes in process line
type ProcessorNode interface {
	// Context return context in process
	Context() Context

	// SetNextNode combine a process node to the end of current node
	// and return it
	SetNextNode() ProcessNode

	// InputChan return a chan to receive Object for process
	InputChan() chan<- Object

	// OutputChan return a chan to provide Object for process,
	// If the ProcessNode has no output, the returned value is nil
	OutputChan() <-chan Object

	// Copy deeply copy the ProcessNode to create a new ProcessNode
	Copy() ProcessNode

	// Run ProcessNode
	Run() error

	// Stop ProcessNode
	Stop() error
}

// PipeNode nodes combined to pipe processor
type PipeNode interface {
	ProcessorNode
}

// BifurcateProcessorNode nodes combined to bifurcate processor
type BifurcateProcessorNode interface {
	ProcessorNode

	// SetProcessorNum config number of bifurcate processor
	SetProcesorNum(num int)
}

// AOVNode nodes combined to aov processor
type AOVNode interface {
	ProcessorNode

	// SetProcessorNum config number of processors in aov
	SetProcesorNum(num int)
}

// ConditionNode nodes combined to condition processor
type ConditionNode interface {
	ProcessorNode

	// SetOKProcessor combine a processor run while condition is ok
	SetOKProcessor(p ProcessorNode) ProcessorNode

	// SetOKProcessor combine a processor run while condition is nil
	SetNilProcessor(p ProcessorNode) ProcessorNode
}
