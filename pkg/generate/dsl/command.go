package dsl

import (
	"github.com/wojnosystems/go-optional/v2"
	"github.com/wojnosystems/okey-dokey/bad"
)

type Command struct {
	// Usage explains how to use this option in a one-liner
	// e.g. --name=Bob
	Usage optional.String `yaml:"usage"`

	// Description is a longer version of usage, include multi-line descriptions here
	Description optional.String `yaml:"description"`

	// Options defines what options this command will parse
	Options []OptionOrReference `yaml:"options"`

	// Commands are the sub-commands of this command
	Commands NamedCommands `yaml:"commands"`

	// MinArgs is the minimum number of arguments that this command takes,
	// this is incompatible with Commands, as commands are "arguments" and is ignored when Commands is not empty
	MinArgs uint `yaml:"minArgs"`

	// MaxArgs is the maximum number of arguments that this command takes,
	// this is incompatible with Commands, as commands are "arguments" and is ignored when Commands is not empty
	MaxArgs uint `yaml:"maxArgs"`
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
