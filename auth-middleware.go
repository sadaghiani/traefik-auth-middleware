package traefik_auth_middleware

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
	NameOfAuthorizationHeader string `json:"nameOfAuthorizationHeader"`
	NameOfUserIDClaim         string `json:"nameOfUserIDClaim"`
	NameOfRoleIDClaim         string `json:"nameOfRoleIDClaim"`
	NameForUserIDHeader       string `json:"nameForUserIDHeader"`
	NameForRoleIDHeader       string `json:"nameForRoleIDHeader"`
}

func CreateConfig() *Config {
	return &Config{}
}

type AuthMiddleware struct {
	next                      http.Handler
	name                      string
	SecretKey                 string
	NameOfAuthorizationHeader string
	NameOfUserIDClaim         string
	NameOfRoleIDClaim         string
	NameForUserIDHeader       string
	NameForRoleIDHeader       string
	logger                    *log.Logger
}

func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {

	return &AuthMiddleware{
		next:                      next,
		name:                      name,
		SecretKey:                 config.SecretKey,
		NameOfAuthorizationHeader: config.NameOfAuthorizationHeader,
		NameOfUserIDClaim:         config.NameOfUserIDClaim,
		NameOfRoleIDClaim:         config.NameOfRoleIDClaim,
		NameForUserIDHeader:       config.NameForUserIDHeader,
		NameForRoleIDHeader:       config.NameForRoleIDHeader,
		logger:                    log.New(os.Stdout, "", log.LstdFlags),
	}, nil
}

func (a *AuthMiddleware) ServeHTTP(rw http.ResponseWriter, req *http.Request) {

	if len(req.Header.Get(a.NameOfAuthorizationHeader)) != 0 {
		splitToken := strings.Fields(req.Header.Get(a.NameOfAuthorizationHeader))
		if len(splitToken) == 2 || strings.ToLower(splitToken[0]) == "bearer" {
			res, err := a.getPayloadData(splitToken[1])
			if err == nil {
				req.Header.Set(a.NameForUserIDHeader, res[a.NameForUserIDHeader])
				req.Header.Set(a.NameForRoleIDHeader, res[a.NameForRoleIDHeader])
			} else {
				a.logger.Println(err)
			}
		}
	}

	a.next.ServeHTTP(rw, req)
}

func (a *AuthMiddleware) getPayloadData(tokenInput string) (map[string]string, error) {

	data := make(map[string]string)
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(tokenInput, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(a.SecretKey), nil
	})
	if err != nil {
		return nil, err
	}

	if fmt.Sprint(claims[a.NameOfUserIDClaim]) == "<nil>" || fmt.Sprint(claims[a.NameOfRoleIDClaim]) == "<nil>" {
		return nil, fmt.Errorf("invalid token payload")
	}
	data[a.NameForUserIDHeader] = fmt.Sprint(claims[a.NameOfUserIDClaim])
	data[a.NameForRoleIDHeader] = fmt.Sprint(claims[a.NameOfRoleIDClaim])
	return data, nil
}
