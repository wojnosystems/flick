package parse

import flag_unmarshaler "github.com/wojnosystems/go-flag-unmarshaler"

func FlagsWithIgnore(globalGroup *flag_unmarshaler.Group, ignore []string) EnvUnmarshaler {
	return &flags{
		globalGroup: globalGroup,
	}
}

func FlagsIgnoreUndefined(globalGroup *flag_unmarshaler.Group) EnvUnmarshaler {
	return FlagsWithIgnore(globalGroup, []string{})
}
