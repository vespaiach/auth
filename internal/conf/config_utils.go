package conf

import (
	"fmt"
)

// BuildMysqlDSN returns mysqldsn
func (config *AppConfig) BuildMysqlDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s", config.DbUser, config.DbPass, config.DbHost,
		config.DbPort, config.DbName, config.DbOption)
}
