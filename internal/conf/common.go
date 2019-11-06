package conf

import (
	"crypto/rsa"
	"os"
	"path"

	log "github.com/sirupsen/logrus"
	"github.com/vespaiach/auth/pkg/gotils"
)

var (
	defaultBcryptCost    = 10
	defaultAppDir        = "/home/vespaiach/working/go/src/github.com/vespaiach/auth"
	defaultErrorFile     = "./error.yml"
	defaultServerAddress = ":4000"
	defaultSigningText   = "key_signing"

	defaultDbhost   = "localhost"
	defaultDbport   = "3306"
	defaultDbname   = "auth"
	defaultDbuser   = "root"
	defaultDbpass   = "password"
	defaultDboption = "charset=utf8&parseTime=True&loc=Local&multiStatements=True&maxAllowedPacket=0"

	defaultPrivateKeyPath = "/configs/rsa/id_rsa"
	defaultPublicKeyPath  = "/configs/rsa/id_rsa_pub"

	defaultAccessTokenDuration  = 120   // minutes
	defaultRefreshTokenDuration = 3600  // minutes
	defaultUseRefreshToken      = false // minutes
)

// AppConfig hold all app's settings and will be read from env
type AppConfig struct {
	AppDir        string
	BcryptCost    int
	ErrorFile     string
	ServerAddress string
	SigningText   string

	DbHost   string
	DbPort   string
	DbName   string
	DbUser   string
	DbPass   string
	DbOption string

	PrivateKeyPath string
	PublicKeyPath  string
	PrivateKey     *rsa.PrivateKey
	PublicKey      *rsa.PublicKey

	AccessTokenDuration  int
	RefreshTokenDuration int
	UseRefreshToken      bool
}

// LoadAppConfig returns service's configuration
func LoadAppConfig() *AppConfig {
	AppDir, err := gotils.GetEnvString("APP_DIR")
	if err != nil {
		AppDir, err = os.Getwd()
	}

	BcryptCost, err := gotils.GetEnvInt("BCRYPT_COST")
	if err != nil {
		log.Println(err)
		BcryptCost = defaultBcryptCost
	}

	SigningText, err := gotils.GetEnvString("SIGNING_TEXT")
	if err != nil {
		log.Println(err)
		SigningText = defaultSigningText
	}

	ErrorFile, err := gotils.GetEnvString("ERROR_FILE")
	if err != nil {
		log.Println(err)
		ErrorFile = path.Join(AppDir, defaultErrorFile)
	}

	ServerAddress, err := gotils.GetEnvString("SERVER_ADDRESS")
	if err != nil {
		log.Println(err)
		ServerAddress = defaultServerAddress
	}

	DbHost, err := gotils.GetEnvString("DB_HOST")
	if err != nil {
		log.Println(err)
		DbHost = defaultDbhost
	}

	DbPort, err := gotils.GetEnvString("DB_PORT")
	if err != nil {
		log.Println(err)
		DbPort = defaultDbport
	}

	DbName, err := gotils.GetEnvString("DB_NAME")
	if err != nil {
		log.Println(err)
		DbName = defaultDbname
	}

	DbUser, err := gotils.GetEnvString("DB_USER")
	if err != nil {
		log.Println(err)
		DbUser = defaultDbuser
	}

	DbPass, err := gotils.GetEnvString("DB_PASS")
	if err != nil {
		log.Println(err)
		DbPass = defaultDbpass
	}

	DbOption, err := gotils.GetEnvString("DB_OPTION")
	if err != nil {
		log.Println(err)
		DbOption = defaultDboption
	}

	PrivateKeyPath, err := gotils.GetEnvString("PRIVATE_KEY_PATH")
	if err != nil {
		log.Println(err)
		PrivateKeyPath = path.Join(AppDir, defaultPrivateKeyPath)
	}

	PublicKeyPath, err := gotils.GetEnvString("PUBLIC_KEY_PATH")
	if err != nil {
		log.Println(err)
		PublicKeyPath = path.Join(AppDir, defaultPublicKeyPath)
	}

	PrivateKey, err := gotils.LoadRsaPrivateKey(PrivateKeyPath)
	if err != nil {
		log.Fatal(err.Error())
	}

	PublicKey, err := gotils.LoadRsaPublicKey(PublicKeyPath)
	if err != nil {
		log.Fatal(err.Error())
	}

	AccessTokenDuration, err := gotils.GetEnvInt("ACCESS_TOKEN_DURATION")
	if err != nil {
		log.Println(err)
		AccessTokenDuration = defaultAccessTokenDuration
	}

	RefreshTokenDuration, err := gotils.GetEnvInt("REFRESH_TOKEN_DURATION")
	if err != nil {
		log.Println(err)
		RefreshTokenDuration = defaultRefreshTokenDuration
	}

	UseRefreshToken, err := gotils.GetEnvBool("USE_REFRESH_TOKEN")
	if err != nil {
		log.Println(err)
		UseRefreshToken = defaultUseRefreshToken
	}

	return &AppConfig{
		AppDir,
		BcryptCost,
		ErrorFile,
		ServerAddress,
		SigningText,

		DbHost,
		DbPort,
		DbName,
		DbUser,
		DbPass,
		DbOption,

		PrivateKeyPath,
		PublicKeyPath,
		PrivateKey,
		PublicKey,

		AccessTokenDuration,
		RefreshTokenDuration,
		UseRefreshToken,
	}
}
