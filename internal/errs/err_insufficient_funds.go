package errs

type InsufficientFundsError struct {
	Err error
}

var ErrInsufficientFundsError = InsufficientFundsError{}

func (err InsufficientFundsError) Error() string {
	return err.Err.Error()
}

func (err InsufficientFundsError) Is(target error) bool {
	_, ok := target.(InsufficientFundsError)
	return ok
}

type IncorrectAmountError struct {
	Err error
}

var ErrIncorrectAmountError = IncorrectAmountError{}

func (err IncorrectAmountError) Error() string {
	return err.Err.Error()
}

func (err IncorrectAmountError) Is(target error) bool {
	_, ok := target.(IncorrectAmountError)
	return ok
}
