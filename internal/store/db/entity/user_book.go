package entity

import "github.com/jackc/pgx/v5/pgtype"

type UserBook struct {
	UserId     int              `json:"user_id"`
	BookId     int              `json:"book_id"`
	TakenData  pgtype.Timestamp `json:"taken_data"`
	ReturnData pgtype.Timestamp `json:"return_data"`
}
