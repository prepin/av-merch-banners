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
		CreditUseCase:   NewCreditUseCase(r.TransactionManager, r.TransactionRepo, r.UserRepo, r.UserInfoCache),
		SendCoinUseCase: NewSendCoinUseCase(r.TransactionManager, r.TransactionRepo, r.UserRepo, r.UserInfoCache),
		OrderUseCase:    NewOrderUseCase(r.TransactionManager, r.TransactionRepo, r.UserRepo, r.ItemRepo, r.OrderRepo, r.UserInfoCache),
		InfoUseCase:     NewInfoUseCase(r.UserRepo, r.TransactionRepo, r.OrderRepo, r.UserInfoCache),
	}
}
