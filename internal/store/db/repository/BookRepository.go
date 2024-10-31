package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Ablyamitov/simple-rest/internal/store/db/entity"

	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
)

const (
	SELECT_ALL_BOOKS = `
				  SELECT * 
				  FROM books`
	SELECT_BOOK_BY_ID = `
				  SELECT id, title, author, available 
				  FROM books 
				  WHERE id=$1`
	INSERT_BOOK = `
				  INSERT INTO books (title, author) 
				  VALUES ($1, $2) 
				  RETURNING books.id,books.available`
	UPDATE_BOOK = `
				  UPDATE books 
				  SET title = $1, author = $2 
				  WHERE id = $3`
	DELETE_BOOK = `
				  DELETE 
				  FROM books 
				  WHERE id = $1`
)

type BookRepository interface {
	GetALL(ctx context.Context) ([]entity.Book, error)
	GetByID(ctx context.Context, id int) (*entity.Book, error)
	Create(ctx context.Context, book *entity.Book) error
	Update(ctx context.Context, book *entity.Book) (*entity.Book, error)
	Delete(ctx context.Context, id int) error
}

type BookRepositoryImpl struct {
	Conn        *pgx.Conn
	RedisClient *redis.Client
}

func NewBookRepository(conn *pgx.Conn, redisClient *redis.Client) BookRepository {
	return &BookRepositoryImpl{Conn: conn, RedisClient: redisClient}
}

func (bookRepository *BookRepositoryImpl) GetALL(ctx context.Context) ([]entity.Book, error) {
	rows, err := bookRepository.Conn.Query(ctx, SELECT_ALL_BOOKS)
	if err != nil {
		return nil, err
	}
	var books []entity.Book
	for rows.Next() {
		var book entity.Book
		err = rows.Scan(&book.ID, &book.Title, &book.Author, &book.Available)
		if err != nil {
			return nil, err
		}
		books = append(books, book)
	}
	return books, nil
}

func (bookRepository *BookRepositoryImpl) GetByID(ctx context.Context, id int) (*entity.Book, error) {
	cacheKey := fmt.Sprintf("book:%d", id)
	cachedUser, err := bookRepository.RedisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var book entity.Book
		err = json.Unmarshal([]byte(cachedUser), &book)
		if err == nil {
			return &book, nil
		}
	}

	book := &entity.Book{}
	err = bookRepository.Conn.QueryRow(ctx, SELECT_BOOK_BY_ID, id).Scan(&book.ID, &book.Title, &book.Author, &book.Available)
	if err != nil {
		return nil, err
	}
	//Сохранение кеша
	userData, err := json.Marshal(book)
	if err == nil {
		bookRepository.RedisClient.Set(ctx, cacheKey, userData, 0)
	}
	return book, nil

}

func (bookRepository *BookRepositoryImpl) Create(ctx context.Context, book *entity.Book) error {
	err := bookRepository.Conn.QueryRow(ctx,
		INSERT_BOOK,
		book.Title, book.Author).Scan(&book.ID, &book.Available)
	return err
}

func (bookRepository *BookRepositoryImpl) Update(ctx context.Context, book *entity.Book) (*entity.Book, error) {

	_, err := bookRepository.Conn.Exec(ctx,
		UPDATE_BOOK, book.Title, book.Author, book.ID)
	if err != nil {
		return nil, err
	}

	err = bookRepository.Conn.QueryRow(ctx,
		SELECT_BOOK_BY_ID, book.ID).
		Scan(&book.ID, &book.Title, &book.Author, &book.Available)

	if err != nil {
		return nil, err
	}
	// Удаление книги с кеша
	bookCacheKey := fmt.Sprintf("book:%d", book.ID)
	err = bookRepository.RedisClient.Del(ctx, bookCacheKey).Err()
	if err != nil {
		return nil, err
	}

	return book, nil
}

func (bookRepository *BookRepositoryImpl) Delete(ctx context.Context, id int) error {
	_, err := bookRepository.Conn.Exec(ctx, DELETE_BOOK, id)
	// Удаление книги с кеша
	bookCacheKey := fmt.Sprintf("book:%d", id)
	err = bookRepository.RedisClient.Del(ctx, bookCacheKey).Err()
	if err != nil {
		return err
	}
	return err
}
