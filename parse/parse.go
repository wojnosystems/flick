package parse

import (
	"errors"
	"github.com/wojnosystems/flick/pkg/cmd_definitions"
	env_parser "github.com/wojnosystems/go-env/v2"
	flag_unmarshaler "github.com/wojnosystems/go-flag-unmarshaler"
	"strings"
)

type EnvFlagParser struct {
	service   cmd_definitions.ServiceDesc
	envParser env_parser.Env
	args      []string
}

func (e *EnvFlagParser) Parse(callback func(path []string)) (exec Exec, err error) {
	path := make([]string, 0, 10)

	commands := flag_unmarshaler.Split(e.args)
	var commandItem cmd_definitions.MethodDesc
	for _, command := range commands {
		path = append(path, command.CommandName)
		callback(path)
		var ok bool
		commandItem, ok = e.service.Methods.Get(path...)
		if !ok {
			err = errors.New("unsupported command: " + strings.Join(path, " "))
			return
		}
		if commandItem.ObjectMaker != nil {
			config := commandItem.ObjectMaker()
			exec.configObjects = append(exec.configObjects, config)
		}
	}
	err = e.envParser.Unmarshall(&config)
	if err != nil {
		return
	}
	err = flag_unmarshaler.New(&command).Unmarshal(&config)
	if err != nil {
		return
	}
	return
}

type Exec struct {
	paths         []string
	configObjects []interface{}
}
