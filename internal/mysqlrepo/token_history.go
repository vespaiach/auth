package mysqlrepo

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
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
WHERE token_histories.id =?
LIMIT 1;
`

// GetByID find a token-history by its ID
func (r *MysqlTokenHistoryRepo) GetByID(id int64) (*model.TokenHistory, error) {
	rows, err := r.DbClient.Queryx(sqlGetTokenHistoryByID, id)
	if err != nil {
		log.Error("MysqlTokenHistoryRepo - GetByID:", err)
		return nil, comtype.ErrQueryDataFailed
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, comtype.ErrDataNotFound
	}

	tokenHistory := new(model.TokenHistory)
	err = rows.StructScan(tokenHistory)
	if err != nil {
		log.Error("MysqlTokenHistoryRepo - GetByID:", err)
		return nil, comtype.ErrQueryDataFailed
	}

	return tokenHistory, nil
}

var sqlCreateTokenHistory = `
INSERT INTO token_histories(uid, user_id, access_token, refresh_token, created_at) VALUES(?, ?, ?, ?, ?);
`

// Create a new token-history
func (r *MysqlTokenHistoryRepo) Create(uid string, userID int64, accessToken string, refreshToken string, createdAt time.Time) error {
	stmt, err := r.DbClient.Prepare(sqlCreateTokenHistory)
	if err != nil {
		log.Error("MysqlTokenHistoryRepo - Create:", err)
		return comtype.ErrCreateDataFailed
	}

	res, err := stmt.Exec(uid, userID, accessToken, refreshToken, createdAt)
	if err != nil {
		log.Error("MysqlTokenHistoryRepo - Create:", err)
		return comtype.ErrCreateDataFailed
	}

	rowAffected, err := res.RowsAffected()
	if err != nil || rowAffected == 0 {
		log.Error("MysqlTokenHistoryRepo - Create:", err)
		return comtype.ErrCreateDataFailed
	}

	return nil
}

const sqlListTokenHistory = `
SELECT *
FROM token_histories
%s
ORDER BY %s 
LIMIT :offset, :limit;`

const sqlCountListTokenHistory = `
SELECT Count(*)
FROM token_histories
%s ;`

// Query a list of token_histories
func (r *MysqlTokenHistoryRepo) Query(page int, perPage int, filters map[string]interface{}, sorts map[string]comtype.SortDirection) ([]*model.TokenHistory, int64, error) {
	conditions := sqlWhereBuilder(" AND ", filters)
	sortings := sqlSortingBuilder(sorts)
	filters = sqlLikeConditionFilter(filters)
	filters["offset"] = (page - 1) * perPage
	filters["limit"] = perPage

	ch := make(chan int64)
	go func() {
		var totals int64
		rows, err := r.DbClient.NamedQuery(fmt.Sprintf(sqlCountListTokenHistory, conditions), filters)
		if err != nil {
			log.Error("MysqlTokenHistoryRepo - Query:", err)
			ch <- int64(-1)
			return
		}
		defer rows.Close()

		if !rows.Next() {
			ch <- int64(-1)
			return
		}

		rows.Scan(&totals)
		ch <- totals
		close(ch)
	}()

	rows, err := r.DbClient.NamedQuery(fmt.Sprintf(sqlListTokenHistory, conditions, sortings), filters)
	if err != nil {
		log.Error("MysqlTokenHistoryRepo - Query:", err)
		return nil, 0, comtype.ErrQueryDataFailed
	}
	defer rows.Close()

	results := make([]*model.TokenHistory, 0, perPage)
	for rows.Next() {
		var ac model.TokenHistory
		rows.StructScan(&ac)
		if err != nil {
			log.Error("MysqlTokenHistoryRepo - Query:", err)
			return nil, 0, comtype.ErrQueryDataFailed
		}
		results = append(results, &ac)
	}

	total := <-ch
	if total == -1 {
		return nil, 0, comtype.ErrQueryDataFailed
	}

	return results, total, nil
}
