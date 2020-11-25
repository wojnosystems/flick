package command

import "github.com/wojnosystems/flick/invoke"

type Action struct {
	Name    string
	Usage   string
	Options interface{}
	Action  invoke.Actioner
}
