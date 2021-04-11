package parse

import "github.com/wojnosystems/go-optional/v2"

// OptionSourceTrace is a convenience structure that will turn on tracing
// how options were set to help users of a cli/server debug where configuration went wrong
type OptionSourceTrace struct {
	ConfigTrace optional.Bool `env:"FLAG_TRACE" flag:"flagTrace" help:"true to show how all parameters are being set"`
}

func (f OptionSourceTrace) Flags() []string {
	return []string{"--flagTrace"}
}
