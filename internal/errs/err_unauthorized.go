package errs

type NoAccessError struct {
	Err error
}

var ErrNoAccessError = NoAccessError{}

func (err NoAccessError) Error() string {
	return err.Err.Error()
}

func (err NoAccessError) Is(target error) bool {
	_, ok := target.(NotFoundError)
	return ok
}
