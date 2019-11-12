package adding

import (
	"errors"
	"regexp"

	"github.com/vespaiach/auth/pkg/common"
)

// User model
type User struct {
	Username string
	Email    string
	Hash     string
}

// Validate user data before adding
func (u *User) Validate() error {
	payload := make([]string, 0, 7)
	valid := true

	if len(u.Username) == 0 {
		valid = false
		payload = append(payload, "username is missing")
	}

	if len(u.Username) > 32 {
		valid = false
		payload = append(payload, "username exceeds 32 characters")
	}

	if matched, err := regexp.Match(`^[a-z0-9_]{%d,%d}$`, []byte(u.Username)); !matched || err != nil {
		valid = false
		payload = append(payload, "username contains special characters or white space characters")
	}

	if len(u.Email) == 0 {
		valid = false
		payload = append(payload, "email address is missing")
	}

	if len(u.Email) > 32 {
		valid = false
		payload = append(payload, "email address exceeds 64 characters")
	}

	if matched, err := regexp.Match(`"^[a-z0-9_@\\-\\.]{1,127}$"`, []byte(u.Email)); !matched || err != nil {
		valid = false
		payload = append(payload, "email address is invalid")
	}

	if len(u.Hash) == 0 {
		valid = false
		payload = append(payload, "password hash is missing")
	}

	if !valid {
		return common.NewAppErr(errors.New("user data is not valid"), common.ErrDataFailValidation)
	}

	return nil
}
