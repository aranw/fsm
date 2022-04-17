package simple_test

import (
	"fmt"
	"testing"

	"github.com/aranw/fsm/simple"
	"github.com/go-quicktest/qt"
)

// TurnstileState
type TurnstileState int

const (
	Locked TurnstileState = iota
	Unlocked
)

func (s TurnstileState) String() string {
	switch s {
	case Locked:
		return "Locked"
	case Unlocked:
		return "Unlocked"
	default:
		return fmt.Sprintf("Unknown state: %d", s)
	}
}

func TestMachine_Turnstile(t *testing.T) {
	t.Run("unable to lock locked turnstile", func(t *testing.T) {
		m := simple.NewSimple(Locked, simple.States[TurnstileState]{
			{Name: Locked, Transitions: simple.Transitions[TurnstileState]{Unlocked}},
			{Name: Unlocked, Transitions: simple.Transitions[TurnstileState]{Locked}},
		})

		_, err := m.Transition(Locked)
		qt.Assert(t, qt.IsNotNil(err))
	})

	t.Run("should be able to unlock turnstile", func(t *testing.T) {
		m := simple.NewSimple(Locked, simple.States[TurnstileState]{
			{Name: Locked, Transitions: simple.Transitions[TurnstileState]{Unlocked}},
			{Name: Unlocked, Transitions: simple.Transitions[TurnstileState]{Locked}},
		})

		ok, err := m.Transition(Unlocked)
		qt.Assert(t, qt.IsNil(err))
		qt.Assert(t, qt.IsTrue(ok))
	})

	t.Run("should be able to transition", func(t *testing.T) {
		m := simple.NewSimple(Locked, simple.States[TurnstileState]{
			{Name: Locked, Transitions: simple.Transitions[TurnstileState]{Unlocked}},
			{Name: Unlocked, Transitions: simple.Transitions[TurnstileState]{Locked}},
		})

		ok, err := m.CanTransition()
		qt.Assert(t, qt.IsNil(err))
		qt.Assert(t, qt.IsTrue(ok))
	})
}

type State int

const (
	New State = iota
	Starting
	Running
	Stopping
	Terminated
	Failed
)

func (s State) String() string {
	switch s {
	case New:
		return "New"
	case Starting:
		return "Starting"
	case Running:
		return "Running"
	case Stopping:
		return "Stopping"
	case Terminated:
		return "Terminated"
	case Failed:
		return "Failed"
	default:
		return fmt.Sprintf("Unknown state: %d", s)
	}
}

func TestSimple_Transition(t *testing.T) {
	s := simple.NewSimple(New, simple.States[State]{
		{
			Name: New,
			Transitions: simple.Transitions[State]{
				Starting,
				Terminated,
			},
		},
		{
			Name: Starting,
			Transitions: simple.Transitions[State]{
				Running,
				Failed,
			},
		},
		{
			Name: Running,
			Transitions: simple.Transitions[State]{
				Stopping,
			},
		},
		{
			Name: Stopping,
			Transitions: simple.Transitions[State]{
				Terminated,
			},
		},
		{
			Name:        Terminated,
			Transitions: simple.Transitions[State]{},
		},
		{
			Name:        Failed,
			Transitions: simple.Transitions[State]{},
		},
	})

	successCheck := func(to State) func(t *testing.T) {
		return func(t *testing.T) {
			got, err := s.Transition(to)
			qt.Assert(t, qt.IsNil(err))
			qt.Assert(t, qt.IsTrue(got))
		}
	}

	failCheck := func(to State) func(t *testing.T) {
		return func(t *testing.T) {
			got, err := s.Transition(to)
			qt.Assert(t, qt.IsNotNil(err))
			qt.Assert(t, qt.IsFalse(got))
		}
	}

	// success
	t.Run("new->starting", successCheck(Starting))
	t.Run("starting->running", successCheck(Running))
	t.Run("running->stopping", successCheck(Stopping))
	t.Run("stopping->terminated", successCheck(Terminated))

	// should not be able to transition to any states when Terminated

	// fail
	t.Run("terminated->new", failCheck(New))
	t.Run("terminated->starting", failCheck(Starting))
	t.Run("terminated->running", failCheck(Running))
	t.Run("terminated->terminated", failCheck(Terminated))
	t.Run("terminated->failed", failCheck(Failed))
}
