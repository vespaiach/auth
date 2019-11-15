package listing

import "time"

type Bunch struct {
	ID        int64
	Name      string
	Desc      string
	Active    bool
	UpdatedAt time.Time
	CreatedAt time.Time
}
