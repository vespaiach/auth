package mysqlrepo

import (
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/vespaiach/auth/internal/comtype"
	"github.com/vespaiach/auth/internal/model"
)

// MysqlTokenHistoryRepo will implement model.TokenHistoryRepo
type MysqlTokenHistoryRepo struct {
	DbClient *sqlx.DB
}

// NewMysqlTokenHistoryRepo create new instance of MysqlTokenHistoryRepo
func NewMysqlTokenHistoryRepo(db *sqlx.DB) model.TokenHistoryRepo {
	return &MysqlTokenHistoryRepo{
		db,
	}
}

var sqlGetTokenHistoryByID = `
SELECT *
FROM token_histories
WHERE token_histories.user_id = ?
;`

// GetByUserID find a token-history by its user's ID
func (r *MysqlTokenHistoryRepo) GetByUserID(userID int64) ([]*model.TokenHistory, *comtype.CommonError) {
	rows, err := r.DbClient.Queryx(sqlGetTokenHistoryByID, userID)
	if err != nil {
		return nil, comtype.NewCommonError(err, "MysqlTokenHistoryRepo - GetByUserID:", comtype.ErrQueryDataFail, nil)
	}
	defer rows.Close()

	results := make([]*model.TokenHistory, 0)
	for rows.Next() {
		var ac model.TokenHistory
		rows.StructScan(&ac)
		if err != nil {
			return nil, comtype.NewCommonError(err, "MysqlTokenHistoryRepo - GetByUserID:", comtype.ErrQueryDataFail, nil)
		}
		results = append(results, &ac)
	}

	err = rows.Err()
	if err != nil {
		return nil, comtype.NewCommonError(err, "MysqlTokenHistoryRepo - GetByUserID:", comtype.ErrQueryDataFail, nil)
	}

	return results, nil
}

var sqlCreateTokenHistory = `
INSERT INTO token_histories(uid, user_id, access_token, refresh_token, remote_addr, x_forwarded_for, x_real_ip, user_agent, 
	created_at, expired_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
`

// Save a token-history
func (r *MysqlTokenHistoryRepo) Save(uid string, userID int64, accessToken string, refreshToken string, remoteAddr string,
	xForwardedFor string, xRealIP string, userAgent string, createdAt time.Time, expiredAt time.Time) *comtype.CommonError {
	stmt, err := r.DbClient.Prepare(sqlCreateTokenHistory)
	if err != nil {
		return comtype.NewCommonError(err, "MysqlTokenHistoryRepo - Create:", comtype.ErrHandleDataFail, nil)
	}

	res, err := stmt.Exec(uid, userID, accessToken, refreshToken, remoteAddr, xForwardedFor, xRealIP, userAgent, createdAt, expiredAt)
	if err != nil {
		return comtype.NewCommonError(err, "MysqlTokenHistoryRepo - Create:", comtype.ErrHandleDataFail, nil)
	}

	rowAffected, err := res.RowsAffected()
	if err != nil || rowAffected == 0 {
		return comtype.NewCommonError(err, "MysqlTokenHistoryRepo - Create:", comtype.ErrHandleDataFail, nil)
	}

	return nil
}
