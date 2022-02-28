package database

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

func NewPostgresConn() (*sqlx.DB, error) {
	connString := fmt.Sprintf("user=%s dbname=%s password=%s host=%s port=%s connect_timeout=%d sslmode=%s",
		viper.GetString("postgres.username"),
		viper.GetString("postgres.database"),
		viper.GetString("postgres.password"),
		viper.GetString("postgres.host"),
		viper.GetString("postgres.port"),
		viper.GetInt("postgres.timeout"),
		viper.GetString("postgres.sslmode"),
	)
	conn, err := sqlx.Open(viper.GetString("postgres.type"), connString)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
