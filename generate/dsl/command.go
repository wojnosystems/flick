package dsl

import (
	"github.com/wojnosystems/go-optional/v2"
	"github.com/wojnosystems/okey-dokey/bad"
)

type Command struct {
	Usage       optional.String     `yaml:"usage"`
	Description optional.String     `yaml:"description"`
	Options     []OptionOrReference `yaml:"options"`
	Commands    NamedCommands       `yaml:"commands"`
	MinArgs     uint                `yaml:"minArgs"`
	MaxArgs     uint                `yaml:"maxArgs"`
}

var commandValidations = commandValidationDefs{}

type commandValidationDefs struct {
}

func (d commandValidationDefs) Validate(on *Command, emitter bad.MemberEmitter) {
	validateMinMaxArgs(on.MinArgs, on.MaxArgs, emitter)
	validateMaxArgsWithSubCommands(on.MaxArgs, on.Commands, emitter)
	for commandName, command := range on.Commands {
		commandValidations.Validate(&command, emitter.Into(commandName))
	}
}

func validateMinMaxArgs(minArgs, maxArgs uint, emitter bad.Emitter) {
	if minArgs > maxArgs {
		emitter.Emit("minArgs must be less than maxArgs")
	}
}

func validateMaxArgsWithSubCommands(maxArgs uint, commands NamedCommands, emitter bad.MemberEmitter) {
	if commands.HasAny() && maxArgs != 0 {
		emitter.Emit("when sub-commands are specified, maxArgs must be 0")
	}
}
