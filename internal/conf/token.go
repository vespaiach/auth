package conf

import (
	log "github.com/sirupsen/logrus"
	"github.com/vespaiach/auth/internal/comtype"
	"github.com/vespaiach/auth/pkg/gotils"
)

var (
	defaultAccessTokenDuration  = 120   // minutes
	defaultRefreshTokenDuration = 3600  // minutes
	defaultUseRefreshToken      = false // minutes
)

// TokenConfig holds all token's configuration
type TokenConfig struct {
	AccessTokenDuration  int
	RefreshTokenDuration int
	UseRefreshToken      bool
}

func loadTokenConfig() (config *TokenConfig, err error) {
	AccessTokenDuration, err := gotils.GetEnvInt("ACCESS_TOKEN_DURATION")
	if err != nil {
		log.Println(err)
		AccessTokenDuration = defaultAccessTokenDuration
		err = comtype.ErrAppConfigMissingOrWrongSet
	}

	RefreshTokenDuration, err := gotils.GetEnvInt("REFRESH_TOKEN_DURATION")
	if err != nil {
		log.Println(err)
		RefreshTokenDuration = defaultRefreshTokenDuration
		err = comtype.ErrAppConfigMissingOrWrongSet
	}

	UseRefreshToken, err := gotils.GetEnvBool("USE_REFRESH_TOKEN")
	if err != nil {
		log.Println(err)
		UseRefreshToken = defaultUseRefreshToken
		err = comtype.ErrAppConfigMissingOrWrongSet
	}

	config = &TokenConfig{
		AccessTokenDuration,
		RefreshTokenDuration,
		UseRefreshToken,
	}

	return
}
