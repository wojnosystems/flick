package invoke

type WithoutOptions struct {
	Action func() (err error)
}

func (a *WithoutOptions) Invoke() (err error) {
	return a.Action()
}

type WithoutOptionsOrErrors struct {
	Action func()
}

func (a *WithoutOptionsOrErrors) Invoke() (err error) {
	a.Action()
	return nil
}
