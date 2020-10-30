package parse

import "io"

type Unmarshaler interface {
	Unmarshal(config interface{}) (err error)
}

type FileUnmarshaler interface {
	UnmarshalFile(r io.Reader, config interface{}) (err error)
}

type FlagUnmarshaler interface {
	Unmarshaler
}

type EnvUnmarshaler interface {
	Unmarshaler
}
