package parse

import "github.com/wojnosystems/go-optional/v2"

// ConfigFile is a convenience structure that represents a configuration file for your application. Just embed this into your configuration file struct
// to get instant support for a -c config file path
type ConfigFile struct {
	ConfigFilePath optional.String `env:"CONFIG_FILE_PATH" flag:"config-file-path" flag-short:"c" usage:"PATH" help:"file path to the configuration file that will set the global defaults"`
}

func (f ConfigFile) Flags() []string {
	return []string{"-c", "--config-file-path"}
}
