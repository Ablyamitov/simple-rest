package dto

type UserDTO struct {
	ID       int        `json:"id"`
	Name     string     `json:"name" validate:"required,notblank"`
	Email    string     `json:"email" validate:"email,required,notblank"`
	Books    []*BookDTO `json:"books"`
	Password string     `json:"password" validate:"required,notblank"`
	Role     string     `json:"role"`
}
