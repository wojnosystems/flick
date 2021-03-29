package cli

import (
	"context"
	env_parser "github.com/wojnosystems/go-env/v2"
)

type Commander interface {
	Switch(ctx context.Context, args []string, receiver env_parser.EnvReader) (err error)
}
