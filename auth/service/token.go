package service

import (
	"time"

	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/vespaiach/auth/config"
	"github.com/vespaiach/auth/store"
)

var (
	signingAlgorithm = "RS256"
	issuer           = "Vespaiach"
)

type accessTokenMapClaims struct {
	jwtgo.StandardClaims
	actions []string
}

func issueAccessToken(u *store.User, id string, config *config.ServiceConfig) (token string, err error) {
	duration := time.Duration(config.AccessTokenDuration) * time.Minute
	actions := make([]string, len(u.Actions))

	for _, a := range u.Actions {
		actions = append(actions, a.Name)
	}

	for _, r := range u.Roles {
		tmp := make([]string, len(r.Actions))
		for _, a := range r.Actions {
			tmp = append(tmp, a.Name)
		}
		actions = append(actions, tmp...)
	}

	accessTokenObj := jwtgo.New(jwtgo.GetSigningMethod(signingAlgorithm))

	accessTokenObj.Claims = accessTokenMapClaims{
		jwtgo.StandardClaims{
			Id:        id,
			Issuer:    issuer,
			Audience:  u.Username,
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(duration).Unix(),
		},
		actions,
	}

	return accessTokenObj.SignedString(config.PrivateKey)
}

func issueRefreshToken(u *store.User, id string, config *config.ServiceConfig) (token string, err error) {
	tokenObj := jwtgo.New(jwtgo.GetSigningMethod(signingAlgorithm))

	tokenObj.Claims = jwtgo.StandardClaims{
		Id:        id,
		Issuer:    issuer,
		Audience:  u.Username,
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(time.Duration(config.RefreshTokenDuration) * time.Minute).Unix(),
	}

	return tokenObj.SignedString(config.PrivateKey)
}
