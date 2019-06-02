package main

import (
	"database/sql"
	"log"
	"text/template"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	Name     string `form:"name"`
	Password string `form:"password"`
}

func signup(c *gin.Context) {
	var user User
	if c.ShouldBind(&user) == nil {
		log.Println(user.Name)
		log.Println(user.Password)

		// Connect to mysql
		db, err := sql.Open("mysql", "asymptoter:password@localhost:3306")
		if err != nil {
			log.Println(err)
		}

		log.Println("sql.Open OK1")

		if _, err = db.Exec("INSERT INTO users (name, password) VALUES (?, ?)", user.Name, user.Password); err != nil {
			log.Println(err)
		}

		log.Println("db.Exec OK")

		t, err := template.ParseFiles("./home.html")
		if err != nil {
			log.Println(err)
		}

		if err := t.Execute(c.Writer, nil); err != nil {
			log.Println(err)
		}
	}
}

func login(c *gin.Context) {
}

func home(c *gin.Context) {
	t, err := template.ParseFiles("./index.html")
	if err != nil {
		log.Println(err)
	}

	if err := t.Execute(c.Writer, nil); err != nil {
		log.Println(err)
	}
}

func main() {
	r := gin.Default()
	// Connect to mysql
	_, err := sql.Open("mysql", "asymptoter:password@localhost:3306")
	if err != nil {
		log.Println(err)
	}
	//fmt.Println("err:", err)
	//fmt.Println("XDDDDD")
	r.GET("/homea", home)
	r.POST("/signup", signup)
	r.POST("/login", login)
	r.Run(":8080")
}
