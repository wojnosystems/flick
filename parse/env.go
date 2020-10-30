package parse

import (
	env_parser "github.com/wojnosystems/go-env"
)

type env struct {
	parser env_parser.Env
}

func Env() EnvUnmarshaler {
	return &env{
		parser: env_parser.Env{
			ParseRegistry: defaultYamlParseRegistry,
		},
	}
}

func (e *env) Unmarshal(config interface{}) (err error) {
	return e.parser.Unmarshall(config)
}
