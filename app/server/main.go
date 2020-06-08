package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	authApi "github.com/asymptoter/practice-backend/apis/auth"
	triviaApi "github.com/asymptoter/practice-backend/apis/trivia"
	"github.com/asymptoter/practice-backend/base/config"
	"github.com/asymptoter/practice-backend/base/db"
	_ "github.com/asymptoter/practice-backend/base/email"
	"github.com/asymptoter/practice-backend/base/redis"
	authStore "github.com/asymptoter/practice-backend/store/auth"
	triviaStore "github.com/asymptoter/practice-backend/store/trivia"
	userStore "github.com/asymptoter/practice-backend/store/user"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/robfig/cron/v3"
)

func newHttpServer(db *sqlx.DB, redisService redis.Service) *http.Server {
	cfg := config.Value.Server
	r := gin.Default()
	/*
		r.GET("/", func(c *gin.Context) {
			c.Writer.Write([]byte("<!doctype html><html><head><title>This is the title of the webpage!</title></head><body>OK</body></html>"))
		})

		r.GET("/index.html", func(c *gin.Context) {
			c.Writer.Write([]byte("<!doctype html><html><head><title>This is the title of the webpage!</title></head><body>OK</body></html>"))
		})
	*/
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	v1 := r.Group("/v1")
	userStore := userStore.New(db, redisService)
	triviaStore := triviaStore.New(db, redisService)
	authStore := authStore.New(userStore)
	authApi.SetHttpHandler(v1.Group("/auth"), db, redisService, userStore, authStore)
	triviaApi.SetHttpHandler(v1.Group("/trivia"), triviaStore, userStore)

	return &http.Server{
		Addr:         cfg.Address,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
}

func main() {
	flag.Parse()
	pwd, err := exec.Command("pwd").Output()
	if err != nil {
		log.Fatal("exec.Command failed", err)
	}

	cfg := config.Init(string(pwd))
	cronJob := cron.New(cron.WithSeconds())
	cronJob.AddFunc("*/30 * * * * *", func() {
		config.Init(string(pwd))
	})
	go func() {
		cronJob.Start()
	}()

	dbCfg := cfg.Database
	db := db.MustNew("postgres", fmt.Sprintf(dbCfg.ConnectionString, dbCfg.Host, dbCfg.Port))
	defer db.Close()

	redisCfg := cfg.Redis
	redisService := redis.NewService(fmt.Sprintf(redisCfg.ConnectionString, redisCfg.Host, redisCfg.Port))

	httpServer := newHttpServer(db, redisService)
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
