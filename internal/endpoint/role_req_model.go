package endpoint

import (
	"time"

	"github.com/vespaiach/auth/internal/comtype"
)

// CreateRoleRequest model
type CreateRoleRequest struct {
	RoleName string `json:"role_name"`
	RoleDesc string `json:"role_desc"`
}

// UpdateRoleRequest model
type UpdateRoleRequest struct {
	ID       int64
	RoleName string `json:"role_name"`
	RoleDesc string `json:"role_desc"`
	Active   *bool  `json:"active"`
}

// QueryRoleRequest model
type QueryRoleRequest struct {
	Take     int
	RoleName string
	Active   *bool
	SortBy   string
}

// RoleResponse model
type RoleResponse struct {
	ID        int64     `json:"id"`
	RoleName  string    `json:"role_name"`
	RoleDesc  string    `json:"role_desc"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

// Validate CreateRoleRequest
func (req CreateRoleRequest) Validate() *comtype.CommonError {
	valid := true
	result := make(map[string]string)

	if len(req.RoleName) == 0 {
		valid = false
		result["role_name"] = "role name is required"
	}

	if len(req.RoleName) > 63 {
		valid = false
		result["role_name"] = "role name must be less than 63 characters"
	}

	if len(req.RoleDesc) > 254 {
		valid = false
		result["role_desc"] = "role desc must be less than 254 characters"
	}

	if !valid {
		return comtype.NewCommonError(nil, "CreateRoleRequest - Validate:", comtype.ErrDataValidationFail, result)
	}

	return nil
}

// Validate UpdateRoleRequest
func (req UpdateRoleRequest) Validate() *comtype.CommonError {
	valid := true
	result := make(map[string]string)

	if req.ID == 0 {
		valid = false
		result["id"] = "missing id field"
	}

	if len(req.RoleName) > 63 {
		valid = false
		result["role_name"] = "role name must be less than 63 characters"
	}

	if len(req.RoleDesc) > 254 {
		valid = false
		result["role_desc"] = "role desc must be less than 254 characters"
	}

	if !valid {
		return comtype.NewCommonError(nil, "UpdateRoleRequest - Validate:", comtype.ErrDataValidationFail, result)
	}

	return nil
}
