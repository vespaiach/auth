package conf

import (
	log "github.com/sirupsen/logrus"
	"github.com/vespaiach/auth/pkg/gotils"
	"gopkg.in/go-playground/validator.v9"
)

var (
	defaultBcryptCost = 10
	defaultAppDir     = "/home/vespaiach/working/go/src/github.com/vespaiach/auth"
)

// CommonConfig holds all common's configuration
type CommonConfig struct {
	BcryptCost int
	AppDir     string
}

// AppConfig holad all app's configurations
type AppConfig struct {
	DbConfig     *DbConfig
	RsaKeyConfig *RsaKeyConfig
	TokenConfig  *TokenConfig
	CommonConfig *CommonConfig
	Validator    *validator.Validate
}

// LoadAppConfig returns service's configuration
func LoadAppConfig() *AppConfig {
	BcryptCost, err := gotils.GetEnvInt("BCRYPT_COST")
	if err != nil {
		log.Println(err)
		BcryptCost = defaultBcryptCost
	}

	AppDir, err := gotils.GetEnvString("APP_DIR")
	if err != nil {
		log.Println(err)
		AppDir = defaultAppDir
	}

	dbConfig, _ := loadDbConfig()
	rsaKeyConfig, _ := loadRsaConfig()
	tokenConfig, _ := loadTokenConfig()
	commonConfig := &CommonConfig{
		BcryptCost,
		AppDir,
	}

	v := validator.New()

	config := &AppConfig{
		dbConfig,
		rsaKeyConfig,
		tokenConfig,
		commonConfig,
		v,
	}

	return config
}
