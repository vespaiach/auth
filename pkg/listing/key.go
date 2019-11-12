package listing

import "time"

type Key struct {
	ID        int64
	Key       string
	Desc      string
	UpdatedAt time.Time
	CreatedAt time.Time
}
