package dsl

type NamedOptionsOrReferences map[string]OptionOrReference

func (o NamedOptionsOrReferences) HasAny() bool {
	return len(o) != 0
}
