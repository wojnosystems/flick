package parse

import (
	flag_unmarshaler "github.com/wojnosystems/go-flag-unmarshaler"
)

type OptionUnmarshaler interface {
	UnmarshalRoot(into interface{}, rootGroup flag_unmarshaler.Group) (err error)
	UnmarshalCmd(into interface{}, cmdGroup flag_unmarshaler.Group) (err error)
}

type BeforeAndAfter struct {
	BeforeFlags []Unmarshaler
	Flags       FlagUnmarshaler
	AfterFlags  []Unmarshaler
}

func (b *BeforeAndAfter) UnmarshalRoot(into interface{}, rootGroup flag_unmarshaler.Group) (err error) {
	if into != nil {
		// load options
		for _, unmarshaler := range b.BeforeFlags {
			err = unmarshaler.Unmarshal(into)
			if err != nil {
				return
			}
		}
		err = b.Flags.Unmarshal(into, rootGroup)
		if err != nil {
			return
		}
		for _, unmarshaler := range b.AfterFlags {
			err = unmarshaler.Unmarshal(into)
			if err != nil {
				return
			}
		}
	}
	return
}

func (b *BeforeAndAfter) UnmarshalCmd(into interface{}, cmdGroup flag_unmarshaler.Group) (err error) {
	if b.Flags != nil {
		err = b.Flags.Unmarshal(into, cmdGroup)
	}
	return
}
