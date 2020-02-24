package db

import (
	"fmt"
	"time"

	"github.com/asymptoter/practice-backend/base/config"
	"github.com/asymptoter/practice-backend/base/ctx"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func MustNew(dbType string) *sqlx.DB {
	res, err := NewDB(dbType)
	if err != nil {
		panic("New" + dbType + " failed by " + err.Error())
	}
	return res
}

func NewDB(dbType string) (*sqlx.DB, error) {
	cfg := config.Value.MySQL
	connectionString := ""
	switch dbType {
	case "mysql":
		connectionString = fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=true&multiStatements=true", cfg.Username, cfg.Password, cfg.Address, cfg.DatabaseName)
	case "postgresql":
		connectionString = fmt.Sprintf("")
	default:
		panic("Must specify dbType")
	}

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
