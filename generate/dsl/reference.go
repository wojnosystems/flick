package dsl

type OptionOrReference struct {
	Reference string `yaml:"$ref"`
	Option
}
