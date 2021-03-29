package dsl

import "github.com/wojnosystems/okey-dokey/bad"

type trackedEmitter struct {
	hasErrors bool
	wrapped   bad.MemberEmitter
	parent    *trackedEmitter
}

func newTrackedEmitter(emitter bad.MemberEmitter) *trackedEmitter {
	return &trackedEmitter{
		hasErrors: false,
		wrapped:   emitter,
	}
}

func (t trackedEmitter) isInvalid() bool {
	return t.hasErrors
}

func (t *trackedEmitter) Emit(msg string) {
	if t.parent == nil {
		t.hasErrors = true
	} else {
		t.parent.hasErrors = true
	}
	t.wrapped.Emit(msg)
}

func (t *trackedEmitter) Into(field string) bad.MemberEmitter {
	p := t.parent
	if p == nil {
		p = t
	}
	return &trackedEmitter{
		wrapped: t.wrapped.Into(field),
		parent:  p,
	}
}
