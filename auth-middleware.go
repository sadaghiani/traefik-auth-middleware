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
	GetUserIDHeaderKey        string `json:"getUserIDHeaderKey"`
	GetRoleIDHeaderKey        string `json:"getRoleIDHeaderKey"`
	SetUserIDHeaderKey        string `json:"setUserIDHeaderKey"`
	SetRoleIDHeaderKey        string `json:"setRoleIDHeaderKey"`
}

func CreateConfig() *Config {
	return &Config{}
}

type AuthMiddleware struct {
	next                      http.Handler
	name                      string
	SecretKey                 string
	GetAuthorizationHeaderKey string
	GetUserIDHeaderKey        string
	GetRoleIDHeaderKey        string
	SetUserIDHeaderKey        string
	SetRoleIDHeaderKey        string
	logger                    *log.Logger
}

func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {

	return &AuthMiddleware{
		next:                      next,
		name:                      name,
		SecretKey:                 config.SecretKey,
		GetAuthorizationHeaderKey: config.GetAuthorizationHeaderKey,
		GetUserIDHeaderKey:        config.GetUserIDHeaderKey,
		GetRoleIDHeaderKey:        config.GetRoleIDHeaderKey,
		SetUserIDHeaderKey:        config.SetUserIDHeaderKey,
		SetRoleIDHeaderKey:        config.SetRoleIDHeaderKey,
		logger:                    log.New(os.Stdout, "", log.LstdFlags),
	}, nil
}

func (a *AuthMiddleware) ServeHTTP(rw http.ResponseWriter, req *http.Request) {

	if len(req.Header.Get(a.GetAuthorizationHeaderKey)) != 0 {
		splitToken := strings.Fields(req.Header.Get(a.GetAuthorizationHeaderKey))
		if len(splitToken) == 2 || strings.ToLower(splitToken[0]) == "bearer" {
			res, err := a.getUserID(splitToken[1])
			if err == nil {
				req.Header.Set(a.SetUserIDHeaderKey, res[a.SetUserIDHeaderKey])
				req.Header.Set(a.SetRoleIDHeaderKey, res[a.SetRoleIDHeaderKey])
			} else {
				a.logger.Println(err)
			}
		}
	}

	a.next.ServeHTTP(rw, req)
}

func (a *AuthMiddleware) getUserID(tokenInput string) (map[string]string, error) {

	data := make(map[string]string)
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(tokenInput, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(a.SecretKey), nil
	})
	if err != nil {
		return nil, err
	}

	if fmt.Sprint(claims[a.GetUserIDHeaderKey]) == "<nil>" || fmt.Sprint(claims[a.GetRoleIDHeaderKey]) == "<nil>" {
		return nil, fmt.Errorf("invalid token payload")
	}
	data[a.SetUserIDHeaderKey] = fmt.Sprint(claims[a.GetUserIDHeaderKey])
	data[a.SetRoleIDHeaderKey] = fmt.Sprint(claims[a.GetRoleIDHeaderKey])
	return data, nil
}
