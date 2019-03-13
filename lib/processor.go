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
	SetStatus(Status)
	SignalChan() <-chan Signal
	Copy() Context
}

// Object describe data structs and could be compared by result of Id()
type Object interface {
	Id() string
}

// ProcessFunc is the type of function for ProcessorNode
type ProcessFunc func(Context, Object) (Object, error)

// ProcessorNode is the basic interface of nodes in process line
type ProcessorNode interface {
	// Context return context in process
	Context() Context

	// SetNextNode combine a process node to the end of current node
	// and return it
	SetNextNode(ProcessorNode) ProcessorNode

	// InputChan return a chan to receive Object for process
	InputChan() chan<- Object

	// OutputChan return a chan to provide Object for process,
	// If the ProcessNode has no output, the returned value is nil
	OutputChan() <-chan Object

	// Copy deeply copy the ProcessNode to create a new ProcessNode
	Copy() ProcessorNode

	// Run ProcessNode
	Run() error

	// Stop ProcessNode
	Stop() error
}

type baseNode struct {
	context Context
	next    ProcessorNode
	input   chan Object
	output  chan Object
}

// Context  context in current baseNode
func (bn *baseNode) Context() Context {
	return bn.context
}

// SetNextNode set a ProcessorNode as the next node of current baseNode
func (bn *baseNode) SetNextNode(node ProcessorNode) ProcessorNode {
	bn.next = node
	return node
}

// InputChan input channel in current baseNode
func (bn *baseNode) InputChan() chan<- Object {
	return bn.input
}

// OutputChan output channel in current baseNode
func (bn *baseNode) OutputChan() <-chan Object {
	return bn.output
}

// EndNext end the next node if it exists
func (bn *baseNode) EndNext() {
	if bn.next != nil {
		bn.next.Context().SetStatus(STOP_SIGNAL)
	}
}

// Copy  ...
func (bn *baseNode) Copy() *baseNode {
	clone := &baseNode{
		context: bn.context.Copy(),
		next:    nil,
		input:   make(chan Object, len(bn.input)),
		output:  make(chan Object, len(bn.output)),
	}

	if bn.next != nil {
		clone.next = bn.next.Copy()
	}

	return clone
}

// PipeNode nodes combined to pipe processor
type PipeNode interface {
	ProcessorNode
}

type pipeNode struct {
	*baseNode
	process ProcessFunc
}

// dealStopSignal  ...
func (pn *pipeNode) dealStopSignal() error {
	return nil
}

// Run ...
func (pn *pipeNode) Run() (err error) {
	select {
	case signal := <-pn.Context().SignalChan():
		if signal == STOP_SIGNAL {
			err = pn.dealStopSignal()
		}
		if err != nil {
			pn.EndNext()
			return
		}
	case obj := <-pn.baseNode.input:
		if pn.Context().Status() == RUNNING {
			output, err := pn.process(pn.Context(), obj)
			output = output
			// TODO maybe I could create a CatchNode to process error
			return err
		}
	}

	return nil
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
