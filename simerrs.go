package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

type state int

const (
	Start = iota
	Err
	OK
)

type stateMachine struct {
	s  state
	ne uint8
}

func (s *stateMachine) Advance() {
	switch s.s {
	case Start:
		s.s = Err
		s.ne = 1
	case Err:
		s.ne++
		if int(s.ne) == *ErrCnt {
			s.s = OK
		}
	case OK:
		s.s = Start
	}
}

type simErrs struct {
	mu    sync.Mutex
	sm    stateMachine
	start time.Time
}

func (s *simErrs) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var status int
	s.mu.Lock()
	defer s.mu.Unlock()
	switch s.sm.s {
	case Start:
		s.start = time.Now()
		fallthrough
	case Err:
		status = http.StatusInternalServerError
	case OK:
		status = http.StatusOK
	}
	fmt.Printf("%s Received at t=(%d) | Respond: %d\n", r.Method, time.Since(s.start).Milliseconds(), status)
	s.sm.Advance()
	w.WriteHeader(status)
}
