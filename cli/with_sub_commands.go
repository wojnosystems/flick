package cli

import (
	"github.com/wojnosystems/flick/command"
	"github.com/wojnosystems/flick/parse"
	"github.com/wojnosystems/flick/validate"
	flag_unmarshaler "github.com/wojnosystems/go-flag-unmarshaler"
)

type WithSubCommands struct {
	Name          string
	Usage         string
	Options       validate.Er
	OptionsParser []parse.Unmarshaler
	Commands      []command.Er
}

func (n *WithSubCommands) Run(groups []flag_unmarshaler.Group) (err error) {
	return
}
