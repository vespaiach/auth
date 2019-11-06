package endpoint

import (
	"github.com/vespaiach/auth/internal/comtype"
)

// CreateUserRoleRequest model
type CreateUserRoleRequest struct {
	UserID int64 `json:"user_id"`
	RoleID int64 `json:"role_id"`
}

// DeleteUserRoleRequest model
type DeleteUserRoleRequest struct {
	UserID int64 `json:"user_id"`
	RoleID int64 `json:"role_id"`
}

// QueryUserRoleRequest model
type QueryUserRoleRequest struct {
	Take   int
	UserID int64
	RoleID int64
}

// UserRoleResponse model
type UserRoleResponse struct {
	ID     int64        `json:"id"`
	UserID int64        `json:"user_id"`
	RoleID int64        `json:"role_id"`
	User   UserResponse `json:"user"`
	Role   RoleResponse `json:"role"`
}

// Validate CreateUserRoleRequest
func (req CreateUserRoleRequest) Validate() *comtype.CommonError {
	valid := true
	result := make(map[string]string)

	if req.RoleID <= 0 {
		valid = false
		result["role_id"] = "role_id is required"
	}

	if req.UserID <= 0 {
		valid = false
		result["user_id"] = "user_id is required"
	}

	if !valid {
		return comtype.NewCommonError(nil, "CreateUserRoleRequest - Validate:", comtype.ErrDataValidationFail, result)
	}

	return nil
}

// Validate DeleteUserRoleRequest
func (req DeleteUserRoleRequest) Validate() *comtype.CommonError {
	valid := true
	result := make(map[string]string)

	if req.RoleID <= 0 {
		valid = false
		result["role_id"] = "role_id is required"
	}

	if req.UserID <= 0 {
		valid = false
		result["user_id"] = "user_id is required"
	}

	if !valid {
		return comtype.NewCommonError(nil, "DeleteUserRoleRequest - Validate:", comtype.ErrDataValidationFail, result)
	}

	return nil
}
