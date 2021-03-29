package parse

import (
	envParser "github.com/wojnosystems/go-env/v2"
)

type env struct {
	parser *envParser.Env
}

func Env() EnvUnmarshaler {
	return &env{
		parser: envParser.NewWithParseRegistry(defaultYamlParseRegistry),
	}
}

func (e *env) Unmarshal(config interface{}) (err error) {
	return e.parser.Unmarshall(config)
}
