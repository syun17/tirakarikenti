package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"

    "gorm.io/driver/postgres"
    "gorm.io/gorm"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type ScoreHistory struct {
	ID        string `dynamodbav:"id" json:"id"`
	Score     int    `dynamodbav:"score" json:"score"`
	ImageURL  string `dynamodbav:"image_url" json:"image_url"`
	CreatedAt string `dynamodbav:"created_at" json:"created_at"`
}

type User struct {
    gorm.Model
    Name  string `json:"name"`
    Email string `json:"email" gorm:"unique"`
}

var db *gorm.DB

func initDB() {
    dsn := os.Getenv("DATABASE_URL")
    var err error
    db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal("DB接続失敗:", err)
    }
    db.AutoMigrate(&User{})
}


func main() {
	initDB()

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		AllowCredentials: true,
	}))

	r.POST("/analyze-python", func(c *gin.Context) {
    	fileHeader, err := c.FormFile("image")
    	if err != nil {
    	    c.JSON(http.StatusBadRequest, gin.H{"error": "image file is required"})
    	    return
    	}

    	file, err := fileHeader.Open()
    	if err != nil {
    	    c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to open image"})
    	    return
    	}
    	defer file.Close()

    	// multipart bodyを構築
    	var buf bytes.Buffer
    	writer := multipart.NewWriter(&buf)

    	part, err := writer.CreateFormFile("file", fileHeader.Filename) // FastAPIのパラメータ名 "file" に合わせる
    	if err != nil {
    	    c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create form file"})
    	    return
    	}

    	if _, err = io.Copy(part, file); err != nil {
    	    c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to copy file"})
    	    return
    	}
    	writer.Close()

    	resp, err := http.Post(
    	    "http://ai:8000/analyze/",  // タイポ修正 + dockerネットワーク内のホスト名
    	    writer.FormDataContentType(),
    	    &buf,
    	)
    	if err != nil {
    	    log.Printf("failed to call python api: %v", err)
    	    c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to call python api"})
    	    return
    	}
    	defer resp.Body.Close()

    	// レスポンスをデコード（抜けていた）
    	var result map[string]interface{}
    	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
    	    c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to decode response"})
    	    return
    	}

    	c.JSON(http.StatusOK, result)
	})

	r.GET("/users", func(c *gin.Context) {
        var users []User
        db.Find(&users)
        c.JSON(200, users)
    })

    r.POST("/users", func(c *gin.Context) {
        var user User
        if err := c.ShouldBindJSON(&user); err != nil {
            c.JSON(400, gin.H{"error": err.Error()})
            return
        }
        db.Create(&user)
        c.JSON(201, user)
    })

	r.Run(":8080")
}
