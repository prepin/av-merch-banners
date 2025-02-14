package usecase

type Usecases struct {
	AuthUseCase     AuthUseCase
	CreditUseCase   CreditUseCase
	SendCoinUseCase SendCoinUseCase
	OrderUseCase    OrderUseCase
	InfoUseCase     InfoUseCase
}

func NewUsecases(r Repos, s Services) Usecases {
	return Usecases{
		AuthUseCase:     NewAuthUsecase(r.TransactionManager, r.UserRepo, r.TransactionRepo, s.TokenService, s.HashService),
		CreditUseCase:   NewCreditUseCase(r.TransactionManager, r.TransactionRepo, r.UserRepo),
		SendCoinUseCase: NewSendCoinUseCase(r.TransactionManager, r.TransactionRepo, r.UserRepo),
		OrderUseCase:    NewOrderUseCase(r.TransactionManager, r.TransactionRepo, r.UserRepo, r.ItemRepo, r.OrderRepo),
		InfoUseCase:     NewInfoUseCase(r.UserRepo, r.TransactionRepo, r.OrderRepo),
	}
}
