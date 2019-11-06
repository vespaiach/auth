package endpoint

import (
	"time"

	checker "github.com/asaskevich/govalidator"
	"github.com/vespaiach/auth/internal/comtype"
)

// VerifyUserRequest model
type VerifyUserRequest struct {
	Username      string `json:"username"`
	Password      string `json:"password"`
	RemoteAddr    string
	XForwardedFor string
	XRealIP       string
	UserAgent     string
}

// RegisterUserRequest model
type RegisterUserRequest struct {
	FullName string `json:"full_name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// UpdateUserRequest model
type UpdateUserRequest struct {
	ID       int64
	FullName string `json:"full_name"`
}

// QueryUsersRequest model
type QueryUsersRequest struct {
	Take     int
	FullName string
	Username string
	Verified *bool
	Active   *bool
	SortBy   string
}

// ChangeUserPasswordRequest model
type ChangeUserPasswordRequest struct {
	ID          int64
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

// VerifyUserResponse model
type VerifyUserResponse struct {
	Token string `json:"access_token"`
}

// UserResponse model
type UserResponse struct {
	ID        int64     `json:"id,omitempty"`
	FullName  string    `json:"full_name,omitempty"`
	Username  string    `json:"username,omitempty"`
	Email     string    `json:"email,omitempty"`
	Verified  bool      `json:"verified,omitempty"`
	Active    bool      `json:"active,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

// Validate RegisterUserRequest
func (reg RegisterUserRequest) Validate() *comtype.CommonError {
	valid := true
	result := make(map[string]string)

	if len(reg.Email) == 0 || !checker.IsEmail(reg.Email) {
		valid = false
		result["email"] = "email is not valid"
	}

	if len(reg.Username) == 0 || len(reg.Username) > 64 {
		valid = false
		result["username"] = "username is not valid"
	} else if !checker.IsAlphanumeric(reg.Username) {
		valid = false
		result["username"] = "username is not alphanumeric"
	}

	if len(reg.Password) == 0 {
		valid = false
		result["password"] = "password is required"
	}

	if len(reg.FullName) == 0 || len(reg.FullName) > 64 {
		valid = false
		result["full_name"] = "full_name is not valid"
	}

	if !valid {
		return comtype.NewCommonError(nil, "RegisterUserRequest - Validate:", comtype.ErrDataValidationFail, result)
	}

	return nil
}

// Validate UpdateUserRequest
func (reg UpdateUserRequest) Validate() *comtype.CommonError {
	valid := true
	result := make(map[string]string)

	if reg.ID <= 0 {
		valid = false
		result["id"] = "missing id field"
	}

	if len(reg.FullName) == 0 || len(reg.FullName) > 64 {
		valid = false
		result["full_name"] = "full_name is not valid"
	}

	if !valid {
		return comtype.NewCommonError(nil, "UpdateUserRequest - Validate:", comtype.ErrDataValidationFail, result)
	}

	return nil
}

// Validate ChangeUserPasswordRequest
func (reg ChangeUserPasswordRequest) Validate() *comtype.CommonError {
	valid := true
	result := make(map[string]string)

	if reg.ID <= 0 {
		valid = false
		result["id"] = "missing id field"
	}

	if len(reg.OldPassword) == 0 {
		valid = false
		result["old_password"] = "missing old_password"
	}

	if len(reg.NewPassword) == 0 {
		valid = false
		result["new_password"] = "missing new_password"
	}

	if !valid {
		return comtype.NewCommonError(nil, "ChangeUserPasswordRequest - Validate:", comtype.ErrDataValidationFail, result)
	}

	return nil
}

// Validate VerifyUserRequest
func (reg VerifyUserRequest) Validate() *comtype.CommonError {
	valid := true
	result := make(map[string]string)

	if len(reg.Username) == 0 || len(reg.Username) > 64 {
		valid = false
		result["username"] = "invalid username or password"
	} else if !checker.IsAlphanumeric(reg.Username) {
		valid = false
		result["username"] = "invalid username or password"
	}

	if len(reg.Password) == 0 {
		valid = false
		result["username"] = "invalid username or password"
	}

	if !valid {
		return comtype.NewCommonError(nil, "RegisterUserRequest - Validate:", comtype.ErrDataValidationFail, result)
	}

	return nil
}
