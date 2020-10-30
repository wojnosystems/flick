package action

type WithoutFlagsArgs struct {
	Action func() (err error)
}

func (a *WithoutFlagsArgs) Invoke() (err error) {
	return a.Action()
}

type WithoutFlagsArgErrors struct {
	Action func()
}

func (a *WithoutFlagsArgErrors) Invoke() (err error) {
	a.Action()
	return nil
}
