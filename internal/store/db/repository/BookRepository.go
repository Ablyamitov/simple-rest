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

func (r *BookRepository) GetByID(ctx context.Context, id int) (*entity.Book, error) {
	book := &entity.Book{}
	err := r.Conn.QueryRow(ctx, "SELECT id, title, author, available FROM books WHERE id=$1", id).Scan(&book.ID, &book.Title, &book.Author, &book.Available)
	if err != nil {
		return nil, err
	}
	return book, nil

}

func (r *BookRepository) Create(ctx context.Context, book *entity.Book) error {
	err := r.Conn.QueryRow(ctx,
		"INSERT INTO books (title, author) VALUES ($1, $2) RETURNING books.id,books.available",
		book.Title, book.Author).Scan(&book.ID, &book.Available)
	return err
}

func (r *BookRepository) Update(ctx context.Context, book *entity.Book) (*entity.Book, error) {

	_, err := r.Conn.Exec(ctx,
		"UPDATE books SET title = $1, author = $2 WHERE id = $3", book.Title, book.Author, book.ID)
	if err != nil {
		return nil, err
	}

	err = r.Conn.QueryRow(ctx,
		"SELECT id, title, author, available FROM books WHERE id = $1", book.ID).
		Scan(&book.ID, &book.Title, &book.Author, &book.Available)

	if err != nil {
		return nil, err
	}
	return book, nil
}

func (r *BookRepository) Delete(ctx context.Context, id int) error {
	_, err := r.Conn.Exec(ctx, "DELETE FROM books WHERE id = $1", id)
	return err
}
