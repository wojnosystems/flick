package cli

import "github.com/wojnosystems/flick/action"

type Command struct {
	Name   string
	Usage  string
	Action action.Invoker
}
