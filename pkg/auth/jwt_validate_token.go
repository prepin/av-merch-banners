package auth

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

func (js *JWTService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString, &Claims{},
		func(token *jwt.Token,
		) (interface{}, error) {
			return js.secretKey, nil
		})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
