package adding

import (
	"errors"

	"github.com/vespaiach/auth/pkg/common"
)

// Bunch model
type Bunch struct {
	Name string
	Desc string
}

func (b *Bunch) Validate() error {
	payload := make([]string, 0, 7)
	valid := true

	if len(b.Name) == 0 || len(b.Name) > 32 {
		valid = false
		payload = append(payload, "bunch name must be from 1 to 32 characters")
	}

	if len(b.Desc) > 64 {
		valid = false
		payload = append(payload, "bunch description must be less than 64 characters")
	}

	if !valid {
		return common.NewAppErr(errors.New("bunch data is not valid"), common.ErrDataFailValidation)
	}

	return nil
}
