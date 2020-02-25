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

	authApi "github.com/asymptoter/practice-backend/apis/auth"
	triviaApi "github.com/asymptoter/practice-backend/apis/trivia"
	"github.com/asymptoter/practice-backend/base/config"
	"github.com/asymptoter/practice-backend/base/db"
	_ "github.com/asymptoter/practice-backend/base/email"
	"github.com/asymptoter/practice-backend/base/redis"
	"github.com/asymptoter/practice-backend/store/trivia"
	"github.com/asymptoter/practice-backend/store/user"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func newHttpServer(db *sqlx.DB, redisService redis.Service) *http.Server {
	cfg := config.Value.Server
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	v1 := r.Group("/api/v1")
	userStore := user.NewStore(db, redisService)
	triviaStore := trivia.NewStore(db, redisService)
	authApi.SetHttpHandler(v1.Group("/auth"), db, redisService, userStore)
	triviaApi.SetHttpHandler(v1.Group("/trivia"), db, redisService, triviaStore, userStore)

	return &http.Server{
		Addr:    cfg.Address,
		Handler: r,
	}
}

func main() {
	flag.Parse()
	config.Init()

	db := db.MustNew("postgres")
	defer db.Close()

	redisService := redis.NewService()

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
