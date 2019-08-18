package transport

import (
	"context"
	"errors"
	"sort"

	stdjwt "github.com/dgrijalva/jwt-go"
	"github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/endpoint"

	"github.com/vespaiach/auth/internal/comtype"
	"github.com/vespaiach/auth/internal/conf"
	"github.com/vespaiach/auth/internal/service"
)

func jwtParser(next endpoint.Endpoint) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {

		appConfig := ctx.Value(comtype.CommonKeyAppConfiguration).(*conf.AppConfig)

		tokenString, ok := ctx.Value(jwt.JWTTokenContextKey).(string)
		if !ok {
			return nil, comtype.NewCommonError(nil, "jwtParser - middleware", comtype.ErrMissingCredential, nil)
		}

		token, err := stdjwt.ParseWithClaims(tokenString, service.TokenClaims{}, func(token *stdjwt.Token) (interface{}, error) {
			if token.Method != stdjwt.SigningMethodHS256 {
				return nil, comtype.NewCommonError(nil, "jwtParser - middleware", comtype.ErrInvalidCredential, nil)
			}

			return []byte(appConfig.SigningText), nil
		})

		if err != nil {
			if e, ok := err.(*stdjwt.ValidationError); ok {
				switch {
				case e.Errors&stdjwt.ValidationErrorMalformed != 0:
					return nil, comtype.NewCommonError(errors.New("Token is malformed"), "jwtParser - middleware",
						comtype.ErrInvalidCredential, nil)
				case e.Errors&stdjwt.ValidationErrorExpired != 0:
					return nil, comtype.NewCommonError(errors.New("Token is expired"), "jwtParser - middleware",
						comtype.ErrInvalidCredential, nil)
				case e.Errors&stdjwt.ValidationErrorNotValidYet != 0:
					return nil, comtype.NewCommonError(errors.New("Token is not active yet"), "jwtParser - middleware",
						comtype.ErrInvalidCredential, nil)
				case e.Inner != nil:
					return nil, comtype.NewCommonError(e.Inner, "jwtParser - middleware",
						comtype.ErrInvalidCredential, nil)
				}
			}
			return nil, comtype.NewCommonError(err, "jwtParser - middleware", comtype.ErrInvalidCredential, nil)
		}

		if !token.Valid {
			return nil, comtype.NewCommonError(nil, "jwtParser - middleware", comtype.ErrInvalidCredential, nil)
		}

		ctx = context.WithValue(ctx, jwt.JWTClaimsContextKey, token.Claims)

		return next(ctx, request)
	}
}

func checkAuth(action string) endpoint.Middleware {
	check := func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			claims, ok := ctx.Value(jwt.JWTClaimsContextKey).(service.TokenClaims)
			if !ok {
				return nil, comtype.NewCommonError(errors.New("Missing credential"), "checkAuth - middleware",
					comtype.ErrMissingCredential, nil)
			}

			if sort.SearchStrings(claims.Actions, action) < len(claims.Actions) {
				return next(ctx, request)
			}
			return nil, comtype.NewCommonError(nil, "checkAuth", comtype.ErrNoPermision, nil)
		}
	}

	return check
}
