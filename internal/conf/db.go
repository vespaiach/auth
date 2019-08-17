package conf

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/vespaiach/auth/internal/comtype"
	"github.com/vespaiach/auth/pkg/gotils"
)

var (
	defaultDbhost   = "localhost"
	defaultDbport   = "3306"
	defaultDbname   = "auth"
	defaultDbuser   = "root"
	defaultDbpass   = "123"
	defaultDboption = "charset=utf8&parseTime=True&loc=Local&multiStatements=True&maxAllowedPacket=0"
)

// DbConfig holds all db's configuration
type DbConfig struct {
	DbHost   string
	DbPort   string
	DbName   string
	DbUser   string
	DbPass   string
	DbOption string
}

// BuildMysqlDSN returns mysqldsn
func (config *DbConfig) BuildMysqlDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s", config.DbUser, config.DbPass, config.DbHost,
		config.DbPort, config.DbName, config.DbOption)
}

// LoadDbConfig returns db's configuration
func loadDbConfig() (config *DbConfig, err error) {
	DbHost, e := gotils.GetEnvString("DB_HOST")
	if e != nil {
		log.Println(e)
		DbHost = defaultDbhost
		err = comtype.ErrAppConfigMissingOrWrongSet
	}

	DbPort, e := gotils.GetEnvString("DB_PORT")
	if e != nil {
		log.Println(e)
		DbPort = defaultDbport
		err = comtype.ErrAppConfigMissingOrWrongSet
	}

	DbName, e := gotils.GetEnvString("DB_NAME")
	if e != nil {
		log.Println(e)
		DbName = defaultDbname
		err = comtype.ErrAppConfigMissingOrWrongSet
	}

	DbUser, e := gotils.GetEnvString("DB_USER")
	if e != nil {
		log.Println(e)
		DbUser = defaultDbuser
		err = comtype.ErrAppConfigMissingOrWrongSet
	}

	DbPass, e := gotils.GetEnvString("DB_PASS")
	if e != nil {
		log.Println(e)
		DbPass = defaultDbpass
		err = comtype.ErrAppConfigMissingOrWrongSet
	}

	DbOption, e := gotils.GetEnvString("DB_OPTION")
	if e != nil {
		log.Println(e)
		DbOption = defaultDboption
		err = comtype.ErrAppConfigMissingOrWrongSet
	}

	config = &DbConfig{
		DbHost,
		DbPort,
		DbName,
		DbUser,
		DbPass,
		DbOption,
	}

	return
}
