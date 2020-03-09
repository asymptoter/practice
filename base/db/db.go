package db

import (
	"fmt"
	"time"

	"github.com/asymptoter/practice-backend/base/config"
	"github.com/asymptoter/practice-backend/base/ctx"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func MustNew(dbType string, isContainer bool) *sqlx.DB {
	res, err := NewDB(dbType, isContainer)
	if err != nil {
		panic("New" + dbType + " failed by " + err.Error())
	}
	return res
}

func NewDB(dbType string, isContainer bool) (*sqlx.DB, error) {
	cfg := config.Value.Database
	connectionString := ""
	address := cfg.Address
	if isContainer {
		address = "docker.for.mac." + address
	}
	switch dbType {
	case "mysql":
		connectionString = fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=true&multiStatements=true", cfg.UserName, cfg.Password, address, cfg.DatabaseName)
	case "postgres":
		connectionString = fmt.Sprintf("user=%s password=%s host=%s port=%d dbname=%s sslmode=disable", cfg.UserName, cfg.Password, address, cfg.Port, cfg.DatabaseName)
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
		db, err = sqlx.Connect(dbType, connectionString)
		if db != nil && err == nil {
			break
		}
		context.Error("sqlx.Connect failed ", err)
		connectionCount++
		time.Sleep(5 * time.Second)
	}

	return db, err
}
