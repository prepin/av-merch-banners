package errs

type ErrInsufficientFunds struct {
	Err error
}

var ErrInsufficientFundsError = ErrInsufficientFunds{}

func (err ErrInsufficientFunds) Error() string {
	return err.Err.Error()
}

func (err ErrInsufficientFunds) Is(target error) bool {
	_, ok := target.(ErrInsufficientFunds)
	return ok
}

type ErrIncorrectAmount struct {
	Err error
}

var ErrIncorrectAmountError = ErrIncorrectAmount{}

func (err ErrIncorrectAmount) Error() string {
	return err.Err.Error()
}

func (err ErrIncorrectAmount) Is(target error) bool {
	_, ok := target.(ErrIncorrectAmount)
	return ok
}
