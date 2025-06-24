package jwttoken

import (
	"SimpleMessaging/pkg/env"
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type ClaimToken struct {
	Username string `json:"username"`
	Fullname string `json:"full_name"`
	jwt.RegisteredClaims
}

var mapTokenType = map[string]time.Duration{
	"token":         time.Hour * 3,
	"refresh_token": time.Hour * 72,
}

var jwtSecret = []byte(env.GetEnv("APP_SECRET", ""))

func GenerateToken(ctx context.Context, username, fullname, tokenType string) (string, error) {

	claimToken := ClaimToken{
		Username: username,
		Fullname: fullname,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    env.GetEnv("APP_NAME", ""),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(mapTokenType[tokenType])),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claimToken)
	result, err := token.SignedString(jwtSecret)
	if err != nil {
		return result, fmt.Errorf("failed to generate token: %v", err)
	}

	return result, err
}

func ValidateToken(ctx context.Context, tokenStr string) (*ClaimToken, error) {
	var (
		claimToken *ClaimToken
		ok         bool
	)

	jwtToken, err := jwt.ParseWithClaims(tokenStr, &ClaimToken{}, func(t *jwt.Token) (interface{}, error) {
		if t.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse jwt: %v", err)
	}

	if claimToken, ok = jwtToken.Claims.(*ClaimToken); !ok || !jwtToken.Valid {
		return nil, fmt.Errorf("token invalid")
	}

	return claimToken, nil
}
