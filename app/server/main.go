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

	"github.com/asymptoter/practice-backend/base/ctx"
	_ "github.com/asymptoter/practice-backend/external/email"

	"github.com/gin-gonic/gin"
)

func newHttpServer() *http.Server {
	r := gin.Default()
	v1 := r.Group("/api/v1")
	v1.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	return &http.Server{
		Addr:         "localhost:8080",
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
}

func main() {
	flag.Parse()
	ctx := ctx.Background()

	/*
		pwd, err := exec.Command("pwd").Output()
		if err != nil {
			log.Fatal("exec.Command failed", err)
		}

		cfg := config.Init(ctx, string(pwd))

		dbCfg := cfg.Database
		db := db.MustNew("postgres", fmt.Sprintf(dbCfg.ConnectionString, dbCfg.Host, dbCfg.Port))
		defer db.Close()

		redisCfg := cfg.Redis
		redisService := redis.NewService(fmt.Sprintf(redisCfg.ConnectionString, redisCfg.Host, redisCfg.Port))
	*/

	httpServer := newHttpServer()
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
	gctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(gctx); err != nil {
		log.Println("main: http server shutdown error: %v", err)
	}
	log.Println("main: gracefully stopped")
}
