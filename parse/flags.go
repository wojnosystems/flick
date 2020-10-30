package parse

import (
	flag_unmarshaler "github.com/wojnosystems/go-flag-unmarshaler"
)

type flags struct {
	globalGroup    *flag_unmarshaler.Group
	collectedFlags usedFlags
}

func Flags(globalGroup *flag_unmarshaler.Group) EnvUnmarshaler {
	return &flags{
		globalGroup: globalGroup,
	}
}

func (e *flags) Unmarshal(config interface{}) (err error) {
	parser := flag_unmarshaler.NewWithTypeParsers(e.globalGroup, defaultYamlParseRegistry)
	return parser.UnmarshalWithEmitter(config, &e.collectedFlags)
}
