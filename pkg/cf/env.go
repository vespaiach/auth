package cf


import (
	"fmt"
	"log"
	"os"
	"path"
)

var (
	defaultBcryptCost           = 10
	defaultErrorFile            = "./error.yml"
	defaultServerAddress        = ":4000"
	defaultSigningText          = "key_signing"
	defaultDbhost               = "localhost"
	defaultDbport               = "3306"
	defaultDbname               = "auth"
	defaultDbuser               = "root"
	defaultDbpass               = "password"
	defaultDboption             = "charset=utf8&parseTime=True&loc=Local&multiStatements=True&maxAllowedPacket=0"
	defaultAccessTokenDuration  = 120   // minutes
	defaultRefreshTokenDuration = 0  // minutes
)

// AppConfig holds all app's settings and will be read from env
type AppConfig struct {
	AppDir               string
	ErrorFilePath        string
	ServerAddress        string
	BcryptCost           int
	SigningText          string
	DbHost               string
	DbPort               string
	DbName               string
	DbUser               string
	DbPass               string
	DbOption             string
	AccessTokenDuration  int
	RefreshTokenDuration int
}

// BuildMysqlDSN returns mysqldsn
func (config *AppConfig) BuildMysqlDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s", config.DbUser, config.DbPass, config.DbHost,
		config.DbPort, config.DbName, config.DbOption)
}

// LoadAppConfig returns service's configuration
func LoadAppConfig() *AppConfig {
	AppDir, err := getEnvString("APP_DIR")
	if err != nil {
		AppDir, err = os.Getwd()
	}

	ErrorFile, err := getEnvString("ERROR_FILE")
	if err != nil {
		log.Println(err)
		ErrorFile = path.Join(AppDir, defaultErrorFile)
	}

	ServerAddress, err := getEnvString("SERVER_ADDRESS")
	if err != nil {
		log.Println(err)
		ServerAddress = defaultServerAddress
	}

	BcryptCost, err := getEnvInt("BCRYPT_COST")
	if err != nil {
		log.Println(err)
		BcryptCost = defaultBcryptCost
	}

	SigningText, err := getEnvString("SIGNING_TEXT")
	if err != nil {
		log.Println(err)
		SigningText = defaultSigningText
	}

	DbHost, err := getEnvString("DB_HOST")
	if err != nil {
		log.Println(err)
		DbHost = defaultDbhost
	}

	DbPort, err := getEnvString("DB_PORT")
	if err != nil {
		log.Println(err)
		DbPort = defaultDbport
	}

	DbName, err := getEnvString("DB_NAME")
	if err != nil {
		log.Println(err)
		DbName = defaultDbname
	}

	DbUser, err := getEnvString("DB_USER")
	if err != nil {
		log.Println(err)
		DbUser = defaultDbuser
	}

	DbPass, err := getEnvString("DB_PASS")
	if err != nil {
		log.Println(err)
		DbPass = defaultDbpass
	}

	DbOption, err := getEnvString("DB_OPTION")
	if err != nil {
		log.Println(err)
		DbOption = defaultDboption
	}

	AccessTokenDuration, err := getEnvInt("ACCESS_TOKEN_DURATION")
	if err != nil {
		log.Println(err)
		AccessTokenDuration = defaultAccessTokenDuration
	}

	RefreshTokenDuration, err := getEnvInt("REFRESH_TOKEN_DURATION")
	if err != nil {
		log.Println(err)
		RefreshTokenDuration = defaultRefreshTokenDuration
	}

	return &AppConfig{
		AppDir,
		ErrorFile,
		ServerAddress,
		BcryptCost,
		SigningText,
		DbHost,
		DbPort,
		DbName,
		DbUser,
		DbPass,
		DbOption,
		AccessTokenDuration,
		RefreshTokenDuration,
	}
}

