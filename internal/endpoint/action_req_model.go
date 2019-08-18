package endpoint

import (
	"time"

	"github.com/vespaiach/auth/internal/comtype"
)

// CreateActionRequest model
type CreateActionRequest struct {
	ActionName string `json:"action_name"`
	ActionDesc string `json:"action_desc"`
}

// UpdateActionRequest model
type UpdateActionRequest struct {
	ID         int64
	ActionName string `json:"action_name"`
	ActionDesc string `json:"action_desc"`
	Active     *bool  `json:"active"`
}

// QueryActionRequest model
type QueryActionRequest struct {
	Take       int
	ActionName string
	Active     *bool
	SortBy     string
}

// ActionResponse model
type ActionResponse struct {
	ID         int64     `json:"id"`
	ActionName string    `json:"action_name"`
	ActionDesc string    `json:"action_desc"`
	Active     bool      `json:"active"`
	CreatedAt  time.Time `json:"created_at,omitempty"`
	UpdatedAt  time.Time `json:"updated_at,omitempty"`
}

// Validate CreateActionRequest
func (req CreateActionRequest) Validate() *comtype.CommonError {
	valid := true
	result := make(map[string]string)

	if len(req.ActionName) == 0 {
		valid = false
		result["action_name"] = "action name is required"
	}

	if len(req.ActionName) > 63 {
		valid = false
		result["action_name"] = "action name must be less than 63 characters"
	}

	if len(req.ActionDesc) > 254 {
		valid = false
		result["action_desc"] = "action desc must be less than 254 characters"
	}

	if !valid {
		return comtype.NewCommonError(nil, "CreateActionRequest - Validate:", comtype.ErrDataValidationFail, result)
	}

	return nil
}

// Validate UpdateActionRequest
func (req UpdateActionRequest) Validate() *comtype.CommonError {
	valid := true
	result := make(map[string]string)

	if req.ID == 0 {
		valid = false
		result["id"] = "missing id field"
	}

	if len(req.ActionName) > 63 {
		valid = false
		result["action_name"] = "action name must be less than 63 characters"
	}

	if len(req.ActionDesc) > 254 {
		valid = false
		result["action_desc"] = "action desc must be less than 254 characters"
	}

	if !valid {
		return comtype.NewCommonError(nil, "UpdateActionRequest - Validate:", comtype.ErrDataValidationFail, result)
	}

	return nil
}
