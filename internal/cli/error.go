package cli

type DetailedError struct {
	msg string
	err error
}

func (e DetailedError) Error() string {
	return e.err.Error()
}
