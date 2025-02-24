package models

import (
	"time"
)

// Song представляет собой модель песни.
// @Description Модель песни
type Song struct {
	Id          int         `swaggerignore:"true" gorm:"primaryKey"`
	Group       string      `json:"group"`                                                    //Название группы
	Song        string      `json:"song"`                                                     //Название песни
	SongDetails SongDetails `json:"SongDetail" swaggerignore:"true" gorm:"foreignKey:SongId"` //связь один к одному
}

// SongDetails представляет собой модель дополнительных данных песни.
// @Description Модель дополнительных данных песни
type SongDetails struct {
	SongId      int       `swaggerignore:"true"` //Внешний ключ
	Text        string    `json:"text"`          //Текст песни
	ReleaseDate time.Time `json:"releaseDate"`   //Дата выхода песни
	Link        string    `json:"link"`          //Ссылка на песню
}
