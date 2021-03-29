package dsl

import "github.com/wojnosystems/go-optional/v2"

type Option struct {
	Type        string          `yaml:"type"`
	Description optional.String `yaml:"description"`
	Usage       optional.String `yaml:"usage"`
	Env         EnvDef          `yaml:"env"`
	Flag        FlagDef         `yaml:"flag"`
	Default     optional.String `yaml:"default"`
	Required    bool            `yaml:"required"`
}
