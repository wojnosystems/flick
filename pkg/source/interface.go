package source

import "github.com/wojnosystems/okey-dokey/bad"

type Parser interface {
	Unmarshall(into interface{}, emitter bad.Emitter) (err error)
}
