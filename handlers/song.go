package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
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
// @Failure 400 {string} string "Неверный формат параметров запроса"
// @Failure 500 {string} string "Внутренняя ошибка сервера"
// @Router /songs [get]
func (h *SongHandler) GetSongs(c *gin.Context) {

	page, err := strconv.Atoi(c.DefaultQuery("page", "1")) //Номер страницы (по умолчанию 1)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Некоректное значение page",
		})
		log.Printf("Некоректное значение page, %v", err)
		return
	}
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10")) //Количество песен на странице (по умолчанию 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Некоректное значение limit",
		})
		log.Printf("Некоректное значение limit, %v", err)
		return
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
		query = query.Order("TO_DATE(release_date, 'DD.MM.YYYY') ASC") //сортировка по возростанию
	} else if sortOrder == "desc" {
		query = query.Order("TO_DATE(release_date, 'DD.MM.YYYY') DESC") //сортировка по убыванию
	} else {
		query = query.Order("TO_DATE(release_date, 'DD.MM.YYYY') ASC") //по умолчанию сортировка по возростанию
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
// @Failure 400 {string} string "Неверный формат параметров запроса"
// @Failure 404 {string} string "Песня не найдена"
// @Failure 404 {string} string "Текст песни отсутствует"
// @Failure 500 {string} string "Внутренняя ошибка сервера"
// @Router /songs/{id}/text [get]
func (h *SongHandler) GetSongText(c *gin.Context) {
	songId := c.Param("id")

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Некоректное значение page",
		})
		log.Printf("Некоректное значение page, %v", err)
		return
	}
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "5"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Некоректное значение limit",
		})
		log.Printf("Некоректное значение limit, %v", err)
		return
	}

	if page < 1 {
		page = 1
	}
	if limit < 0 {
		limit = 10
	}

	var songDetails models.SongDetails

	if h.DB.First(&songDetails, songId).Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Песня не найдена",
		})
		log.Println(err)
		return
	}

	if songDetails.Text == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "Текст песни отсутствует"})
		return
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
// @Param song body models.SongWithDetails false "Данные песни"
// @Success 200 {object} models.Song
// @Failure 400 {string} string "Неверный формат данных"
// @Failure 400 {string} string "Неверный формат даты. Ожидаемый формат: DD.MM.YYYY"
// @Failure 404 {string} string "Песня не найдена"
// @Failure 404 {string} string "Детали песни не найдены"
// @Failure 500 {string} string "Внутренняя ошибка сервера"
// @Router /songs/{id} [put]
func (h *SongHandler) EditSong(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Некоректное значение id",
		})
		log.Printf("Некоректное значение id, %v\n", err)
		return
	}

	var songWithDetails models.SongWithDetails //использую специальную структуру для получения тела запроса

	if err := c.BindJSON(&songWithDetails); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат данных"})
		log.Println(err)
		return
	}

	if songWithDetails.SongDetails.ReleaseDate != "" {
		_, err = time.Parse("02.01.2006", songWithDetails.SongDetails.ReleaseDate) //проверяем коректность введеной даты
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат даты. Ожидаемый формат: DD.MM.YYYY"})
			log.Printf("Неверный формат даты, %v\n", err)
			return
		}
	}

	log.Printf("Тело запроса: %+v\n", songWithDetails)

	var song models.Song
	var songDetails models.SongDetails

	if err := h.DB.First(&song, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Песня не найдена"})
		log.Printf("Детали песни не найдены, %v", err)
		return
	}

	if err := h.DB.First(&songDetails, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Детали песни не найдены"})
		log.Printf("Детали песни не найдены, %v", err)
		return
	}

	song = models.Song{
		Id:    id,
		Song:  songWithDetails.Song,
		Group: songWithDetails.Group,
	}

	songDetails = models.SongDetails{
		SongId:      id,
		Text:        songWithDetails.SongDetails.Text,
		Link:        songWithDetails.SongDetails.Link,
		ReleaseDate: songWithDetails.SongDetails.ReleaseDate,
	}

	tx := h.DB.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось обновить данные песни"})
		log.Printf("Детали песни не найдены, %v\n", tx.Error.Error())
		return
	}

	if err := tx.Model(&song).Where("id = ?", id).Updates(&song).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось обновить данные песни"})
		log.Printf("Детали песни не найдены, %v\n", err)
		return
	}

	if err := tx.Model(&songDetails).Where("song_id = ?", id).Updates(&songDetails).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось обновить детали песни"})
		log.Printf("Детали песни не найдены, %v\n", err)
		return
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		log.Printf("Детали песни не найдены, %v\n", err)
		return
	}

	if err := h.DB.First(&song, id).Error; err != nil { //для вывода обновленных данных
		c.JSON(http.StatusNotFound, gin.H{"error": "Песня не найдена"})
		log.Printf("Детали песни не найдены, %v\n", err)
		return
	}

	if err := h.DB.First(&songDetails, id).Error; err != nil { //для вывода обновленных данных
		c.JSON(http.StatusNotFound, gin.H{"error": "Детали песни не найдены"})
		log.Printf("Детали песни не найдены, %v\n", err)
		return
	}

	songWithDetails = models.SongWithDetails{
		Song:  song.Song,
		Group: song.Group,
		SongDetails: models.SongDetails{
			SongId:      songDetails.SongId,
			Text:        songDetails.Text,
			Link:        songDetails.Link,
			ReleaseDate: songDetails.ReleaseDate,
		},
	}

	c.JSON(http.StatusOK, songWithDetails)
}

// Удалить песню по ID
// @Summary Удалить песню
// @Description Удалить песню по её ID
// @Tags songs
// @Accept json
// @Produce json
// @Param id path int true "ID песни"
// @Success 200 {string} string "Песня успешно удалена"
// @Failure 404 {string} string "Песня не найдена"
// @Failure 500 {string} string "Ошибка при удалении песни"
// @Router /songs/{id} [delete]
func (h *SongHandler) DeleteSong(c *gin.Context) {
	id := c.Param("id")
	var song models.Song
	var songDetails models.SongDetails

	if err := h.DB.First(&songDetails, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Песня не найдена"})
		log.Printf("Песня не найдена, %v\n", err)
		return
	}

	tx := h.DB.Begin() //Начало транзакции
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Не удалось удалить данные песни",
		})
		log.Printf("Не удалось удалить данные песни, %v\n", tx.Error)
		return
	}

	if tx.Delete(&songDetails, id).Error != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Не удалось удалить данные песни",
		})
		log.Printf("Не удалось удалить данные песни, %v\n", tx.Error)
		return
	}

	if tx.Delete(&song, id).Error != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Не удалось удалить данные песни",
		})
		log.Printf("Не удалось удалить данные песни, %v\n", tx.Error)
		return
	}

	if tx.Commit().Error != nil { //подтверждение транзакции
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Не удалось удалить данные песни",
		})
		log.Printf("Не удалось удалить данные песни, %v\n", tx.Error)
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
// @Failure 400 {string} string "Неверный формат данных"
// @Failure 500 {string} string "Ошибка при добавлении песни"
// @Failure 500 {string} string "Ошибка при получении данных от внешнего API"
// @Failure 400 {string} string "Песня уже добавлена"
// @Router /songs [post]
func (h *SongHandler) AddSong(c *gin.Context) {

	var song models.Song

	err := c.BindJSON(&song)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат данных"})
		log.Printf("Неверный формат данных, %v\n", err)
		return
	}

	if h.DB.First(&song).Error == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Песня уже добавлена"})
		log.Println("Песня уже добавлена")
		return
	}

	songDetails, err := GetSongInfo(song.Song, song.Group) //получаем доп данные песни
	if err != nil {
		log.Println(err)
	}

	tx := h.DB.Begin() //начало транзакции, чтобы данные добавлялись в обе таблицы атомарно
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Не удалось добавить песню",
		})
		log.Printf("Не удалось добавить данные песни, %v\n", tx.Error)
		return
	}

	if tx.Create(&song).Error != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Не удалось добавить песню",
		})
		log.Printf("Не удалось добавить данные песни, %v\n", tx.Error)
		return
	}

	songDetails.SongId = song.Id

	if tx.Create(&songDetails).Error != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Не удалось добавить песню",
		})
		log.Printf("Не удалось добавить данные песни, %v\n", tx.Error)
		return
	}

	if tx.Commit().Error != nil { //подтверждение транзакции
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Не удалось добавить песню",
		})
		log.Printf("Не удалось добавить данные песни, %v\n", tx.Error)
		return
	}

	song.SongDetails = *songDetails

	c.JSON(http.StatusOK, song)
}

func GetSongInfo(song string, group string) (*models.SongDetails, error) { //отправляет запрос на сторонее API для получения информации о песне

	var songDetails models.SongDetails

	url := fmt.Sprintf("%s/info?song=%s&group=%s", os.Getenv("API_DOMAIN"), url.QueryEscape(song), url.QueryEscape(song))

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении ответа от сервера: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ошибка при получении данных от API: статус %d", resp.StatusCode)
	}

	log.Println("Получен ответ от API", url)

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении тела ответа от сервера: %v", err)
	}

	if len(data) == 0 { //на случай если в API не будет доп данных о песне, мы их просто не добавляем
		songDetails = models.SongDetails{Link: "", Text: "", ReleaseDate: ""}
		return &songDetails, fmt.Errorf("дополнительные данные песни не найдены")
	}

	log.Printf("Получен тело ответа от API %v\n%v", url, string(data))

	err = json.Unmarshal(data, &songDetails)
	if err != nil {
		return nil, fmt.Errorf("ошибка при анмаршалинге ответа от сервера: %v", err)
	}

	return &songDetails, nil
}
