package parse

import "github.com/wojnosystems/go-optional/v2"

type appConfig struct {
	Hostname optional.String   `yaml:"hostname"`
	Delay    optional.Duration `yaml:"delay"`
}
