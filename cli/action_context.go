package cli

import (
	"fmt"
	"os"
)

type actionContext struct {
	positionalArgs []string
}

func (c actionContext) Exit(code int) {
	os.Exit(code)
}

func (c actionContext) ExitWithMessage(code int, msg string) {
	if code == 0 {
		fmt.Sprintln(msg)
	} else {
		_, _ = fmt.Fprintln(os.Stderr, msg)
	}
	os.Exit(code)
}

func (c *actionContext) PositionalArgs() []string {
	return c.positionalArgs
}
