package db

import (
	"fmt"
	"time"

	"github.com/asymptoter/geochallenge-backend/base/config"
	"github.com/asymptoter/geochallenge-backend/base/ctx"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func NewMySQL() (*sqlx.DB, error) {
	cfg := config.Value.MySQL
	connectionString := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=true&multiStatements=true", cfg.Username, cfg.Password, cfg.Address, cfg.DatabaseName)

	var err error
	var db *sqlx.DB
	connectionCount := 0
	context := ctx.Background()
	fmt.Println(connectionString)
	fmt.Println("retry:", cfg.ConnectionRetry)
	// Connect to MySQL
	for connectionCount < cfg.ConnectionRetry {
		db, err = sqlx.Connect("mysql", connectionString)
		if db != nil && err == nil {
			break
		}
		context.Error("sqlx.Connect failed ", err)
		connectionCount++
		time.Sleep(5 * time.Second)
	}

	return db, err
}
