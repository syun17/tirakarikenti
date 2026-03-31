package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	_ "backend/docs"
)

// ========================================
// Swagger API Info
// ========================================

// @title Tirakari API
// @version 1.0
// @description 部屋の画像解析と履歴管理を行うAPI
// @host localhost:8080
// @BasePath /

// ========================================
// GORM Models
// ========================================

type AppUser struct {
	ID       uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Password string `json:"password" binding:"required"`
}

func (AppUser) TableName() string {
	return "app_user"
}

type Room struct {
	ID       uint    `gorm:"primaryKey;autoIncrement" json:"id"`
	RoomName string  `json:"room_name" binding:"required"`
	Score    float64 `json:"score"`
	UserID   uint    `json:"user_id" binding:"required"`
}

func (Room) TableName() string {
	return "room"
}

type Img struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Score     float64   `json:"score"`
	Note      string    `json:"note"` // 物体検出結果などのメモ
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	RoomID    uint      `json:"room_id"`
}

func (Img) TableName() string {
	return "img"
}

type Detail struct {
	ImgID  uint `gorm:"primaryKey" json:"img_id"`
	ObjID  uint `gorm:"primaryKey" json:"obj_id"`
	ObjCnt int  `json:"obj_cnt"`
}

func (Detail) TableName() string {
	return "detail"
}

// ========================================
// Request / Response structs (Swagger用)
// ========================================

type LoginRequest struct {
	ID       uint   `json:"id" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	Password string `json:"password" binding:"required"`
}

type RoomCreateRequest struct {
	RoomName string `json:"room_name" binding:"required"`
	UserID   uint   `json:"user_id" binding:"required"`
}

type MyRoomsRequest struct {
	UserID uint `json:"user_id" binding:"required"`
}

type RoomHistoryRequest struct {
	RoomID uint `json:"room_id" binding:"required"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type LoginResponse struct {
	Message string `json:"message"`
	UserID  uint   `json:"user_id"`
}

type MyRoomItem struct {
	RoomID   uint   `json:"room_id"`
	RoomName string `json:"room_name"`
}

type AnalyzeResponse struct {
	ImgID uint           `json:"img_id"`
	Score float64        `json:"score"`
	YOLO  map[string]int `json:"yolo"`
}

type PythonAnalyzeResult struct {
	Score float64        `json:"score"`
	YOLO  map[string]int `json:"yolo"`
}

// ========================================
// DB
// ========================================

var db *gorm.DB

func initDB() {
	dsn := os.Getenv("DATABASE_URL")
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("DB接続失敗:", err)
	}

	// 自動マイグレーション
	db.AutoMigrate(&Img{})
}

// ========================================
// Handlers
// ========================================

// Login godoc
// @Summary ログイン
// @Description ユーザーIDとパスワードでログイン認証を行う
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "ログイン情報"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /login [post]
func LoginHandler(c *gin.Context) {
	var input LoginRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	var user AppUser
	if err := db.First(&user, "id = ? AND password = ?", input.ID, input.Password).Error; err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "IDまたはパスワードが正しくありません"})
		return
	}

	c.JSON(http.StatusOK, LoginResponse{
		Message: "ログイン成功",
		UserID:  user.ID,
	})
}

// Register godoc
// @Summary ユーザー登録
// @Description 新規ユーザーを登録する
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "登録情報"
// @Success 201 {object} AppUser
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /register [post]
func RegisterHandler(c *gin.Context) {
	var user AppUser
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	if err := db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "ユーザー登録に失敗しました"})
		return
	}

	c.JSON(http.StatusCreated, user)
}

// CreateRoom godoc
// @Summary 部屋作成
// @Description 新しい部屋を作成する
// @Tags rooms
// @Accept json
// @Produce json
// @Param request body RoomCreateRequest true "部屋情報"
// @Success 201 {object} Room
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /rooms [post]
func CreateRoomHandler(c *gin.Context) {
	var room Room
	if err := c.ShouldBindJSON(&room); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	if err := db.Create(&room).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "部屋の登録に失敗しました"})
		return
	}

	c.JSON(http.StatusCreated, room)
}

// GetMyRooms godoc
// @Summary 自分の部屋一覧取得
// @Description 指定したユーザーの部屋一覧を取得する
// @Tags rooms
// @Accept json
// @Produce json
// @Param request body MyRoomsRequest true "ユーザーID"
// @Success 200 {array} MyRoomItem
// @Failure 400 {object} ErrorResponse
// @Router /rooms/my-rooms [post]
func GetMyRoomsHandler(c *gin.Context) {
	var input MyRoomsRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	var rooms []Room
	db.Where("user_id = ?", input.UserID).Find(&rooms)

	var result []MyRoomItem
	for _, room := range rooms {
		result = append(result, MyRoomItem{
			RoomID:   room.ID,
			RoomName: room.RoomName,
		})
	}

	c.JSON(http.StatusOK, result)
}

// GetRoomHistory godoc
// @Summary 部屋の履歴取得
// @Description room_idを基に過去の画像スコアと作成日を取得する
// @Tags rooms
// @Accept json
// @Produce json
// @Param request body RoomHistoryRequest true "部屋ID"
// @Success 200 {array} Img
// @Failure 400 {object} ErrorResponse
// @Router /rooms/history [post]
func GetRoomHistoryHandler(c *gin.Context) {
	var input RoomHistoryRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	var imgs []Img
	db.Where("room_id = ?", input.RoomID).Order("created_at desc").Find(&imgs)
	c.JSON(http.StatusOK, imgs)
}

// AnalyzePython godoc
// @Summary 画像解析
// @Description 画像をPython APIに送信し、スコアと検出結果を保存する
// @Tags analyze
// @Accept multipart/form-data
// @Produce json
// @Param room_id formData int true "部屋ID"
// @Param image formData file true "アップロード画像"
// @Success 200 {object} AnalyzeResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /analyze-python [post]
func AnalyzePythonHandler(c *gin.Context) {
	roomIDStr := c.PostForm("room_id")
	if roomIDStr == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "room_id is required"})
		return
	}

	fileHeader, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "image file is required"})
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to open image"})
		return
	}
	defer file.Close()

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	if err := writer.WriteField("room_id", roomIDStr); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to write room_id"})
		return
	}

	part, err := writer.CreateFormFile("file", fileHeader.Filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to create form file"})
		return
	}

	if _, err = io.Copy(part, file); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to copy file"})
		return
	}

	writer.Close()

	resp, err := http.Post(
		"http://ai:8000/analyze/",
		writer.FormDataContentType(),
		&buf,
	)
	if err != nil {
		log.Printf("failed to call python api: %v", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to call python api"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		log.Printf("python api returned status %d: %s", resp.StatusCode, string(bodyBytes))
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "python api returned error"})
		return
	}

	var pythonResult PythonAnalyzeResult
	if err := json.NewDecoder(resp.Body).Decode(&pythonResult); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to decode response"})
		return
	}

	var roomID uint
	if _, err := fmt.Sscanf(roomIDStr, "%d", &roomID); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid room_id"})
		return
	}

	// YOLOの結果をメモ用の文字列に変換
	var noteStr string
	for name, count := range pythonResult.YOLO {
		if noteStr != "" {
			noteStr += ", "
		}
		noteStr += fmt.Sprintf("%s: %d", name, count)
	}

	img := Img{
		Score:  pythonResult.Score,
		Note:   noteStr,
		RoomID: roomID,
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&img).Error; err != nil {
			return err
		}

		objMap := map[string]uint{
			"bottle":     1,
			"cup":        2,
			"wine glass": 3,
			"bowl":       4,
			"can":        5,
		}

		for name, count := range pythonResult.YOLO {
			if objID, ok := objMap[name]; ok {
				detail := Detail{
					ImgID:  img.ID,
					ObjID:  objID,
					ObjCnt: count,
				}
				if err := tx.Create(&detail).Error; err != nil {
					return err
				}
			}
		}

		if err := tx.Model(&Room{}).Where("id = ?", roomID).Update("score", pythonResult.Score).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to save to database"})
		return
	}

	c.JSON(http.StatusOK, AnalyzeResponse{
		ImgID: img.ID,
		Score: img.Score,
		YOLO:  pythonResult.YOLO,
	})
}

// ========================================
// Main
// ========================================

// DeleteImage godoc
// @Summary 画像の削除
// @Description 指定されたIDの画像を削除する
// @Tags images
// @Accept json
// @Produce json
// @Param id path int true "画像ID"
// @Success 200 {object} MessageResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /images/{id} [delete]
func DeleteImageHandler(c *gin.Context) {
	id := c.Param("id")
	if err := db.Delete(&Img{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "削除に失敗しました"})
		return
	}
	c.JSON(http.StatusOK, MessageResponse{Message: "削除しました"})
}

func main() {
	initDB()

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		AllowCredentials: false,
	}))

	// Swagger UI
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API Routes
	r.POST("/login", LoginHandler)
	r.POST("/register", RegisterHandler)
	r.POST("/rooms", CreateRoomHandler)
	r.POST("/rooms/my-rooms", GetMyRoomsHandler)
	r.POST("/rooms/history", GetRoomHistoryHandler)
	r.POST("/analyze-python", AnalyzePythonHandler)
	r.DELETE("/images/:id", DeleteImageHandler)

	r.Run(":8080")
}