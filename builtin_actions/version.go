package builtin_actions

import (
	"github.com/wojnosystems/flick/cli"
	"github.com/wojnosystems/flick/context"
	"github.com/wojnosystems/flick/invoke"
)

func CommandPrintVersion(v string) *cli.Action {
	return &cli.Action{
		Name: "version",
		Action: &invoke.WithoutOptionsOrErrors{
			Action: func(c context.Context) {
				c.ExitWithMessage(0, v)
			},
		},
	}
}
