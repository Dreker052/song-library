package handlers

import (
	"log"
	"net/http"
	"strconv"
	"time"

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
// @Param group query string false "Фильтр по группе"
// @Param song query string false "Фильтр по названию песни"
// @Param sort query string false "Поле для сортировки"
// @Param page query int false "Номер страницы" default(1)
// @Param limit query int false "Лимит записей на странице" default(10)
// @Success 200 {array} models.Song
// @Router /songs [get]
func (h *SongHandler) GetSongs(c *gin.Context) {

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))    //Номер страницы (по умолчанию 1)
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10")) //Количество песен на странице (по умолчанию 10)
	offset := limit * (page - 1)

	group := c.Query("group") //фильтр по группе
	song := c.Query("song")   //фмльтр по песне
	// releaseDate := c.Query("releaseDate") //фильтр по дате релиза

	query := h.DB.Model(&models.Song{})

	if group != "" {
		query = query.Where(`"group" = ?`, group) // Экранирую "group"
	}
	if song != "" {
		query = query.Where("song = ?", song)
	}
	// if releaseDate != "" {
	// 	query = query.Where("releaseDate = ?", releaseDate)
	// }

	sortOrder := c.Query("sort")

	if sortOrder == "asc" {
		query = query.Order("release_date ASC") //сортировка по возростанию
	} else if sortOrder == "desc" {
		query = query.Order("release_date DESC") //сортировка по убыванию
	} else {
		query = query.Order("release_date ASC") //по умолчанию сортировка по возростанию
	}

	var songs []models.Song
	query.Offset(offset).Limit(limit).Find(&songs)

	c.JSON(http.StatusOK, gin.H{
		"page":  page,
		"limit": limit,
		"sort":  sortOrder,
		"songs": songs,
	})
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

	var request struct {
		Group       string `json:"group"`
		Song        string `json:"song"`
		ReleaseDate string `json:"releaseDate"` // Используем строку для разбора
		Text        string `json:"text"`
		Link        string `json:"link"`
	}

	err := c.BindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		log.Println(err)
		return
	}

	releaseDate, err := time.Parse("02.01.2006", request.ReleaseDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		log.Println(err)
		return
	}

	song := models.Song{
		Group:       request.Group,
		Song:        request.Song,
		ReleaseDate: releaseDate,
		Text:        request.Text,
		Link:        request.Link,
	}

	h.DB.Create(&song)
	c.JSON(http.StatusOK, song)
}
