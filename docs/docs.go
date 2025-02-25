// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/songs": {
            "get": {
                "description": "Получить список всех песен",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "songs"
                ],
                "summary": "Получить все песни",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Фильтр по группе",
                        "name": "group",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Фильтр по названию песни",
                        "name": "song",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Фильтр по ссылке",
                        "name": "link",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Фильтр по тексту или фрагменту тектса",
                        "name": "text",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "Поле для сортировки asc для возрастания и desc для убывания",
                        "name": "sort",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "default": 1,
                        "description": "Номер страницы(пагинация)",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "default": 10,
                        "description": "Лимит записей на странице",
                        "name": "limit",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.Song"
                            }
                        }
                    },
                    "400": {
                        "description": "Неверный формат параметров запроса",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Внутренняя ошибка сервера",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            },
            "post": {
                "description": "Добавить новую песню",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "songs"
                ],
                "summary": "Добавить песню",
                "parameters": [
                    {
                        "description": "Данные песни",
                        "name": "song",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.Song"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/models.Song"
                        }
                    },
                    "400": {
                        "description": "Песня уже добавлена",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Ошибка при получении данных от внешнего API",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/songs/{id}": {
            "put": {
                "description": "Редактировать песню по её ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "songs"
                ],
                "summary": "Редактировать песню",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "ID песни",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Данные песни",
                        "name": "song",
                        "in": "body",
                        "schema": {
                            "$ref": "#/definitions/models.SongWithDetails"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.Song"
                        }
                    },
                    "400": {
                        "description": "Неверный формат даты. Ожидаемый формат: DD.MM.YYYY",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Детали песни не найдены",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Внутренняя ошибка сервера",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            },
            "delete": {
                "description": "Удалить песню по её ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "songs"
                ],
                "summary": "Удалить песню",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "ID песни",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Песня успешно удалена",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Песня не найдена",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Ошибка при удалении песни",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/songs/{id}/text": {
            "get": {
                "description": "Получить текст песни по её ID",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "songs"
                ],
                "summary": "Получить текст песни",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "ID песни",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "default": 1,
                        "description": "Номер страницы(пагинация)",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "default": 5,
                        "description": "Лимит куплетов на странице",
                        "name": "limit",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Текст песни",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Неверный формат параметров запроса",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "404": {
                        "description": "Текст песни отсутствует",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Внутренняя ошибка сервера",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "models.Song": {
            "description": "Модель песни",
            "type": "object",
            "properties": {
                "group": {
                    "description": "Название группы",
                    "type": "string"
                },
                "song": {
                    "description": "Название песни",
                    "type": "string"
                }
            }
        },
        "models.SongDetails": {
            "description": "Модель дополнительных данных песни",
            "type": "object",
            "properties": {
                "link": {
                    "description": "Ссылка на песню",
                    "type": "string",
                    "example": ""
                },
                "releaseDate": {
                    "description": "Дата выхода песни",
                    "type": "string",
                    "example": ""
                },
                "text": {
                    "description": "Текст песни",
                    "type": "string",
                    "example": ""
                }
            }
        },
        "models.SongWithDetails": {
            "description": "Модель песни c дополнительными данными песни",
            "type": "object",
            "properties": {
                "SongDetails": {
                    "$ref": "#/definitions/models.SongDetails"
                },
                "group": {
                    "type": "string",
                    "example": ""
                },
                "song": {
                    "type": "string",
                    "example": ""
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "",
	Description:      "API для управления библиотекой песен.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
