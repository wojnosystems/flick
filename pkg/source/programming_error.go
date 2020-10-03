package source

type ErrProgramming struct {
	msg string
}

func NewErrProgramming(msg string) *ErrProgramming {
	return &ErrProgramming{
		msg: msg,
	}
}

func (p ErrProgramming) Error() string {
	return "programming error: " + p.msg
}
