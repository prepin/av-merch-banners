package errs

type InternalError struct {
	Err error
}

func (err InternalError) Error() string {
	return err.Err.Error()
}

func (err InternalError) Unwrap() error {
	return err.Err
}

func (err InternalError) Is(target error) bool {
	_, ok := target.(InternalError)
	return ok
}
