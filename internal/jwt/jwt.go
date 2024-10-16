package jwt

import (
	"errors"
	"fmt"
	"os"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/goldenfealla/gear-manager/domain"
	"github.com/google/uuid"
)

type TokenType = string
type AccessToken = string
type RefrestToken = string

const (
	ACCESS_TOKEN_SECRET  TokenType = "ACCESS_TOKEN_SECRET"
	REFRESH_TOKEN_SECRET TokenType = "REFRESH_TOKEN_SECRET"
)

func generate(data map[string]interface{}, tt TokenType, d time.Duration) (string, error) {
	claims := make(jwt.MapClaims)

	for k, v := range data {
		claims[k] = v
	}

	// Expiration
	claims["exp"] = jwt.NewNumericDate(time.Now().Add(d))

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	ts, err := t.SignedString([]byte(os.Getenv(tt)))

	if err != nil {
		return "", err
	}

	return ts, nil
}

func generateKeyFunc(tt TokenType) jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv(tt)), nil
	}
}

func parse(ts string, tt TokenType) (jwt.MapClaims, error) {
	token, err := jwt.Parse(ts, generateKeyFunc(tt))

	switch {
	case token.Valid:
		{
			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				return claims, nil
			}
			return nil, fmt.Errorf("unknown error")
		}
	case errors.Is(err, jwt.ErrTokenMalformed):
		return nil, fmt.Errorf("this is not a token")
	case errors.Is(err, jwt.ErrTokenSignatureInvalid):
		return nil, fmt.Errorf("invalid signature")
	case errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet):
		return nil, fmt.Errorf("token expired or not active")
	}

	return nil, fmt.Errorf("can't handle this token")
}

func ValidateRefreshToken(token string) (jwt.MapClaims, error) {
	claims, err := parse(token, REFRESH_TOKEN_SECRET)

	if err != nil {
		return claims, nil
	}

	return claims, nil
}

func GenerateRefreshToken(u *domain.UserInfo) (string, error) {
	// 30 days
	return generate(map[string]interface{}{
		"id":       u.ID,
		"username": u.Username,
		"email":    u.Email,
	}, REFRESH_TOKEN_SECRET, time.Second*2592000)
}

func GenerateAccessToken(rt RefrestToken) (string, error) {
	claims, err := parse(rt, REFRESH_TOKEN_SECRET)

	if err != nil {
		return "", err
	}

	// 5 mins
	return generate(map[string]interface{}{
		"username": claims["username"],
		"email":    claims["email"],
	}, ACCESS_TOKEN_SECRET, time.Second*300)
}

func ParseAccessToken(at AccessToken) (*domain.UserInfo, error) {
	claims, err := parse(at, REFRESH_TOKEN_SECRET)

	if err != nil {
		return nil, err
	}

	return &domain.UserInfo{
		ID:       claims["id"].(uuid.UUID),
		Username: claims["username"].(string),
		Email:    claims["email"].(string),
	}, nil
}
