package keymgr

import (
	"time"
)

type Key struct {
	ID        int64
	Key       string
	Desc      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
