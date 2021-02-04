package parse

import "github.com/wojnosystems/go-optional/v2"

type ConfigFile struct {
	ConfigFilePath optional.String `env:"CONFIG_FILE_PATH" flag:"config-file-path" flag-short:"c" usage:"PATH" help:"file path to the configuration file that will set the global defaults"`
	ConfigTrace    optional.Bool   `env:"FLAG_TRACE" flag:"flagTrace" help:"true to show how all parameters are being set"`
}

func (f ConfigFile) Flags() []string {
	return []string{"-c", "--config-file-path", "--flagTrace"}
}
