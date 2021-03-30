package cli

import "errors"

var ErrUnimplemented = errors.New("command was declared, but not implemented")
