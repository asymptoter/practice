package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"text/template"
	"time"

	"github.com/asymptoter/geochallenge-backend/apis/auth"
	"github.com/asymptoter/geochallenge-backend/base/config"
	"github.com/asymptoter/geochallenge-backend/base/ctx"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func home(c *gin.Context) {
	t, err := template.ParseFiles("./index.html")
	if err != nil {
		log.Println(err)
	}

	if err := t.Execute(c.Writer, nil); err != nil {
		log.Println(err)
	}
}

func setupMySQL() (*gorm.DB, error) {
	cfg := config.Value.MySQL
	connectionString := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=true&multiStatements=true", cfg.Username, cfg.Password, cfg.Address, cfg.DatabaseName)

	var err error
	var db *gorm.DB
	connectionCount := 0
	context := ctx.Background()
	context.Info(connectionString)
	// Connect to mysql
	for connectionCount < cfg.ConnectionRetry {
		db, err = gorm.Open("mysql", connectionString)
		if db != nil && err == nil {
			break
		}
		context.Error("gorm.Open failed ", err)
		connectionCount++
		time.Sleep(5 * time.Second)
	}

	return db, err
}

func setupRedis() (*redis.Client, error) {
	cfg := config.Value.Redis
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Address,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err := client.Ping().Result()
	if err != nil {
		return nil, err
	}
	return client, nil
}

func newHttpServer(db *gorm.DB, redisClient *redis.Client) *http.Server {
	cfg := config.Value.Server
	r := gin.Default()
	auth.SetHttpHandler(r, db, redisClient)
	r.GET("/home", home)

	return &http.Server{
		Addr:    cfg.Address,
		Handler: r,
	}
}

func main() {
	flag.Parse()

	db, err := setupMySQL()
	if err != nil {
		log.Println("setup MySQL failed ", err)
		return
	}
	defer db.Close()

	redisClient, err := setupRedis()
	if err != nil {
		log.Println("setup Redis failed ", err)
		return
	}
	defer redisClient.Close()

	httpServer := newHttpServer(db, redisClient)
	// Start http server
	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			log.Println("ListenAndServe failed ", err)
		}
	}()

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGTERM)
	<-stopChan
	log.Println("main: shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Println("main: http server shutdown error: %v", err)
	}
	log.Println("main: gracefully stopped")
}
