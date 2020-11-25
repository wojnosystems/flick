package invoke

import (
	"github.com/wojnosystems/flick/context"
)

type Actioner interface {
	Invoke(context context.Action) (err error)
}

type SubCommander interface {
	Invoke(context context.SubCommand) (err error)
}
