package dsl

type FlagDef struct {
	Name    string   `yaml:"name"`
	Aliases []string `yaml:"aliases"`
}
