package usecase

type Usecases struct {
	AuthUseCase      *AuthUseCase
	CreditUseCase    *CreditUseCase
	SendCoinsUseCase *SendCoinsUseCase
}

func NewUsecases(r Repos, s Services) Usecases {
	return Usecases{
		AuthUseCase:      NewAuthUsecase(r.TransactionManager, r.UserRepo, r.TransactionRepo, s.TokenService, s.HashService),
		CreditUseCase:    NewCreditUseCase(r.TransactionManager, r.TransactionRepo, r.UserRepo),
		SendCoinsUseCase: (*SendCoinsUseCase)(NewSendCoinUseCase(r.TransactionManager, r.TransactionRepo, r.UserRepo)),
	}
}
