package machine

import (
	"errors"
)

type state[A, S comparable] struct {
	transitions map[A]S
}

// Transition describes what Input Action can transition a State From and To
// An optional Output Action can also be set
type Transition[A, S comparable] struct {
	Input A
	To    S
}

type Transitions[A, S comparable] []Transition[A, S]

// StateMap describes a State's possible Transitions
type StateMap[A, S comparable] struct {
	Name        S
	Transitions Transitions[A, S]
}

// States is a slice of States and it's possible Transitions
type States[A, S comparable] []StateMap[A, S]

type Machine[A, S comparable] struct {
	name   string
	states map[S]state[A, S]
}

func NewMachine[A, S comparable](name string, states States[A, S]) *Machine[A, S] {
	machine := &Machine[A, S]{
		name:   name,
		states: make(map[S]state[A, S]),
	}

	for _, s := range states {
		st := state[A, S]{
			transitions: make(map[A]S),
		}

		for _, t := range s.Transitions {
			st.transitions[t.Input] = t.To
		}

		machine.states[s.Name] = st
	}

	return machine
}

func (m *Machine[A, S]) next(current S, action A) (S, error) {
	next, ok := m.states[current].transitions[action]
	if !ok {
		return *new(S), errors.New("cannot transition")
	}

	return next, nil
}

func (m *Machine[A, S]) Transition(state S, action A) (S, error) {
	return m.next(state, action)
}

func (m *Machine[A, S]) CanTransition(state S) (bool, error) {
	next, ok := m.states[state]
	if !ok {
		return false, errors.New("unknown state")
	}

	return len(next.transitions) > 0, nil
}
