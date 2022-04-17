package machine

import (
	"errors"
	"sync"
)

type state[E, S comparable] struct {
	transitions map[E]S
}

// Transition describes what Input Event can transition a State the current State and To
type Transition[E, S comparable] struct {
	Event E
	To    S
}

type Transitions[E, S comparable] []Transition[E, S]

// StateMap describes a State's possible Transitions
type StateMap[E, S comparable] struct {
	Name        S
	Transitions Transitions[E, S]
}

// States is a slice of States and it's possible Transitions
type States[E, S comparable] []StateMap[E, S]

type Machine[E, S comparable] struct {
	states map[S]state[E, S]

	// mu protects the values below
	mu      sync.Mutex
	current S
}

// NewMachine creates a new Finite State Machine
func NewMachine[E, S comparable](current S, states States[E, S]) *Machine[E, S] {
	machine := &Machine[E, S]{
		states:  make(map[S]state[E, S]),
		current: current,
		mu:      sync.Mutex{},
	}

	for _, s := range states {
		st := state[E, S]{
			transitions: make(map[E]S),
		}

		for _, t := range s.Transitions {
			st.transitions[t.Event] = t.To
		}

		machine.states[s.Name] = st
	}

	return machine
}

// Transition takes the provided event and attempts to transition to the next State
// If successfully it'll return the new State
// If the current state can not be transitioned an error will be returned
func (m *Machine[E, S]) Transition(event E) (S, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	next, ok := m.states[m.current].transitions[event]
	if !ok {
		return *new(S), errors.New("cannot transition")
	}

	m.current = next

	return next, nil
}

// CanTransition checks to see whether the current State can be transitioned to another State
func (m *Machine[A, S]) CanTransition() (bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	next, ok := m.states[m.current]
	if !ok {
		return false, errors.New("unknown state")
	}

	return len(next.transitions) > 0, nil
}

// State returns the current state of the machine
func (m *Machine[E, S]) State() S {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.current
}
