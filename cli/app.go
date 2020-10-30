package cli

import flag_unmarshaler "github.com/wojnosystems/go-flag-unmarshaler"

type New struct {
	Name     string
	Usage    string
	Commands []Command
}

func (n *New) Run(groups []flag_unmarshaler.Group) (err error) {
	return
}
