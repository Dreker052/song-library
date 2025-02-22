package handlers

import (
	"net/http"

	_ "song-library/docs"
	"song-library/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SongHandler struct {
	DB *gorm.DB
}

func NewSongHandler(db *gorm.DB) *SongHandler {
	return &SongHandler{DB: db}
}

// Получить все песни
// @Summary Получить все песни
// @Description Получить список всех песен
// @Tags songs
// @Accept json
// @Produce json
// @Success 200 {array} models.Song
// @Router /songs [get]
func (h *SongHandler) GetSongs(c *gin.Context) {
	var songs []models.Song
	h.DB.Find(&songs)
	c.JSON(http.StatusOK, songs)
}

// Получить текст песни по ID
// @Summary Получить текст песни
// @Description Получить текст песни по её ID
// @Tags songs
// @Accept json
// @Produce json
// @Param id path int true "ID песни"
// @Success 200 {string} string "Текст песни"
// @Router /songs/{id}/text [get]
func (h *SongHandler) GetSongText(c *gin.Context) {
	id := c.Param("id")
	var song models.Song
	h.DB.First(&song, id)
	c.JSON(http.StatusOK, gin.H{"text": song.Text})
}

// Редактировать песню по ID
// @Summary Редактировать песню
// @Description Редактировать песню по её ID
// @Tags songs
// @Accept json
// @Produce json
// @Param id path int true "ID песни"
// @Param song body models.Song true "Данные песни"
// @Success 200 {object} models.Song
// @Router /songs/{id} [put]
func (h *SongHandler) EditSong(c *gin.Context) {
	id := c.Param("id")
	var song models.Song
	h.DB.First(&song, id)
	c.BindJSON(&song)
	h.DB.Save(&song)
	c.JSON(http.StatusOK, song)
}

// Удалить песню по ID
// @Summary Удалить песню
// @Description Удалить песню по её ID
// @Tags songs
// @Accept json
// @Produce json
// @Param id path int true "ID песни"
// @Success 200 {string} string "Песня успешно удалена"
// @Router /songs/{id} [delete]
func (h *SongHandler) DeleteSong(c *gin.Context) {
	id := c.Param("id")
	var song models.Song
	h.DB.Delete(&song, id)
	c.JSON(http.StatusOK, gin.H{"message": "Песня удалена"})
}

// Добавить новую песню
// @Summary Добавить песню
// @Description Добавить новую песню
// @Tags songs
// @Accept json
// @Produce json
// @Param song body models.Song true "Данные песни"
// @Success 201 {object} models.Song
// @Router /songs [post]
func (h *SongHandler) AddSong(c *gin.Context) {
	var song models.Song
	c.BindJSON(&song)
	h.DB.Create(&song)
	c.JSON(http.StatusOK, song)
}
