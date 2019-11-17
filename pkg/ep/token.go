package ep

import (
	"context"
	"github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/endpoint"
	"github.com/google/uuid"
	"github.com/vespaiach/auth/pkg/cf"
	"github.com/vespaiach/auth/pkg/common"
	"github.com/vespaiach/auth/pkg/usrmgr"
	"golang.org/x/crypto/bcrypt"
	"sort"
	"sync"
	"time"

	jwtgo "github.com/dgrijalva/jwt-go"
)

// TokenClaims token's claims
type TokenClaims struct {
	jwtgo.StandardClaims
	Bunches []string
	Keys    []string
}

type Token struct {
	AccessToken string `json:"access_token"`
}

type VerifyingUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// TokenParser will read and parse jwt token from context
func TokenParserMiddleware(ep endpoint.Endpoint) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		appConfig := ctx.Value(common.AppConfigContextKey).(*cf.AppConfig)

		tokenStr, ok := ctx.Value(jwt.JWTTokenContextKey).(string)
		if !ok {
			return nil, common.ErrMissingJWTToken
		}

		token, err := jwtgo.ParseWithClaims(tokenStr, TokenClaims{}, func(token *jwtgo.Token) (interface{}, error) {
			if token.Method != jwtgo.SigningMethodHS256 {
				return nil, common.ErrWrongJWTToken
			}

			return []byte(appConfig.SigningText), nil
		})

		if err != nil {
			return nil, err
		}

		if !token.Valid {
			return nil, common.ErrWrongJWTToken
		}

		ctx = context.WithValue(ctx, jwt.JWTClaimsContextKey, token.Claims)

		return ep(ctx, request)
	}
}

func VerifyingUserMiddleware(ep endpoint.Endpoint) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		erch := make(chan error)
		uch := make(chan *usrmgr.User)
		userv := ctx.Value(common.UserManagementService).(usrmgr.Service)

		go func() {
			req, ok := request.(*VerifyingUser)
			if !ok {
				erch <- common.ErrWrongInputDatatype
				return
			}

			user, err := userv.GetUserByUsername(req.Username)
			if err != nil {
				erch <- err
				return
			}
			if user == nil {
				erch <- common.ErrUserNotFound
				return
			}

			if bcrypt.CompareHashAndPassword([]byte(user.Hash), []byte(req.Password)) != nil {
				erch <- common.ErrWrongCredentials
				return
			}

			uch <- user
		}()

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case e := <-erch:
			return nil, e
		case u := <-uch:
			return ep(ctx, u)
		}
	}
}

// KeyChecker is a middleware for checking key
func KeyChecker(key string) endpoint.Middleware {
	return func(ep endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			claims, ok := ctx.Value(jwt.JWTClaimsContextKey).(TokenClaims)
			if !ok {
				return nil, common.ErrMissingJWTToken
			}

			if sort.SearchStrings(claims.Keys, key) < len(claims.Keys) {
				return ep(ctx, request)
			}
			return nil, common.ErrWrongJWTToken
		}
	}
}

func IssueTokenEndPoint(ctx context.Context, request interface{}) (interface{}, error) {
	erch := make(chan error)
	ach := make(chan string)
	userv := ctx.Value(common.UserManagementService).(usrmgr.Service)
	appConfig := ctx.Value(common.AppConfigContextKey).(*cf.AppConfig)

	duration, err := time.ParseDuration(appConfig.AccessTokenDuration)
	if err != nil {
		return nil, err
	}

	user, ok := request.(*usrmgr.User)
	if !ok {
		return nil, common.ErrWrongInputDatatype
	}

	go func() {
		var (
			wg          sync.WaitGroup
			bunches     []*usrmgr.Bunch
			keys        []*usrmgr.Key
			errGetBunch error
			errGetKey   error
		)

		wg.Add(1)
		go func() {
			defer wg.Done()
			bunches, errGetBunch = userv.GetBunches(user.Username)
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			keys, errGetKey = userv.GetKeys(user.Username)
		}()

		wg.Wait()

		if errGetBunch != nil {
			erch <- errGetBunch
			return
		}

		if errGetKey != nil {
			erch <- errGetKey
			return
		}

		tokenObj := createToken(user, bunches, keys, duration)
		accessToken, err := tokenObj.SignedString([]byte(appConfig.SigningText))
		if err != nil {
			erch <- err
			return
		}
		ach <- accessToken
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case e := <-erch:
		return nil, e
	case accessToken := <-ach:
		return &Token{AccessToken: accessToken}, nil
	}
}

func createToken(user *usrmgr.User, bunches []*usrmgr.Bunch, keys []*usrmgr.Key, duration time.Duration) *jwtgo.Token {
	klst := make([]string, 0, len(keys))
	blst := make([]string, 0, len(bunches))

	for _, k := range keys {
		klst = append(klst, k.Key)
	}

	for _, b := range blst {
		blst = append(blst, b)
	}

	sort.Strings(klst)
	sort.Strings(blst)

	uid := uuid.New().String()
	createdAt := time.Now()
	expiredAt := createdAt.Add(duration)

	return jwtgo.NewWithClaims(jwtgo.SigningMethodHS256, TokenClaims{
		jwtgo.StandardClaims{
			Id:        uid,
			Issuer:    "",
			Audience:  user.Username,
			IssuedAt:  createdAt.Unix(),
			ExpiresAt: expiredAt.Unix(),
		},
		blst,
		klst,
	})
}
