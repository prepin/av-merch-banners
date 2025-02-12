package usecase

type Usecases struct {
	AuthUseCase   *AuthUseCase
	CreditUseCase *CreditUseCase
}

func NewUsecases(r Repos, s Services) Usecases {
	return Usecases{
		AuthUseCase:   NewAuthUsecase(r.UserRepo, s.TokenService, s.HashService),
		CreditUseCase: NewCreditUseCase(r.TransactionRepo, r.UserRepo),
	}
}
