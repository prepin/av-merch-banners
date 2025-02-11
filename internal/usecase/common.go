package usecase

type Usecases struct {
	AuthUseCase *AuthUseCase
}

func NewUsecases(r Repos, s Services) Usecases {
	return Usecases{
		AuthUseCase: NewAuthUsecase(r.UserRepo, s.TokenService),
	}
}
