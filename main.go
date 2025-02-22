package main

import (
	"fmt"
	"log"
	"os"
	_ "song-library/docs"
	"song-library/handlers"
	"song-library/models"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// @Song Library API
// @version 1.0
// @description API для управления библиотекой песен.

// @host localhost:8080
// @BasePath /
func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Не удалось загрузить файл конфигурации")
	}

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	database := os.Getenv("DB_NAME")
	password := os.Getenv("DB_PASSWORD")

	connStr := fmt.Sprintf("host=%v port=%v user=%v database=%v password=%v sslmode=disable",
		host,
		port,
		user,
		database,
		password)

	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		log.Fatal("Не удалось подключиться к базе данных", err)
	}

	db.Debug()

	db.AutoMigrate(&models.Song{})

	songHandler := handlers.NewSongHandler(db)

	r := gin.Default()

	r.GET("/songs", songHandler.GetSongs)

	r.GET("/songs/:id/text", songHandler.GetSongText)

	r.PUT("/songs/:id", songHandler.EditSong)

	r.DELETE("/songs/:id", songHandler.DeleteSong)

	r.POST("/songs", songHandler.AddSong)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Run(":8080")
}
