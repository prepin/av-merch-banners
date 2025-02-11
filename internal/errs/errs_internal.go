package errs

type ErrInternal struct {
	Err error
}

func (err ErrInternal) Error() string {
	return err.Err.Error()
}

func (err ErrInternal) Unwrap() error {
	return err.Err
}

func (err ErrInternal) Is(target error) bool {
	_, ok := target.(ErrInternal)
	return ok
}
