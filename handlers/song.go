package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
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
// @Param link query string false "Фильтр по ссылке"
// @Param text query string false "Фильтр по тексту или фрагменту тектса"
// @Param sort query string false "Поле для сортировки asc для возрастания и desc для убывания"
// @Param page query int false "Номер страницы(пагинация)" default(1)
// @Param limit query int false "Лимит записей на странице" default(10)
// @Success 200 {array} models.Song
// @Router /songs [get]
func (h *SongHandler) GetSongs(c *gin.Context) {

	page, err := strconv.Atoi(c.DefaultQuery("page", "1")) //Номер страницы (по умолчанию 1)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10")) //Количество песен на странице (по умолчанию 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}

	if page < 1 {
		page = 1
	}
	if limit < 0 {
		limit = 0
	}

	offset := limit * (page - 1)

	group := c.Query("group") //фильтр по группе
	song := c.Query("song")   //фмльтр по песне
	link := c.Query("link")   //фильтр по ссылке
	text := c.Query("text")   //фильтр по тексту

	text = strings.ReplaceAll(text, `\n`, "\n") // чтобы коректно находились записи в бд

	query := h.DB.Model(&models.Song{}).
		Joins("JOIN song_details ON song_details.song_id = songs.id").
		Preload("SongDetails")

	if group != "" {
		query = query.Where(`"group" = ?`, group) // Экранирую "group"
	}
	if song != "" {
		query = query.Where("song = ?", song)
	}
	if link != "" {
		query = query.Where("link = ?", link)
	}
	if text != "" {
		query = query.Where("text LIKE ?", "%"+text+"%") //LIKE для частичного совподения чтобы искать не по всему тексту а по фрагменту
	}

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
// @Param page query int false "Номер страницы(пагинация)" default(1)
// @Param limit query int false "Лимит куплетов на странице" default(5)
// @Success 200 {string} string "Текст песни"
// @Router /songs/{id}/text [get]
func (h *SongHandler) GetSongText(c *gin.Context) {
	songId := c.Param("id")

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "5"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}

	if page < 1 {
		page = 1
	}
	if limit < 0 {
		limit = 0
	}

	var songDetails models.SongDetails

	if h.DB.First(&songDetails, songId).Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Песня не найдена",
		})
		log.Println(err)
	}

	verses := strings.Split(songDetails.Text, "\n\n")

	totalVerses := len(verses)
	start := (page - 1) * limit
	end := start + limit

	if start > totalVerses {
		start = totalVerses
	}
	if end > totalVerses {
		end = totalVerses
	}

	paginationVerses := verses[start:end]

	c.JSON(http.StatusOK, gin.H{"text": paginationVerses})
}

// Редактировать песню по ID
// @Summary Редактировать песню
// @Description Редактировать песню по её ID
// @Tags songs
// @Accept json
// @Produce json
// @Param id path int true "ID песни"
// @Param song body models.Song false "Данные песни"
// @Param songDetails body models.SongDetails false "Дополнительные данные песни"
// @Success 200 {object} models.Song
// @Router /songs/{id} [put]
func (h *SongHandler) EditSong(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}

	var song models.Song

	// Десериализуем запрос
	if err := c.BindJSON(&song); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат данных"})
		log.Println(err)
		return
	}

	log.Printf("%+v\n", song)

	// releaseDate, err := time.Parse("02.01.2006", request.SongDetails.ReleaseDate) //преобразуем строку в time.Time
	// if err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{
	// 		"error": "Некоректная дата релиза",
	// 	})
	// 	log.Printf("ошибка парсинга даты релиза песни: %v", err)
	// }

	// song := models.Song{
	// 	Id:    id,
	// 	Song:  request.Song,
	// 	Group: request.Group,
	// }

	// songDetails := models.SongDetails{
	// 	SongId:      id,
	// 	Text:        request.SongDetails.Text,
	// 	ReleaseDate: releaseDate,
	// 	Link:        request.SongDetails.Link,
	// }

	log.Printf("%+v\n", song)

	if err := h.DB.First(&song, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Песня не найдена"})
		log.Println(err)
		return
	}

	log.Printf("%+v\n", song)

	// if err := h.DB.First(&songDetails, id).Error; err != nil {
	// 	c.JSON(http.StatusNotFound, gin.H{"error": "Детали песни не найдены"})
	// 	log.Println(err)
	// 	return
	// }

	// log.Printf("%+v\n", songDetails)

	// if err := h.DB.Save(&song).Error; err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось обновить песню"})
	// 	return
	// }

	// // Обновляем данные деталей песни с использованием Save
	// if err := h.DB.Save(&songDetails).Error; err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось обновить детали песни"})
	// 	return
	// }

	// if err := h.DB.Model(&song).Where("id = ?", id).Updates(&song).Error; err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось обновить песню"})
	// 	return
	// }

	// // Обновляем запись деталей песни
	// if err := h.DB.Model(&songDetails).Where("song_id = ?", id).Updates(&songDetails).Error; err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось обновить детали песни"})
	// 	return
	// }

	//c.JSON(http.StatusOK, song)
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
	var songDetails models.SongDetails

	tx := h.DB.Begin() //Начало транзакции
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": tx.Error.Error(),
		})
		log.Println(tx.Error)
		return
	}

	if tx.Delete(&songDetails, id).Error != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": tx.Error.Error(),
		})
		return
	}

	if tx.Delete(&song, id).Error != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": tx.Error.Error(),
		})
		return
	}

	if tx.Commit().Error != nil { //подтверждение транзакции
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": tx.Error.Error(),
		})
		return
	}
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

	err := c.BindJSON(&song)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		log.Println(err)
		return
	}

	songDetails, err := GetSongInfo(song.Song, song.Group)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		log.Println(err)
		return
	}

	tx := h.DB.Begin() //начало транзакции, чтобы данные добавлялись в обе таблицы атомарно
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": tx.Error.Error(),
		})
		log.Println(tx.Error)
		return
	}

	if tx.Create(&song).Error != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": tx.Error.Error(),
		})
		log.Println(tx.Error)
		return
	}

	songDetails.SongId = song.Id

	if tx.Create(&songDetails).Error != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": tx.Error.Error(),
		})
		log.Println(tx.Error)
		return
	}

	if tx.Commit().Error != nil { //подтверждение транзакции
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": tx.Error.Error(),
		})
		log.Println(tx.Error)
		return
	}

	song.SongDetails = *songDetails

	c.JSON(http.StatusOK, song)
}

func GetSongInfo(song string, group string) (*models.SongDetails, error) { //отправляет запрос на сторонее API для получения информации о песне

	var request struct {
		ReleaseDate string `json:"releaseDate"` // Используем строку для разбора
		Text        string `json:"text"`
		Link        string `json:"link"`
	}

	url := fmt.Sprintf("http://localhost:8081/info?song=%s&group=%s", url.QueryEscape(song), url.QueryEscape(song))

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении ответа от сервера: %v", err)
	}
	defer resp.Body.Close()

	log.Println("Получен ответ от API", url)

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении тела ответа от сервера: %v", err)
	}

	// if len(data) == 0 {
	// 	return nil, nil
	// }

	log.Printf("Получен тело ответа от API %v\n%v", url, string(data))

	err = json.Unmarshal(data, &request)
	if err != nil {
		return nil, fmt.Errorf("ошибка при анмаршалинге ответа от сервера: %v", err)
	}

	log.Printf("%+v\n", request)

	releaseDate, err := time.Parse("02.01.2006", request.ReleaseDate) //преобразуем строку в time.Time
	if err != nil {
		return nil, fmt.Errorf("ошибка парсинга даты релиза песни: %v", err)
	}

	songDetail := models.SongDetails{
		Text:        request.Text,
		ReleaseDate: releaseDate,
		Link:        request.Link,
	}

	return &songDetail, nil
}
