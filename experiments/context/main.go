package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
)

func serve(ctx context.Context) error {
	r := gin.Default()
	r.GET("/DoSomething", func(c *gin.Context) {
		for i := 0; i < 10; i++ {
			time.Sleep(time.Second)
			fmt.Println(i)
		}
		c.JSON(200, gin.H{
			"message": "done",
		})
	})

	httpServer := &http.Server{
		Handler: r,
	}

	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			log.Println("ListenAndServe failed", err)
		}
	}()

	log.Printf("server started")
	<-ctx.Done()
	log.Printf("server stopped")

	ctxShutDown, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer func() {
		cancel()
	}()

	if err := httpServer.Shutdown(ctxShutDown); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server Shutdown Failed:%+s", err)
	}

	log.Printf("server exited properly")

	return nil
}

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		oscall := <-c
		log.Printf("system call:%+v", oscall)
		cancel()
	}()

	if err := serve(ctx); err != nil {
		log.Printf("failed to serve:+%v\n", err)
	}

	/*
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer func() {
			cancel()
		}()
		stopChan := make(chan os.Signal, 1)
		signal.Notify(stopChan, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGTERM)
		<-stopChan
		log.Println("main: shutting down server...")
		if err := httpServer.Shutdown(ctx); err != nil {
			log.Fatalf("main: http server shutdown error: %v", err)
		}
		log.Println("main: gracefully stopped")
	*/
}
