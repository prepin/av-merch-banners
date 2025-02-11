package errs

type ErrNoAccess struct {
	Err error
}

var ErrNoAccessError = ErrNoAccess{}

func (err ErrNoAccess) Error() string {
	return err.Err.Error()
}

func (err ErrNoAccess) Is(target error) bool {
	_, ok := target.(ErrNotFound)
	return ok
}
