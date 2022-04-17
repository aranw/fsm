package simple

import (
	"errors"
	"sync"
)

type state[S comparable] struct {
	transitions []S
}

type Transitions[S comparable] []S

// StateMap describes a State's possible Transitions
type StateMap[S comparable] struct {
	Name        S
	Transitions Transitions[S]
}

// States is a slice of States and it's possible Transitions
type States[S comparable] []StateMap[S]

type Simple[S comparable] struct {
	states map[S]state[S]

	// mu protects the values below
	mu      sync.Mutex
	current S
}

func NewSimple[S comparable](current S, states States[S]) *Simple[S] {
	simple := &Simple[S]{
		states:  make(map[S]state[S]),
		current: current,
		mu:      sync.Mutex{},
	}

	for _, s := range states {
		st := state[S]{
			transitions: make([]S, 0),
		}

		for _, t := range s.Transitions {
			st.transitions = append(st.transitions, t)
		}

		simple.states[s.Name] = st

	}

	return simple
}

// CanTransition checks to see whether the current State can be transitioned to another State
func (s *Simple[S]) CanTransition() (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	next, ok := s.states[s.current]
	if !ok {
		return false, errors.New("unknown state")
	}

	return len(next.transitions) > 0, nil
}

func (s *Simple[S]) Transition(to S) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	transitions := s.states[s.current].transitions
	if !contains(transitions, to) {
		return false, errors.New("cannot transition")
	}

	s.current = to

	return true, nil
}

func contains[S comparable](s []S, e S) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
