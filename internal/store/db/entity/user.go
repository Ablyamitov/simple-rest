package entity

import "gorm.io/gorm"

type User struct {
	gorm.Model
	ID       int     `json:"id" gorm:"primaryKey"`
	Name     string  `json:"name" validate:"required,notblank" gorm:""`
	Email    string  `json:"email" validate:"email,required,notblank" gorm:"unique"`
	Books    []*Book `json:"books" gorm:"many2many:user_books;"`
	Password string  `json:"password" validate:"required,notblank" gorm:"unique"`
	Role     string  `json:"role" gorm:"unique"`
}
