package user

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type userRepositoryDB struct {
	db *sqlx.DB
}

func NewUserRepositoryDB(db *sqlx.DB) userRepositoryDB {
	return userRepositoryDB{
		db: db,
	}
}

func (r userRepositoryDB) QueryUser(ctx context.Context, request map[string]interface{}) (*[]User, error) {
	users := make([]User, 0)
	query := `
		SELECT	id,
				first_name,
				last_name,
				phone,
				email,
				balance,
				date_time
		FROM public."user"
		WHERE 1 = 1
	`
	for key := range request {
		query = fmt.Sprintf("%s AND %s = :%s", query, key, key)
	}
	rows, err := r.db.NamedQueryContext(ctx, query, request)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var user User
		if err := rows.StructScan(&user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	defer rows.Close()
	return &users, nil
}
