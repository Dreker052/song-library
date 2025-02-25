package models

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
	SongId      int    `swaggerignore:"true"`          //Внешний ключ
	Text        string `json:"text" example:""`        //Текст песни
	ReleaseDate string `json:"releaseDate" example:""` //Дата выхода песни
	Link        string `json:"link" example:""`        //Ссылка на песню
}

// SongWithDetails представляет собой модель песни с дополнительными данными песни.
// @Description Модель песни c дополнительными данными песни
type SongWithDetails struct {
	Group       string      `json:"group" example:""`
	Song        string      `json:"song" example:""`
	SongDetails SongDetails `json:"SongDetails"`
}
