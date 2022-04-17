package machine_test

import (
	"fmt"
	"testing"

	"github.com/aranw/fsm/machine"
	"github.com/go-quicktest/qt"
)

// TurnstileEvent
type TurnstileEvent int

const (
	Coin TurnstileEvent = iota
	Push
)

func (e TurnstileEvent) String() string {
	switch e {
	case Push:
		return "Push"
	case Coin:
		return "Coin"
	default:
		return fmt.Sprintf("Unknown event: %d", e)
	}
}

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
		m := machine.NewMachine(Locked, machine.States[TurnstileEvent, TurnstileState]{
			{Name: Locked, Transitions: []machine.Transition[TurnstileEvent, TurnstileState]{{Event: Coin, To: Unlocked}}},
			{Name: Unlocked, Transitions: []machine.Transition[TurnstileEvent, TurnstileState]{{Event: Push, To: Locked}}},
		})

		_, err := m.Transition(Push)
		qt.Assert(t, qt.IsNotNil(err))
	})

	t.Run("should be able to unlock turnstile", func(t *testing.T) {
		m := machine.NewMachine(Locked, machine.States[TurnstileEvent, TurnstileState]{
			{Name: Locked, Transitions: []machine.Transition[TurnstileEvent, TurnstileState]{{Event: Coin, To: Unlocked}}},
			{Name: Unlocked, Transitions: []machine.Transition[TurnstileEvent, TurnstileState]{{Event: Push, To: Locked}}},
		})

		s, err := m.Transition(Coin)
		qt.Assert(t, qt.IsNil(err))
		qt.Assert(t, qt.Equals(s, Unlocked))
	})

	t.Run("should be able to transition", func(t *testing.T) {
		m := machine.NewMachine(Locked, machine.States[TurnstileEvent, TurnstileState]{
			{Name: Locked, Transitions: []machine.Transition[TurnstileEvent, TurnstileState]{{Event: Coin, To: Unlocked}}},
			{Name: Unlocked, Transitions: []machine.Transition[TurnstileEvent, TurnstileState]{{Event: Push, To: Locked}}},
		})

		ok, err := m.CanTransition()
		qt.Assert(t, qt.IsNil(err))
		qt.Assert(t, qt.IsTrue(ok))
	})
}

type ServiceEvent int

const (
	Start ServiceEvent = iota
	Stop
)

type ServiceState int

const (
	New ServiceState = iota
	Started
	Running
	Stopped
	Terminated
	Failed
)

func (s ServiceState) String() string {
	switch s {
	case New:
		return "New"
	case Started:
		return "Started"
	case Running:
		return "Running"
	case Stopped:
		return "Stopped"
	case Terminated:
		return "Terminated"
	case Failed:
		return "Failed"
	default:
		return fmt.Sprintf("Unknown state: %d", s)
	}
}

func TestMachine_Complex(t *testing.T) {

}
