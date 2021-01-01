package pipeline

import (
	"container/heap"
	"math"
	"strings"
)

// SimpleFunc stage fn
type SimpleFunc func() error

// SimpleFuncStage simple func stage
type SimpleFuncStage struct {
	fn         SimpleFunc
	priority   PipeLevel
	stages     []*SimpleFuncStage          // children
	stageMap   map[string]*SimpleFuncStage // can reach stage after pop
	defaultPPL PipeLevel                   // default parent pipe level
	// should not changed after init
	name  string
	split string
}

// NewSimpleFuncStage empty root simple func pipe
func NewSimpleFuncStage(name, split string) *SimpleFuncStage {
	return &SimpleFuncStage{
		name:     name,
		split:    split,
		fn:       func() error { return nil },
		priority: Debugger,
		stageMap: make(map[string]*SimpleFuncStage),
	}
}

// SetFunc set func
func (s *SimpleFuncStage) SetFunc(f SimpleFunc) *SimpleFuncStage {
	s.fn = f
	return s
}

// SetPriority set priority
func (s *SimpleFuncStage) SetPriority(p PipeLevel) *SimpleFuncStage {
	s.priority = p
	return s
}

// SetDefaultPPL set default priority for parent(auto added)
func (s *SimpleFuncStage) SetDefaultPPL(p PipeLevel) *SimpleFuncStage {
	s.defaultPPL = p
	return s
}

// AddStage add stage to pipe or offspring of this pipe
func (s *SimpleFuncStage) AddStage(parent string, name string, stage PipeLine) PipeLine {
	if !strings.HasPrefix(name, s.name) {
		return nil
	}

	// 1. s is parent
	if s.name == parent {
		funcStage, ok := stage.(*SimpleFuncStage)
		if !ok {
			return nil
		}

		funcStage.name = name

		// avoid repeat add
		if p := s.stageMap[name]; p != nil {
			return p
		}

		heap.Push(s, funcStage)
		s.stageMap[name] = funcStage
		return funcStage
	}

	// 2. s is not parent
	// map direct
	if p := s.stageMap[parent]; p != nil {
		return p.AddStage(parent, name, stage)
	}
	// slow add
	for _, v := range s.stages {
		// 2.1 parent is a children of v
		if strings.HasPrefix(parent, v.name) {
			return v.AddStage(parent, name, stage)
		}

		// 2.2 origin v should be a children of parent(new one)
		if strings.HasPrefix(v.name, parent+s.split) {
			origin := *v

			// change to parent
			v.name = parent
			v.fn = func() error { return nil }
			v.stages = nil
			v.stageMap = map[string]*SimpleFuncStage{}
			delete(s.stageMap, origin.name)
			s.stageMap[parent] = v

			if v.AddStage(parent, origin.name, &origin) == nil {
				// rollback
				v = &origin
				delete(s.stageMap, parent)
				s.stageMap[origin.name] = v
				return nil
			}

			n, ok := stage.(*SimpleFuncStage)
			if !ok {
				return nil
			}

			return v.AddStage(parent, name, n)
		}
	}

	// 3. no parent before
	// this path is very dangerous because of parent's default PipeLevel
	// It is necessary to confirm parent's pipeLevel by set child's defaultPPL
	p := NewSimpleFuncStage(parent, s.split).SetPriority(s.defaultPPL)
	heap.Push(s, p)
	s.stageMap[parent] = p
	return p.AddStage(parent, name, stage)
}

// GetStage get stage with name
func (s *SimpleFuncStage) GetStage(name string) PipeLine {
	if !strings.HasPrefix(name, s.name) {
		return nil
	}

	// map direct
	if p := s.stageMap[name]; p != nil {
		return p
	}
	// slow get
	for _, v := range s.stages {
		if name == v.name {
			return v
		}

		if strings.HasPrefix(name, v.name+s.split) {
			return v.GetStage(name)
		}
	}

	return nil
}

// CustomizeStage customize stage with new stage
func (s *SimpleFuncStage) CustomizeStage(name string, stage PipeLine) PipeLine {
	if !strings.HasPrefix(name, s.name) {
		return nil
	}

	if s.name == name {
		sfs, ok := stage.(*SimpleFuncStage)
		if !ok {
			return nil
		}

		s.fn = sfs.fn
		s.priority = sfs.priority
		return s
	}

	// map direct
	if p := s.stageMap[name]; p != nil {
		return p.CustomizeStage(name, stage)
	}
	// slow customize
	for _, v := range s.stages {
		if strings.HasPrefix(name, v.name) {
			return v.CustomizeStage(name, stage)
		}
	}

	return nil
}

// RemoveStage remove stage with name
func (s *SimpleFuncStage) RemoveStage(name string) PipeLine {
	if !strings.HasPrefix(name, s.name) {
		return nil
	}

	for i, v := range s.stages {
		if name == v.name {
			heap.Remove(s, i)
			return v
		}

		if strings.HasPrefix(name, v.name+s.split) {
			return v.RemoveStage(name)
		}
	}

	return nil
}

// Exec finally
func (s *SimpleFuncStage) Exec(input interface{}) (output interface{}, err error) {
	err = s.fn()
	if err != nil {
		return input, err
	}

	for s.Len() != 0 {
		input, err = heap.Pop(s).(*SimpleFuncStage).Exec(input)
		if err != nil {
			return input, err
		}
	}

	return input, err
}

// String help debug print exec pipeline
func (s *SimpleFuncStage) String() string {
	p := s.name

	spi := " & "
	pipes := make([]string, 2*(math.MaxInt8+1))
	for _, stage := range s.stages {
		c := stage.String()
		if c != "" {
			pipes[uint8(stage.priority)+math.MaxInt8+1] += c + spi
		}
	}

	for _, pp := range pipes {
		if pp != "" {
			p += "\n => " + strings.TrimSuffix(pp, spi)
		}
	}
	return p
}

// Following functions should not be used directly

func (s *SimpleFuncStage) Len() int {
	return len(s.stages)
}

func (s *SimpleFuncStage) Less(i, j int) bool {
	return s.stages[i].priority < s.stages[j].priority
}

func (s *SimpleFuncStage) Swap(i, j int) {
	s.stages[i], s.stages[j] = s.stages[j], s.stages[i]
}

func (s *SimpleFuncStage) Push(x interface{}) {
	s.stages = append(s.stages, x.(*SimpleFuncStage))
}

func (s *SimpleFuncStage) Pop() interface{} {
	n := s.Len()
	item := s.stages[n-1]
	s.stages[n-1] = nil // avoid memory leak
	s.stages = s.stages[0 : n-1]
	return item
}
