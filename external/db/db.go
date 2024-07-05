package db

import (
	"fmt"
	"time"

	"github.com/asymptoter/practice-backend/base/ctx"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func MustNew(dbType, connectionString string) *sqlx.DB {
	res, err := New(dbType, connectionString)
	if err != nil {
		panic("New" + dbType + " failed by " + err.Error())
	}
	return res
}

func New(dbType, connectionString string) (db *sqlx.DB, err error) {
	connectionCount := 0
	context := ctx.Background()
	fmt.Println(connectionString)
	// Connect to MySQL
	for connectionCount < 10 {
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
