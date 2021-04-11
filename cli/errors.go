package cli

import "errors"

var ErrCommandUnimplemented = errors.New("command was declared, but not implemented")
