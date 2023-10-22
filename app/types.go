package app

import "sync"

type State struct {
	I         int
	Result    float64
	WorkerNum int
}

type Command struct {
	Type  string
	Value int
}

type Stage struct {
	value int
	mx    sync.Mutex
}

func (s *Stage) Set(n int) {
	s.mx.Lock()
	s.value = n
	s.mx.Unlock()
}

func (s *Stage) Get() int {
	return s.value
}
