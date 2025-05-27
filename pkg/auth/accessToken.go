package auth

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

/**
 * @Description
 * @Date 2025/4/23 18:23
 **/
var (
	accessTokenKey  = []byte("your_access_secret_key")
	refreshTokenKey = []byte("your_refresh_secret_key")
)

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// GenerateAccessToken 生成 Access Token
func GenerateAccessToken(username string) (string, error) {
	expirationTime := time.Now().Add(15 * time.Minute)
	claims := &Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(accessTokenKey)
}

// ValidateAccessToken 验证 Access Token
func ValidateAccessToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return accessTokenKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("无效的 Access Token")
	}

	return claims, nil
}
