package cmd_definitions

import (
	"context"
	"github.com/wojnosystems/flick/pkg/generate/dsl"
)

type MethodHandler func(ctx context.Context) error
type ObjectFactory func() interface{}

type MethodDesc struct {
	Handler     MethodHandler
	Meta        dsl.Command
	ObjectMaker ObjectFactory
}
