package parse

import (
	flag_unmarshaler "github.com/wojnosystems/go-flag-unmarshaler"
	"io"
)

type Unmarshaler interface {
	Unmarshal(config interface{}) (err error)
}

type FileUnmarshaler interface {
	UnmarshalFile(r io.Reader, config interface{}) (err error)
}

type FlagUnmarshaler interface {
	Unmarshal(config interface{}, group flag_unmarshaler.Group) (err error)
}

type EnvUnmarshaler interface {
	Unmarshaler
}
