package endpoint

import (
	"github.com/vespaiach/auth/internal/comtype"
)

// CreateRoleActionRequest model
type CreateRoleActionRequest struct {
	RoleID   int64 `json:"role_id"`
	ActionID int64 `json:"action_id"`
}

// DeleteRoleActionRequest model
type DeleteRoleActionRequest struct {
	RoleID   int64 `json:"role_id"`
	ActionID int64 `json:"action_id"`
}

// QueryRoleActionRequest model
type QueryRoleActionRequest struct {
	Take     int
	RoleID   int64
	ActionID int64
}

// RoleActionResponse model
type RoleActionResponse struct {
	ID       int64          `json:"id"`
	RoleID   int64          `json:"action_name"`
	ActionID int64          `json:"action_desc"`
	Role     RoleResponse   `json:"role"`
	Action   ActionResponse `json:"action"`
}

// Validate CreateRoleActionRequest
func (req CreateRoleActionRequest) Validate() *comtype.CommonError {
	valid := true
	result := make(map[string]string)

	if req.RoleID <= 0 {
		valid = false
		result["role_id"] = "role_id is required"
	}

	if req.ActionID <= 0 {
		valid = false
		result["action"] = "action_id is required"
	}

	if !valid {
		return comtype.NewCommonError(nil, "CreateRoleActionRequest - Validate:", comtype.ErrDataValidationFail, result)
	}

	return nil
}

// Validate DeleteRoleActionRequest
func (req DeleteRoleActionRequest) Validate() *comtype.CommonError {
	valid := true
	result := make(map[string]string)

	if req.RoleID <= 0 {
		valid = false
		result["role_id"] = "role_id is required"
	}

	if req.ActionID <= 0 {
		valid = false
		result["action"] = "action_id is required"
	}

	if !valid {
		return comtype.NewCommonError(nil, "DeleteRoleActionRequest - Validate:", comtype.ErrDataValidationFail, result)
	}

	return nil
}
