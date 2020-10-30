package parse

import "github.com/wojnosystems/go-optional"

type appConfig struct {
	Hostname optional.String   `yaml:"hostname"`
	Delay    optional.Duration `yaml:"delay"`
}
