package parse

import flag_unmarshaler "github.com/wojnosystems/go-flag-unmarshaler"

// FlagsWithIgnore parses all the flags in the group except those listed in the ignore list
// Any flags in the ignore list will not be listed as errors, however, if a flag is not unmarshalled but isn't on the
// ignore list, it will be marked as an error
func FlagsWithIgnore(globalGroup *flag_unmarshaler.Group, ignore []string) EnvUnmarshaler {
	return &flags{
		globalGroup: globalGroup,
	}
}

// FlagsIgnoreUndefined parses all flags in the group except those which do not exist in the Umarshal target
// Any flags parsed without a destination are ignored and no error will be returned
func FlagsIgnoreUndefined(globalGroup *flag_unmarshaler.Group) EnvUnmarshaler {
	return FlagsWithIgnore(globalGroup, []string{})
}
