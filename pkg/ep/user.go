package ep

import (
	"context"
	"database/sql"
	"github.com/vespaiach/auth/pkg/cf"
	"github.com/vespaiach/auth/pkg/common"
	"github.com/vespaiach/auth/pkg/usrmgr"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type User struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type AddingUser struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ModifyingUser struct {
	Lookup      string
	Username    string `json:"username"`
	Email       string `json:"email"`
	NewPassword string `json:"new_password"`
	OldPassword string `json:"old_password"`
	Active      *bool  `json:"active"`
}

type QueryingUser struct {
	Username string
	Email    string
	Active   sql.NullBool
	Sort     string
	Page     int64
	PerPage  int64
}

type Users struct {
	Records []*User `json:"records"`
	Total   int64   `json:"total"`
	Page    int64   `json:"page"`
	PerPage int64   `json:"per_page"`
}

type AddingBunchesToUser struct {
	Bunches  []string `json:"bunches"`
	Username string
}

func AddingUserEndPoint(ctx context.Context, request interface{}) (interface{}, error) {
	erch := make(chan error)
	uch := make(chan *usrmgr.User)
	userv := ctx.Value(common.UserManagementService).(usrmgr.Service)
	appConfig := ctx.Value(common.AppConfigContextKey).(*cf.AppConfig)

	go func() {
		req, ok := request.(*AddingUser)
		if !ok {
			erch <- common.ErrWrongInputDatatype
			return
		}

		if len(req.Password) == 0 {
			erch <- common.ErrPasswordMissing
			return
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), appConfig.BcryptCost)
		if err != nil {
			erch <- err
			return
		}

		id, err := userv.AddUser(req.Username, req.Email, string(hash))
		if err != nil {
			erch <- err
			return
		}

		u, err := userv.GetUser(id)
		if err != nil {
			erch <- err
			return
		}
		uch <- u
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case e := <-erch:
		return nil, e
	case u := <-uch:
		return &User{
			u.ID,
			u.Username,
			u.Email,
			u.Active.Bool,
			u.CreatedAt,
			u.UpdatedAt,
		}, nil
	}
}

func GettingUserEndPoint(ctx context.Context, request interface{}) (interface{}, error) {
	erch := make(chan error)
	uch := make(chan *usrmgr.User)
	userv := ctx.Value(common.UserManagementService).(usrmgr.Service)

	go func() {
		name, ok := request.(string)
		if !ok {
			erch <- common.ErrWrongInputDatatype
			return
		}

		user, err := userv.GetUserByUsername(name)
		if err != nil {
			erch <- err
			return
		}
		if user == nil {
			erch <- common.ErrUserNotFound
			return
		}

		uch <- user
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case e := <-erch:
		return nil, e
	case u := <-uch:
		return &User{
			u.ID,
			u.Username,
			u.Email,
			u.Active.Bool,
			u.CreatedAt,
			u.UpdatedAt,
		}, nil
	}
}

func ModifyingUserEndPoint(ctx context.Context, request interface{}) (interface{}, error) {
	erch := make(chan error)
	success := make(chan bool)
	userv := ctx.Value(common.UserManagementService).(usrmgr.Service)

	go func() {
		req, ok := request.(*ModifyingUser)
		if !ok {
			erch <- common.ErrWrongInputDatatype
			return
		}

		user, err := userv.GetUserByUsername(req.Lookup)
		if err != nil {
			erch <- err
			return
		}
		if user == nil {
			erch <- common.ErrUserNotFound
			return
		}

		active := sql.NullBool{}
		if req.Active != nil {
			active.Bool = *req.Active
			active.Valid = true
		}

		err = userv.ModifyUser(user.ID, req.Username, req.Email, "", active)
		if err != nil {
			erch <- err
			return
		}
		success <- true
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case e := <-erch:
		return nil, e
	case <-success:
		return true, nil
	}
}

func QueryingUserEndPoint(ctx context.Context, request interface{}) (interface{}, error) {
	erch := make(chan error)
	uch := make(chan []*usrmgr.User)
	userv := ctx.Value(common.UserManagementService).(usrmgr.Service)
	params, ok := request.(*QueryingUser)

	var total int64

	go func() {
		if !ok {
			erch <- common.ErrWrongInputDatatype
			return
		}

		if params.PerPage == 0 {
			params.PerPage = common.Take
		}

		if params.Page == 0 {
			params.Page = 1
		}

		records, count, err := userv.QueryUsers(params.Page, params.PerPage, params.Username, params.Email,
			params.Active, params.Sort)
		if err != nil {
			erch <- err
			return
		}
		total = count
		uch <- records
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case e := <-erch:
		return nil, e
	case lst := <-uch:
		rows := make([]*User, 0, len(lst))
		for _, row := range lst {
			rows = append(rows, &User{
				row.ID,
				row.Username,
				row.Email,
				row.Active.Bool,
				row.CreatedAt,
				row.UpdatedAt,
			})
		}
		return &Users{
			rows,
			total,
			params.Page,
			params.PerPage,
		}, nil
	}
}

func AddingBunchesToUserEndPoint(ctx context.Context, request interface{}) (interface{}, error) {
	erch := make(chan error)
	qch := make(chan bool)
	userv := ctx.Value(common.UserManagementService).(usrmgr.Service)

	go func() {
		req, ok := request.(*AddingBunchesToUser)
		if !ok {
			erch <- common.ErrWrongInputDatatype
			return
		}

		err := userv.AddBunchesToUser(req.Username, req.Bunches)
		if err != nil {
			erch <- err
			return
		}

		qch <- true
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case e := <-erch:
		return nil, e
	case <-qch:
		return true, nil
	}
}

func GettingBunchesOfUserEndPoint(ctx context.Context, request interface{}) (interface{}, error) {
	erch := make(chan error)
	bch := make(chan []*usrmgr.Bunch)
	userv := ctx.Value(common.UserManagementService).(usrmgr.Service)

	go func() {
		name, ok := request.(string)
		if !ok {
			erch <- common.ErrWrongInputDatatype
			return
		}

		bunches, err := userv.GetBunches(name)
		if err != nil {
			erch <- err
			return
		}

		bch <- bunches
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case e := <-erch:
		return nil, e
	case lst := <-bch:
		rows := make([]*Bunch, 0, len(lst))
		for _, row := range lst {
			rows = append(rows, &Bunch{
				row.ID,
				row.Name,
				row.Desc,
				row.Active.Bool,
				row.CreatedAt,
				row.UpdatedAt,
			})
		}
		return rows, nil
	}
}

func GettingKeysOfUserEndPoint(ctx context.Context, request interface{}) (interface{}, error) {
	erch := make(chan error)
	kch := make(chan []*usrmgr.Key)
	userv := ctx.Value(common.UserManagementService).(usrmgr.Service)

	go func() {
		name, ok := request.(string)
		if !ok {
			erch <- common.ErrWrongInputDatatype
			return
		}

		keys, err := userv.GetKeys(name)
		if err != nil {
			erch <- err
			return
		}

		kch <- keys
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case e := <-erch:
		return nil, e
	case lst := <-kch:
		rows := make([]*Key, 0, len(lst))
		for _, row := range lst {
			rows = append(rows, &Key{
				row.ID,
				row.Key,
				row.Desc,
				row.CreatedAt,
				row.UpdatedAt,
			})
		}
		return rows, nil
	}
}
