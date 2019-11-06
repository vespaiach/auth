package migrate

import (
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/jmoiron/sqlx"
	"github.com/vespaiach/auth/internal/model"
	"golang.org/x/crypto/bcrypt"
)

type uniqInt struct {
	order int
	mux   sync.Mutex
}

var inc *uniqInt = &uniqInt{order: 1}

func (unq *uniqInt) New() int {
	unq.mux.Lock()
	defer unq.mux.Unlock()
	unq.order++
	return unq.order
}

// Script migration script
type Script struct {
	Name string
	Text string
}

// Migrator struct
type Migrator struct {
	db   *sqlx.DB
	up   []*Script
	down []*Script
}

// NewMigrator return struct instance
func NewMigrator(db *sqlx.DB) *Migrator {

	var upScripts = []*Script{
		&Script{Name: "init_database", Text: initDatabase},
	}

	var downScripts = []*Script{
		&Script{Name: "drop_init_database", Text: dropInitDatabase},
	}

	return &Migrator{
		db:   db,
		up:   upScripts,
		down: downScripts,
	}
}

// Up to run migration scripts
func (m *Migrator) Up() {
	tx := m.db.MustBegin()

	for _, s := range m.up {
		tx.MustExec(santizeSQL(s.Text))
	}

	tx.Commit()
}

// Down to run migration scripts
func (m *Migrator) Down() {
	tx := m.db.MustBegin()

	for i := len(m.down) - 1; i >= 0; i-- {
		tx.MustExec(santizeSQL(m.down[i].Text))
	}

	tx.Commit()
}

// SeedTestData to migrate test data
func (m *Migrator) SeedTestData() {
	tx := m.db.MustBegin()
	tx.MustExec(santizeSQL(initTestData))
	tx.Commit()
}

// SeedProdData to init production data
func (m *Migrator) SeedProdData(adminFullName string, adminUsername string, adminEmail string, adminPassword string, bcryptCost int) {
	tx := m.db.MustBegin()
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(adminPassword), bcryptCost)
	tx.MustExec(santizeSQL(fmt.Sprintf(initProdData, adminFullName, adminUsername, adminEmail, hashedPassword)))
	tx.Commit()
}

// CreateSeedingAction is to create a new action record for testing
func (m *Migrator) CreateSeedingAction(beforeCreate func(map[string]interface{})) *model.Action {
	fields := map[string]interface{}{
		"action_name": fmt.Sprintf("action_name_%s", strconv.Itoa(inc.New())),
		"action_desc": fmt.Sprintf("action_desc_%s", strconv.Itoa(inc.New())),
		"active":      true,
	}

	if beforeCreate != nil {
		beforeCreate(fields)
	}

	result, _ := m.db.NamedExec(sqlInsertAction, fields)
	id, _ := result.LastInsertId()

	rows, _ := m.db.Queryx(sqlGetAction, id)
	defer rows.Close()

	action := new(model.Action)
	rows.Next()
	rows.StructScan(action)

	return action
}

// CreateSeedingRole is to create a new role record for testing
func (m *Migrator) CreateSeedingRole(beforeCreate func(map[string]interface{})) *model.Role {
	fields := map[string]interface{}{
		"role_name": fmt.Sprintf("role_name_%s", strconv.Itoa(inc.New())),
		"role_desc": fmt.Sprintf("role_desc_%s", strconv.Itoa(inc.New())),
		"active":    true,
	}

	if beforeCreate != nil {
		beforeCreate(fields)
	}

	result, _ := m.db.NamedExec(sqlInsertRole, fields)
	id, _ := result.LastInsertId()

	rows, _ := m.db.Queryx(sqlGetRole, id)
	defer rows.Close()

	role := new(model.Role)
	rows.Next()
	rows.StructScan(role)

	return role
}

// CreateSeedingRoleAction is to create a new role-action record for testing
func (m *Migrator) CreateSeedingRoleAction(beforeCreate func(map[string]interface{})) *model.RoleAction {
	fields := map[string]interface{}{}

	if beforeCreate != nil {
		beforeCreate(fields)
	} else {
		action := m.CreateSeedingAction(nil)
		role := m.CreateSeedingRole(nil)

		fields["role_id"] = role.ID
		fields["action_id"] = action.ID
	}

	result, _ := m.db.NamedExec(sqlInsertRoleAction, fields)
	id, _ := result.LastInsertId()

	rows, _ := m.db.Queryx(sqlGetRoleAction, id)
	defer rows.Close()

	roleAction := new(model.RoleAction)
	rows.Next()
	rows.StructScan(roleAction)

	return roleAction
}

// CreateSeedingUserRole is to create a new user-role record for testing
func (m *Migrator) CreateSeedingUserRole(beforeCreate func(map[string]interface{})) *model.UserRole {
	fields := map[string]interface{}{}

	if beforeCreate != nil {
		beforeCreate(fields)
	} else {
		user := m.CreateSeedingUser(nil)
		role := m.CreateSeedingRole(nil)

		fields["role_id"] = role.ID
		fields["user_id"] = user.ID
	}

	result, _ := m.db.NamedExec(sqlInsertUserRole, fields)
	id, _ := result.LastInsertId()

	rows, _ := m.db.Queryx(sqlGetUserRole, id)
	defer rows.Close()

	userRole := new(model.UserRole)
	rows.Next()
	rows.StructScan(userRole)

	return userRole
}

// CreateSeedingUser is to create a new user record for testing
func (m *Migrator) CreateSeedingUser(beforeCreate func(map[string]interface{})) *model.User {
	fields := map[string]interface{}{
		"full_name": fmt.Sprintf("full_name%s", strconv.Itoa(inc.New())),
		"username":  fmt.Sprintf("username%s", strconv.Itoa(inc.New())),
		"email":     fmt.Sprintf("email%s@localhost.com", strconv.Itoa(inc.New())),
		"hashed":    "$2a$10$MRwncTIqhLkQTW74tKrwJORwPGGyC5pfpDru9Srf4ORFQUBxhvWLi", // = password
		"verified":  true,
		"active":    true,
	}

	if beforeCreate != nil {
		beforeCreate(fields)
	}

	result, _ := m.db.NamedExec(sqlInsertUser, fields)
	id, _ := result.LastInsertId()

	rows, _ := m.db.Queryx(sqlGetUser, id)
	defer rows.Close()

	user := new(model.User)
	rows.Next()
	rows.StructScan(user)

	return user
}

// CreateUniqueString is to create unique string for testing
func (m *Migrator) CreateUniqueString(prefix string) string {
	return fmt.Sprintf("%s%s", prefix, strconv.Itoa(inc.New()))
}

func santizeSQL(sql string) string {
	return strings.Replace(sql, `"`, "`", -1)
}
