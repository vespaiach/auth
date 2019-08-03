package mysqlrepo

import (
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func TestUserRepo(t *testing.T) {

	prepareFixtures("user_test")

}
