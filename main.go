package main

import (
	"fmt"

	"github.com/aranw/fsm/machine"
)

type Action struct {
	Name string
}

type State struct {
	Name string
}

func main() {
	Schedule := Action{
		Name: "schedule",
	}
	Request := Action{
		Name: "request",
	}
	Success := Action{
		Name: "success",
	}
	Fail := Action{
		Name: "fail",
	}
	Retry := Action{
		Name: "retry",
	}
	Discard := Action{
		Name: "discard",
	}
	New := State{
		Name: "new",
	}
	Scheduled := State{
		Name: "scheduled",
	}
	Processing := State{
		Name: "processing",
	}
	Successful := State{
		Name: "successful",
	}
	Discarded := State{
		Name: "discarded",
	}
	Failed := State{
		Name: "failed",
	}
	Retrying := State{
		Name: "retrying",
	}

	m := machine.NewMachine("delivery", machine.States[Action, State]{
		{
			Name: New,
			Transitions: []machine.Transition[Action, State]{
				{
					Input: Schedule,
					To:    Scheduled,
				},
			},
		},
		{
			Name: Scheduled,
			Transitions: []machine.Transition[Action, State]{
				{
					Input: Request,
					To:    Processing,
				},
			},
		},
		{
			Name: Processing,
			Transitions: []machine.Transition[Action, State]{
				{
					Input: Success,
					To:    Successful,
				},
				{
					Input: Fail,
					To:    Failed,
				},
				{
					Input: Discard,
					To:    Discarded,
				},
			},
		},
		{
			Name:        Successful,
			Transitions: []machine.Transition[Action, State]{},
		},
		{
			Name:        Discarded,
			Transitions: []machine.Transition[Action, State]{},
		},
		{
			Name: Failed,
			Transitions: []machine.Transition[Action, State]{
				{
					Input: Retry,
					To:    Retrying,
				},
				{
					Input: Discard,
					To:    Discarded,
				},
			},
		},
	})

	// Success
	fmt.Println(m.Transition(New, Schedule))
	fmt.Println(m.CanTransition(New))
	fmt.Println(m.Transition(Scheduled, Request))
	fmt.Println(m.CanTransition(Scheduled))
	fmt.Println(m.Transition(Processing, Success))
	fmt.Println(m.Transition(Processing, Fail))
	fmt.Println(m.Transition(Processing, Discard))
	fmt.Println(m.CanTransition(Processing))
	fmt.Println(m.Transition(Failed, Retry))
	fmt.Println(m.Transition(Failed, Discard))
	fmt.Println(m.CanTransition(Failed))
	fmt.Println(m.CanTransition(Discarded))

	// Failure
	fmt.Println(m.Transition(Scheduled, Success))
	fmt.Println(m.Transition(Successful, Discard))
	fmt.Println(m.Transition(Discarded, Success))
	fmt.Println(m.Transition(Failed, Success))
	fmt.Println(m.Transition(Retrying, Success))
	// fmt.Println(m.Transition(, Success))
}
