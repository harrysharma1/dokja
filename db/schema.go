package db

import "gorm.io/gorm"

type WebNovel struct {
	gorm.Model
	Name          string
	AuthorName    string
	TotalChapters int
	Info          string
	UrlPath       string
}
