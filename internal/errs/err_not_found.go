package errs

type ErrNotFound struct {
	Err error
}

var ErrNotFoundError = ErrNotFound{}

func (err ErrNotFound) Error() string {
	return err.Err.Error()
}

func (err ErrNotFound) Is(target error) bool {
	_, ok := target.(ErrNotFound)
	return ok
}
