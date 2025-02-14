package auth

import "golang.org/x/crypto/bcrypt"

type BCryptHashService struct {
}

func NewBCryptHashService() *BCryptHashService {
	return &BCryptHashService{}
}

func (hs *BCryptHashService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func (hs *BCryptHashService) CompareWithPassword(hashed, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password))
	return err == nil
}
