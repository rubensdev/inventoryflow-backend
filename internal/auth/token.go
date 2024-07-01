package auth

import (
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

type UserClaims struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Firstname string `json:"first_name"`
	Lastname  string `json:"last_name"`
	jwt.RegisteredClaims
}

func NewAccessToken(claims UserClaims, secret string) (string, error) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedStr, err := accessToken.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("error getting access token signed string. %w", err)
	}

	return signedStr, nil
}

func NewRefreshToken(claims jwt.RegisteredClaims, secret string) (string, error) {
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedStr, err := refreshToken.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("error getting refresh token signed string. %w", err)
	}
	return signedStr, nil
}

func ParseAccessToken(accessToken string, secret string) (*UserClaims, error) {
	parsedToken, err := jwt.ParseWithClaims(accessToken, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})
	err = validateParse(parsedToken, err)
	if err != nil {
		return nil, err
	}
	return parsedToken.Claims.(*UserClaims), nil
}

func ParseRefreshToken(refreshToken string, secret string) (*jwt.RegisteredClaims, error) {
	parsedToken, err := jwt.ParseWithClaims(refreshToken, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})
	err = validateParse(parsedToken, err)
	if err != nil {
		return nil, err
	}
	return parsedToken.Claims.(*jwt.RegisteredClaims), nil
}

func validateParse(token *jwt.Token, err error) error {
	if token == nil {
		return fmt.Errorf("token invalid")
	}

	switch {
	case token.Valid:
		return nil
	case errors.Is(err, jwt.ErrTokenMalformed):
		return fmt.Errorf("token malformed")
	case errors.Is(err, jwt.ErrTokenSignatureInvalid):
		return fmt.Errorf("invalid signature")
	case errors.Is(err, jwt.ErrTokenExpired):
		return fmt.Errorf("token is expired")
	case errors.Is(err, jwt.ErrTokenNotValidYet):
		return fmt.Errorf("token is not valid yet")
	default:
		return fmt.Errorf("could not handle this token")
	}
}
