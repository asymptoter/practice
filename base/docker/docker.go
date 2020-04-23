/*Package docker wraps common used docker utilities*/
package docker

import (
	"database/sql"
	"fmt"

	"github.com/asymptoter/practice-backend/base/redis"

	"github.com/jmoiron/sqlx"
	"github.com/ory/dockertest"
	dc "github.com/ory/dockertest/docker"
	log "github.com/sirupsen/logrus"
)

var Pool *dockertest.Pool
var postgreSQLResource *dockertest.Resource
var redisResource *dockertest.Resource

func initPool() {
	var err error
	if Pool == nil {
		Pool, err = dockertest.NewPool("")
		if err != nil {
			log.Fatalf("Could not connect to docker: %s", err)
		}
	}
}

func getContainerID(image string) string {
	cs, err := Pool.Client.ListContainers(dc.ListContainersOptions{})
	if err != nil {
		log.Fatalf("Could not list containers: %s", err)
	}

	name := "/unittest_" + image
	for _, c := range cs {
		if c.Names[0] == name {
			return c.ID
		}
	}
	return ""
}

func GetPostgreSQL() *sqlx.DB {
	var db *sql.DB
	var err error

	initPool()

	cID := getContainerID("postgres")
	if len(cID) != 0 {
		postgreSQLResource = &dockertest.Resource{}
		container, err := Pool.Client.InspectContainer(cID)
		if err != nil {
			log.Fatalf("Could not get containers: %s", err)
		}
		postgreSQLResource.Container = container
	} else {
		postgreSQLResource, err = Pool.RunWithOptions(&dockertest.RunOptions{Name: "unittest_postgres", Repository: "postgres", Tag: "latest", Env: []string{"POSTGRES_USER=a", "POSTGRES_PASSWORD=b", "POSTGRES_DB=practice"}})
		if err != nil {
			log.Fatalf("Could not start resource: %s", err)
		}
	}

	if err = Pool.Retry(func() error {
		s := fmt.Sprintf("postgres://a:b@localhost:%s/%s?sslmode=disable", postgreSQLResource.GetPort("5432/tcp"), "practice")
		db, err = sql.Open("postgres", s)
		if err != nil {
			return err
		}
		return nil
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
	initPool()

	cID := getContainerID("redis")

	if len(cID) != 0 {
		redisResource = &dockertest.Resource{}
		container, err := Pool.Client.InspectContainer(cID)
		if err != nil {
			log.Fatalf("Could not get containers: %s", err)
		}
		redisResource.Container = container
	} else {
		redisResource, err = Pool.RunWithOptions(&dockertest.RunOptions{Name: "unittest_redis", Repository: "redis", Tag: "latest"})
		if err != nil {
			log.Fatalf("Could not start resource: %s", err)
		}
	}

	if err = Pool.Retry(func() error {
		addr := fmt.Sprintf("localhost:%s", redisResource.GetPort("6379/tcp"))
		res = redis.NewService(addr)
		return nil
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
