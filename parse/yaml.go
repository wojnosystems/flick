package parse

import (
	optional_parse_registry "github.com/wojnosystems/go-optional-parse-registry/v2"
	parse_register "github.com/wojnosystems/go-parse-register"
	"github.com/wojnosystems/yamlreg"
	"io"
)

type yml struct {
	registry parse_register.ValueSetter
}

var defaultYamlParseRegistry = optional_parse_registry.NewWithGoPrimitives()

func Yaml() FileUnmarshaler {
	return YamlWithParseRegister(defaultYamlParseRegistry)
}

func YamlWithParseRegister(registry parse_register.ValueSetter) FileUnmarshaler {
	return &yml{
		registry: registry,
	}
}

func (y *yml) UnmarshalFile(r io.Reader, config interface{}) (err error) {
	dec := yamlreg.NewDecoder(r, y.registry)
	return dec.Decode(config)
}
