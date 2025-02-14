package errs

type NotFoundError struct {
	Err error
}

var ErrNotFoundError = NotFoundError{}

func (err NotFoundError) Error() string {
	return err.Err.Error()
}

func (err NotFoundError) Is(target error) bool {
	_, ok := target.(NotFoundError)
	return ok
}
