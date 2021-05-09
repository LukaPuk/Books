package driver

import (
	"database/sql"
	"github.com/lib/pq"
	"os"
	"time"
)

var DB *sql.DB

func InitPostgres() error {
	ElephantUrl := os.Getenv("ELEPHANTSQL_URL")
	pgUrl, err := pq.ParseURL(ElephantUrl)
	if err != nil {
		return err
	}
	DB, err = sql.Open("postgres", pgUrl)
	err = DB.Ping()
	if err != nil {
		return err
	}
	DB.SetMaxOpenConns(100)
	DB.SetMaxIdleConns(0)
	DB.SetConnMaxLifetime(time.Millisecond)

	return nil
}
