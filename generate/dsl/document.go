package dsl

import "github.com/wojnosystems/okey-dokey/bad"

type Document struct {
	OptionApi  OptionApi                `yaml:"optionapi"`
	Options    NamedOptionsOrReferences `yaml:"options"`
	Commands   NamedCommands            `yaml:"commands"`
	Components Components               `yaml:"components"`
	MinArgs    uint                     `yaml:"minArgs"`
	MaxArgs    uint                     `yaml:"maxArgs"`
}

var DocumentValidations = DocumentValidationDefs{}

type DocumentValidationDefs struct {
}

func (d DocumentValidationDefs) Validate(on *Document, emitter bad.MemberEmitter) {
	validateMinMaxArgs(on.MinArgs, on.MaxArgs, emitter)
	validateMaxArgsWithSubCommands(on.MaxArgs, on.Commands, emitter)
	for commandName, command := range on.Commands {
		commandValidations.Validate(&command, emitter.Into(commandName))
	}
}
