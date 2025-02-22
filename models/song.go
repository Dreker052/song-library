package models

import "gorm.io/gorm"

// Song представляет собой модель песни.
// @Description Модель песни
type Song struct {
	gorm.Model  `swaggerignore:"true"`
	Group       string `json:"group"`       //Название группы
	Song        string `json:"song"`        //Название песни
	Text        string `json:"text"`        //Текст песни
	ReleaseDate string `json:"releaseDate"` //Дата выхода песни
	Link        string `json:"link"`        //Ссылка на песню
}
