/*Package docker wraps common used docker utilities*/
package docker

import (
	"database/sql"
	"fmt"

	"github.com/asymptoter/practice-backend/base/redis"

	"github.com/jmoiron/sqlx"
	"github.com/ory/dockertest"
	log "github.com/sirupsen/logrus"
)

var Pool *dockertest.Pool
var postgreSQLResource *dockertest.Resource
var redisResource *dockertest.Resource

func GetPostgreSQL() *sqlx.DB {
	var db *sql.DB
	var err error
	if Pool == nil {
		Pool, err = dockertest.NewPool("")
		if err != nil {
			log.Fatalf("Could not connect to docker: %s", err)
		}
	}

	postgreSQLResource, err = Pool.Run("postgres", "9.6", []string{"POSTGRES_USER=a", "POSTGRES_PASSWORD=b", "POSTGRES_DB=practice"})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	if err = Pool.Retry(func() error {
		s := fmt.Sprintf("postgres://a:b@localhost:%s/%s?sslmode=disable", postgreSQLResource.GetPort("5432/tcp"), "practice")
		db, err = sql.Open("postgres", s)
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	return sqlx.NewDb(db, "postgres")
}

func PurgePostgreSQL() {
	if err := Pool.Purge(postgreSQLResource); err != nil {
		log.Fatalf("Purge postgresql failed")
	}
}

func GetRedis() redis.Service {
	var res redis.Service
	var err error

	if Pool == nil {
		Pool, err = dockertest.NewPool("")
		if err != nil {
			log.Fatalf("Could not connect to docker: %s", err)
		}
	}

	redisResource, err = Pool.Run("redis", "3.2", nil)
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	if err = Pool.Retry(func() error {
		addr := fmt.Sprintf("localhost:%s", redisResource.GetPort("6379/tcp"))
		res = redis.NewService(addr)
		return nil
		//return res.redis.Get().Ping().Err()
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}
	return res
}

func PurgeRedis() {
	if err := Pool.Purge(redisResource); err != nil {
		log.Error("Purge redis failed")
	}
}
