package cli

import (
	"github.com/wojnosystems/flick/invoke"
	"github.com/wojnosystems/flick/parse"
	"github.com/wojnosystems/flick/validate"
	flag_unmarshaler "github.com/wojnosystems/go-flag-unmarshaler"
	"github.com/wojnosystems/okey-dokey/bad"
)

type WithAction struct {
	Name    string
	Usage   string
	Options validate.Er
	Parsing parse.OptionUnmarshaler
	Action  invoke.Actioner
}

func (n *WithAction) Run(groups []flag_unmarshaler.Group) (err error) {
	if n.Parsing != nil {
		err = n.Parsing.UnmarshalRoot(n.Options, groups[0])
		if err != nil {
			return
		}
		validationErrors := bad.NewCollection()
		if n.Options != nil {
			err = n.Options.Validate(validationErrors)
			if err != nil {
				return
			}
		}
		if validationErrors.HasAny() {
			// TODO print the validation errors
			return
		}
	}
	ctx := &actionContext{}
	for _, group := range groups {
		ctx.positionalArgs = append(ctx.positionalArgs, group.CommandName)
	}
	return n.Action.Invoke(ctx)
}
