package auth_middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type Config struct {
	SecretKey                 string `json:"secretKey"`
	GetAuthorizationHeaderKey string `json:"getAuthorizationHeaderKey"`
	SetUserIDHeaderKey        string `json:"setUserIDHeaderKey"`
}

func CreateConfig() *Config {
	return &Config{}
}

type AuthMiddleware struct {
	next                      http.Handler
	name                      string
	SecretKey                 string
	GetAuthorizationHeaderKey string
	SetUserIDHeaderKey        string
	logger                    *log.Logger
}

func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {

	return &AuthMiddleware{
		next:                      next,
		name:                      name,
		SecretKey:                 config.SecretKey,
		GetAuthorizationHeaderKey: config.GetAuthorizationHeaderKey,
		SetUserIDHeaderKey:        config.SetUserIDHeaderKey,
		logger:                    log.New(os.Stdout, "", log.LstdFlags),
	}, nil
}

func (a *AuthMiddleware) ServeHTTP(rw http.ResponseWriter, req *http.Request) {

	if len(req.Header.Get(a.GetAuthorizationHeaderKey)) != 0 {
		splitToken := strings.Fields(req.Header.Get(a.GetAuthorizationHeaderKey))
		if len(splitToken) == 2 || strings.ToLower(splitToken[0]) == "bearer" {
			res, err := a.getUserID(splitToken[1])
			if err == nil {
				req.Header.Set(a.SetUserIDHeaderKey, res)
			} else {
				a.logger.Println(err)
			}
		}
	}

	a.next.ServeHTTP(rw, req)
}

func (a *AuthMiddleware) getUserID(tokenInput string) (string, error) {

	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(tokenInput, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(a.SecretKey), nil
	})
	if err != nil {
		return "", err
	}

	if fmt.Sprint(claims["userId"]) == "<nil>" {
		return "", fmt.Errorf("invalid token payload")
	}
	return fmt.Sprint(claims["userId"]), nil
}
