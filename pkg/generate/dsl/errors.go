package dsl

type ErrReference struct {
	refValue string
}

func newErrReference(refValue string) *ErrReference {
	return &ErrReference{
		refValue: refValue,
	}
}

func (e *ErrReference) Error() string {
	return "undefined reference: '" + e.refValue + "'"
}
