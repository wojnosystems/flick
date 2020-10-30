package validate

import "github.com/wojnosystems/okey-dokey/bad"

type Validater interface {
	Validate(emitter bad.MemberEmitter) (err error)
}
