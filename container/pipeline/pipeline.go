package pipeline

import (
	"container/heap"
	"fmt"
	"math"
)

// PipeLine interface
// every stage is also a pipeline
type PipeLine interface {
	heap.Interface // sort by pipeLevel
	fmt.Stringer   // print order by pipeLevel
	AddStage(parent string, name string, stage PipeLine) PipeLine
	GetStage(name string) PipeLine
	CustomizeStage(name string, stage PipeLine) PipeLine
	RemoveStage(name string) PipeLine
	Exec(input interface{}) (output interface{}, err error)
}

// PipeLevel level in pipe
type PipeLevel int8

// priority const
// lower will run before than higher
// you can customize your own priority
const (
	Debugger PipeLevel = iota // default 0:debug wrap
	Enhancer                  // enhancer feature
	Feature                   // user feature

	Initiator PipeLevel = math.MinInt8 // run first,default priority
	Runner    PipeLevel = math.MaxInt8 // final runner
)

// ClosePipeLevel by original pipeLevel
// if register close, flip around 0 would be useful
func ClosePipeLevel(p PipeLevel) PipeLevel {
	// fix -128 to 127
	if p == Initiator {
		p++
	}
	return -p
}
