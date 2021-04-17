package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/joho/godotenv"
)

type Message struct {
	gorm.Model
	Title       string `json:"title"`
	MessageText string `json:"message_text"`
}

func NewMessage() Message {
	return Message{}
}

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func gormConnect() *gorm.DB {
	DBName := os.Getenv("DB_NAME")
	DBUser := os.Getenv("DB_USER")
	DBPass := os.Getenv("DB_PASS")
	fmt.Println(DBName)

	connectTemplate := "%s:%s@%s/%s"
	connect := fmt.Sprintf(connectTemplate, DBUser, DBPass, "tcp(mysql:3306)", DBName)
	db, err := gorm.Open("mysql", connect)

	if err != nil {
		log.Println(err.Error())
	}

	return db
}

func setRouter(db *gorm.DB) *gin.Engine {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*.html")

	data := "Hello Go/Gin!!"

	router.GET("/", func(ctx *gin.Context) {
		ctx.HTML(200, "index.html", gin.H{"data": data})
	})

	//CREATE
	router.POST("/message", func(c *gin.Context) {
		data := NewMessage()
		now := time.Now()
		data.CreatedAt = now
		data.UpdatedAt = now

		if err := c.BindJSON(&data); err != nil {
			c.String(http.StatusBadRequest, "Request is failed: "+err.Error())
		}
		db.NewRecord(data)
		db.Create(&data)
		if db.NewRecord(data) == false {
			c.JSON(http.StatusOK, data)
		}
	})

	//READ
	//全レコード
	router.GET("/messages", func(c *gin.Context) {
		messages := []Message{}
		db.Find(&messages)
		c.JSON(http.StatusOK, messages)
	})

	return router
}

func main() {
	loadEnv()
	db := gormConnect()
	router := setRouter(db)

	defer db.Close()

	db.Set("gorm:table_options", "ENGINE=InnoDB")
	db.AutoMigrate(&Message{})
	db.LogMode(true)

	router.Run()
}
