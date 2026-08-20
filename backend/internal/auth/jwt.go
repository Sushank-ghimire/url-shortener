package auth

import (
	"errors"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService struct {
	secret []byte
}

func NewJWTService(secret string) *JWTService {
	return &JWTService{
		secret: []byte(secret),
	}
}

func (j *JWTService) GenerateAccessToken(
	userID int64,
	duration time.Duration,
) (string, error) {
	now := time.Now()

	claims := jwt.MapClaims{
		"sub": strconv.FormatInt(userID, 10),
		"iat": now.Unix(),
		"exp": now.Add(duration).Unix(),
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)
	return token.SignedString(j.secret)
}

func (j *JWTService) ValidateAccessToken(
	tokenString string,
) (int64, error) {
	token, err := jwt.Parse(
		tokenString,
		func(token *jwt.Token) (interface{}, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, errors.New("unexpected signing method")
			}
			return j.secret, nil
		},
	)

	if err != nil || !token.Valid {
		return 0, errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("invalid claims")
	}

	sub, ok := claims["sub"].(string)

	if !ok {
		return 0, errors.New("invalid subject")
	}
	userID, err := strconv.ParseInt(sub, 10, 64)
	if err != nil {
		return 0, errors.New("invalid user id")
	}

	return userID, nil
}
