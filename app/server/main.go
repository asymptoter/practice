package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"text/template"

	"github.com/asymptoter/geochallenge-backend/apis/auth"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"google.golang.org/grpc"
	yaml "gopkg.in/yaml.v2"
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

type serverConfiguration struct {
	Address string `yaml:"address"`
}

type mysqlConfiguration struct {
	Address      string `yaml:"address"`
	DatabaseName string `yaml:"databaseName"`
	Username     string `yaml:"username"`
	Password     string `yaml:"password"`
}

type redisConfiguration struct {
	Address string `yaml:"address"`
}

type configuration struct {
	Server serverConfiguration `yaml:"server"`
	MySQL  mysqlConfiguration  `yaml:"mysql"`
	Redis  redisConfiguration  `yaml:"redis"`
}

type configurations map[string]configuration

func getConfigYaml(env string) *configuration {
	c := configurations{}
	file, err := ioutil.ReadFile("../../config/config.yml")
	if err != nil {
		log.Println("ioutil.ReadFile failed ", err)
		return nil
	}

	if err := yaml.Unmarshal(file, c); err != nil {
		log.Println("yaml.Unmarshal failed ", err)
		return nil
	}

	res := c[env]

	return &res
}

func setupMySQL(cfg mysqlConfiguration) (*gorm.DB, error) {
	connectionString := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=true&multiStatements=true", cfg.Username, cfg.Password, cfg.Address, cfg.DatabaseName)
	fmt.Println(connectionString)
	db, err := gorm.Open("mysql", connectionString)
	if err != nil {
		log.Println("gorm.Open failed ", err)
		return nil, err
	}

	return db, nil
}

func newHttpServer(cfg serverConfiguration, db *gorm.DB) *http.Server {
	r := gin.New()

	auth.SetHttpHandler(r, db)
	r.GET("/home", home)

	return &http.Server{
		Addr:    cfg.Address,
		Handler: r,
	}
}

func main() {
	flag.Parse()
	cfg := getConfigYaml("local")

	db, err := setupMySQL(cfg.MySQL)
	if err != nil {
		log.Println("setup MySQL failed ", err)
		return
	}

	httpServer := newHttpServer(cfg.Server, db)

	// Start http server
	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			log.Println("ListenAndServe failed ", err)
		}
	}()

	fmt.Println("XD")
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGTERM)
	<-stopChan
	log.Println("main: shutting down server...")
	grpcServer.GracefulStop()
	log.Println("main: gracefully stopped")
}
