package parse

import (
	flag_unmarshaler "github.com/wojnosystems/go-flag-unmarshaler"
)

func UnmarshalWithFile(defaultConfigFilePath string, fileUnmarshaler FileUnmarshaler, commandArgs []string, globalConfig interface{}) (flagGroups []flag_unmarshaler.Group, err error) {
	configFile := ConfigFile{}
	if "" != defaultConfigFilePath {
		configFile.ConfigFilePath.Set(defaultConfigFilePath)
	}
	flagGroups = flag_unmarshaler.Split(commandArgs)
	err = Unmarshall(&configFile,
		Env(),
		FlagsIgnoreUndefined(&flagGroups[0]))
	if err != nil {
		return
	}

	err = Unmarshall(globalConfig,
		FileIsOptional(configFile.ConfigFilePath, fileUnmarshaler),
		Env(),
		FlagsWithIgnore(&flagGroups[0], configFile.Flags()))
	return
}
