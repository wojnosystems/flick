package validate

import "github.com/wojnosystems/okey-dokey/bad"

type Er interface {
	Validate(emitter bad.MemberEmitter) (err error)
}
