package action

type Invoker interface {
	Invoke() (err error)
}
