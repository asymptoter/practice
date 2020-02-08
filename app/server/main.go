package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/asymptoter/geochallenge-backend/apis/auth"
	"github.com/asymptoter/geochallenge-backend/base/config"
	"github.com/asymptoter/geochallenge-backend/base/db"
	_ "github.com/asymptoter/geochallenge-backend/base/email"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
	"github.com/jmoiron/sqlx"
)

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

func newHttpServer(db *sqlx.DB, redisClient *redis.Client) *http.Server {
	cfg := config.Value.Server
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	v1 := r.Group("/api/v1")
	auth.SetHttpHandler(v1.Group("/auth"), db, redisClient)

	return &http.Server{
		Addr:    cfg.Address,
		Handler: r,
	}
}

func main() {
	flag.Parse()
	config.Init()

	db, err := db.NewMySQL()
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
