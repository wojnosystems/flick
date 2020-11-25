package command

type SubCommands struct {
	Name     string
	Usage    string
	Options  interface{}
	Commands []Er
}
