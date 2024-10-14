package repository

import (
	"context"
	"github.com/Ablyamitov/simple-rest/internal/store/db/entity"
	"github.com/jackc/pgx/v5"
)

type BookRepository struct {
	Conn *pgx.Conn
}

func NewBookRepository(conn *pgx.Conn) *BookRepository {
	return &BookRepository{Conn: conn}
}

func (r *BookRepository) GetALL(ctx context.Context) ([]entity.Book, error) {
	rows, err := r.Conn.Query(ctx, "SELECT * FROM books")
	if err != nil {
		return nil, err
	}
	var books []entity.Book
	for rows.Next() {
		var book entity.Book
		err = rows.Scan(&book.ID, &book.Title, &book.Author)
		if err != nil {
			return nil, err
		}
		books = append(books, book)
	}
	return books, nil
}

func (r *BookRepository) Create(ctx context.Context, book *entity.Book) error {
	_, err := r.Conn.Exec(ctx, "INSERT INTO books (title, author) VALUES ($1, $2)", book.Title, book.Author)
	return err
}

func (r *BookRepository) GetByID(ctx context.Context, id int) (*entity.Book, error) {
	book := &entity.Book{}
	err := r.Conn.QueryRow(ctx, "SELECT id, title, author, available FROM books WHERE id=$1", id).Scan(&book.ID, &book.Title, &book.Author, &book.Available)
	if err != nil {
		return nil, err
	}
	return book, nil

}
