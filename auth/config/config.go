package config

import (
	"crypto/rsa"
	"fmt"
	"log"
	"os"

	"github.com/vespaiach/gotils"
)

var (
	defaultPrivateKeyPath       = "/etc/auth/id_rsa"
	defaultPublicKeyPath        = "/etc/auth/id_rsa.pub"
	defaultDbhost               = "localhost"
	defaultDbport               = "3306"
	defaultDbname               = "auth"
	defaultDbuser               = "root"
	defaultDbpass               = "123"
	defaultBcryptCost           = 10
	defaultAccessTokenDuration  = 120   // minutes
	defaultRefreshTokenDuration = 3600  // minutes
	defaultUseRefreshToken      = false // minutes
)

// ServiceConfig holds all service's configuration
type ServiceConfig struct {
	PrivateKeyPath       string
	PublicKeyPath        string
	PrivateKey           *rsa.PrivateKey
	PublicKey            *rsa.PublicKey
	DbHost               string
	DbPort               string
	DbName               string
	DbUser               string
	DbPass               string
	BcryptCost           int
	AccessTokenDuration  int
	RefreshTokenDuration int
	UseRefreshToken      bool
}

// LoadConfig returns service's configuration
func LoadConfig() (config *ServiceConfig, err error) {
	PrivateKeyPath := os.Getenv("PRIVATE_KEY_PATH")
	if len(PrivateKeyPath) == 0 {
		fmt.Println(err)
		PrivateKeyPath = defaultPrivateKeyPath
	}

	PublicKeyPath := os.Getenv("PUBLIC_KEY_PATH")
	if len(PublicKeyPath) == 0 {
		fmt.Println(err)
		PublicKeyPath = defaultPublicKeyPath
	}

	DbHost := os.Getenv("DB_HOST")
	if len(DbHost) == 0 {
		fmt.Println(err)
		DbHost = defaultDbhost
	}

	DbPort := os.Getenv("DB_PORT")
	if len(DbPort) == 0 {
		fmt.Println(err)
		DbPort = defaultDbport
	}

	DbName := os.Getenv("DB_NAME")
	if len(DbName) == 0 {
		fmt.Println(err)
		DbName = defaultDbname
	}

	DbUser := os.Getenv("DB_USER")
	if len(DbUser) == 0 {
		fmt.Println(err)
		DbUser = defaultDbuser
	}

	DbPass := os.Getenv("DB_PASS")
	if len(DbPass) == 0 {
		fmt.Println(err)
		DbPass = defaultDbpass
	}

	BcryptCost, err := gotils.GetEnvInt("BCRYPT_COST")
	if err != nil {
		fmt.Println(err)
		BcryptCost = defaultBcryptCost
	}

	AccessTokenDuration, err := gotils.GetEnvInt("ACCESS_TOKEN_DURATION")
	if err != nil {
		fmt.Println(err)
		AccessTokenDuration = defaultAccessTokenDuration
	}

	RefreshTokenDuration, err := gotils.GetEnvInt("REFRESH_TOKEN_DURATION")
	if err != nil {
		fmt.Println(err)
		RefreshTokenDuration = defaultRefreshTokenDuration
	}

	UseRefreshToken, err := gotils.GetEnvBool("USE_REFRESH_TOKEN")
	if err != nil {
		fmt.Println(err)
		UseRefreshToken = defaultUseRefreshToken
	}

	PrivateKey, err := gotils.LoadRsaPrivateKey(PrivateKeyPath)
	if err != nil {
		log.Fatal(err.Error())
	}

	PublicKey, err := gotils.LoadRsaPublicKey(PublicKeyPath)
	if err != nil {
		log.Fatal(err.Error())
	}

	config = &ServiceConfig{
		PrivateKeyPath,
		PublicKeyPath,
		PrivateKey,
		PublicKey,
		DbHost,
		DbPort,
		DbName,
		DbUser,
		DbPass,
		BcryptCost,
		AccessTokenDuration,
		RefreshTokenDuration,
		UseRefreshToken,
	}

	return
}
