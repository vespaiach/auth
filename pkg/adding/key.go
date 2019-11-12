package adding

import (
	"errors"

	"github.com/vespaiach/auth/pkg/common"
)

// ServiceKey model
type ServiceKey struct {
	Key  string
	Desc string
}

func (sk *ServiceKey) Validate() error {
	payload := make([]string, 0, 2)
	valid := true

	if len(sk.Key) == 0 || len(sk.Key) > 32 {
		valid = false
		payload = append(payload, "key name must be from 1 to 32 characters")
	}

	if len(sk.Desc) > 64 {
		valid = false
		payload = append(payload, "key description must be less than 64 characters")
	}

	if !valid {
		return common.NewAppErr(errors.New("key data is not valid"), common.ErrDataFailValidation)
	}

	return nil
}
