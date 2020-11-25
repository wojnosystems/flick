package context

type Er interface {
	Exit(code int)
	ExitWithMessage(code int, msg string)
}

// ActionContext
// Has arguments as there are no sub-commands. Sub-commands use the positional arguments to determine groupings of commands
type Action interface {
	Er
	PositionalArgs() []string
}

// SubCommandContext
// these have sub-commands, therefore cannot have positional arguments as the next positional argument is going to be a sub-command.
type SubCommand interface {
	Er
}
