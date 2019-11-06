package service

import (
	"sort"
	"time"

	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/vespaiach/auth/internal/comtype"
	"github.com/vespaiach/auth/internal/conf"
	"github.com/vespaiach/auth/internal/model"
)

var (
	signingAlgorithm = "RS256"
	issuer           = "Vespaiach"
)

// TokenService interface
type TokenService interface {
	IssueToken(user *model.User, actions []*model.Action, roles []*model.Role, remoteAddr string, xForwardedFor string,
		xRealIP string, userAgent string) (string, *comtype.CommonError)
}

// NewTokenService creates a struct that implement ITokenService
func NewTokenService(appRepo *model.AppRepo, appConfig *conf.AppConfig) TokenService {
	return &tokenService{
		appConfig,
		appRepo,
	}
}

type tokenService struct {
	appConfig *conf.AppConfig
	appRepo   *model.AppRepo
}

// TokenClaims token's claims
type TokenClaims struct {
	jwtgo.StandardClaims
	Actions []string
	Roles   []string
}

func (ser *tokenService) IssueToken(user *model.User, actions []*model.Action, roles []*model.Role, remoteAddr string,
	xForwardedFor string, xRealIP string, userAgent string) (string, *comtype.CommonError) {
	duration := time.Duration(ser.appConfig.AccessTokenDuration) * time.Minute
	actionStrings := make([]string, 0, len(actions))
	roleStrings := make([]string, 0, len(roles))

	for _, a := range actions {
		actionStrings = append(actionStrings, a.ActionName)
	}

	for _, r := range roles {
		roleStrings = append(roleStrings, r.RoleName)

		tmp := make([]string, 0, len(r.Actions))
		for _, a := range r.Actions {
			tmp = append(tmp, a.ActionName)
		}
		actionStrings = append(actionStrings, tmp...)
	}

	sort.Strings(roleStrings)
	sort.Strings(actionStrings)

	uid := uuid.New().String()
	createdAt := time.Now()
	expiredAt := createdAt.Add(duration)

	accessTokenObj := jwtgo.NewWithClaims(jwtgo.SigningMethodHS256, TokenClaims{
		jwtgo.StandardClaims{
			Id:        uid,
			Issuer:    issuer,
			Audience:  user.Username,
			IssuedAt:  createdAt.Unix(),
			ExpiresAt: expiredAt.Unix(),
		},
		actionStrings,
		roleStrings,
	})

	token, err := accessTokenObj.SignedString([]byte(ser.appConfig.SigningText))
	if err != nil {
		return "", comtype.NewCommonError(err, "TokenService - IssueToken:", comtype.ErrHandleDataFail, nil)
	}

	saveerr := ser.appRepo.TokenHistoryRepo.Save(uid, user.ID, token, "", remoteAddr, xForwardedFor, xRealIP, userAgent,
		createdAt, expiredAt)
	if err != nil {
		return "", saveerr
	}

	return token, nil
}
