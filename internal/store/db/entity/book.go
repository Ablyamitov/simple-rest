package entity

import "gorm.io/gorm"

type Book struct {
	gorm.Model
	ID        int     `json:"id"`
	Title     string  `json:"title" validate:"required,notblank"`
	Author    string  `json:"author" validate:"required,notblank"`
	Available bool    `json:"available"`
	Users     []*User `gorm:"many2many:user_books;" json:"users"`
}
