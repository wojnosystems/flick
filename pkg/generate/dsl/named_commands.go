package dsl

type NamedCommands map[string]Command

func (c NamedCommands) HasAny() bool {
	return len(c) != 0
}
