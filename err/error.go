package err

type InvalidInputError struct {
	Msg string
}

func (e *InvalidInputError) Error() string {
	return e.Msg
}

type ErrInvalidStep struct {
	Msg string
}

func (e *ErrInvalidStep) Error() string {
	return e.Msg
}
