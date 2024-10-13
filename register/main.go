package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB

// Конфигурация базы данных PostgreSQL
const (
	host     = "localhost"
	port     = 5432
	user     = "mirak"
	password = "mirak1991"
	dbname   = "postgresql"
)

func main() {
	var err error

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Successfully connected to the database!")

	router := gin.Default()

	router.Static("/static", "./static")

	router.GET("/register", showRegisterPage)

	router.POST("/register", registerUser)

	router.Run(":8080")
}

func showRegisterPage(c *gin.Context) {
	tmpl, _ := template.ParseFiles("templates/register.html")
	tmpl.Execute(c.Writer, nil)
}

func registerUser(c *gin.Context) {
	username := c.PostForm("username")
	email := c.PostForm("email")
	password := c.PostForm("password")

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error while hashing password")
		return
	}

	_, err = db.Exec("INSERT INTO users (username, email, password_hash) VALUES ($1, $2, $3)", username, email, string(hashedPassword))
	if err != nil {
		c.String(http.StatusInternalServerError, "Error while saving user to the database")
		return
	}

	c.String(http.StatusOK, "Registration successful!")
}
