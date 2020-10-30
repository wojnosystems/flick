package parse

import flag_unmarshaler "github.com/wojnosystems/go-flag-unmarshaler"

type usedFlags struct {
	usedFlags map[string]bool
}

// ReceiveSet is called when a value is set on the object.
// Tracks which flagNames are set when flags are unmarshalled into the structure
func (f *usedFlags) ReceiveSet(structPath string, flagName string, value string) {
	if f.usedFlags == nil {
		f.usedFlags = make(map[string]bool)
	}
	f.usedFlags[flagName] = true
}

// unused returns all the flags that were provided, but were not consumed
// these are the extra flags. If you want to create errors on having too many flags, this is where you do it.
// allFlags: a set of flags to subtract from the usedFlags
// out = allFlags - usedFlags
func (f usedFlags) unused(allFlags []flag_unmarshaler.KeyValue) (out []flag_unmarshaler.KeyValue) {
	for _, flag := range allFlags {
		if _, ok := f.usedFlags[flag.Key]; !ok {
			out = append(out, flag)
		}
	}
	return
}
