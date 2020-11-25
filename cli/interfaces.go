package cli

import flag_unmarshaler "github.com/wojnosystems/go-flag-unmarshaler"

type Runner interface {
	Run(groups []flag_unmarshaler.Group) (err error)
}
