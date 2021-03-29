package dsl

type NamedOptions map[string]Option

func (o NamedOptions) HasAny() bool {
	return len(o) != 0
}
