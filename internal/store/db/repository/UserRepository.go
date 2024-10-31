package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"

	"github.com/Ablyamitov/simple-rest/internal/store/db/entity"

	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
)

const (
	SELECT_ALL_USERS = `
				  SELECT id, name, email, password, role
				  FROM users`

	SELECT_ALL_USERS_BOOKS = `
			 	  SELECT ub.user_id, ub.book_id, b.title, b.author, b.available 
			 	  FROM books AS b 
				  JOIN user_books AS ub ON b.id = ub.book_id`

	SELECT_USER_BY_ID = `
				  SELECT id, name, email, password, role
				  FROM users 
				  WHERE id=$1`

	SELECT_USER_BY_EMAIL = `
				  SELECT id, name, email, password, role
				  FROM users 
				  WHERE email=$1`

	SELECT_ALL_USER_BOOKS_BY_ID = `
				  SELECT b.id, b.title, b.author, b.available 
				  FROM books AS b 
				      JOIN user_books AS ub 
				          ON b.id = ub.book_id 
				  WHERE ub.user_id = $1`

	INSERT_USER = `
				  INSERT INTO users (name, email, password, role) 
				  VALUES ($1, $2, $3, $4) 
				  RETURNING id`

	UPDATE_USER = `
				  UPDATE users 
				  SET name = $1, email=$2
				  WHERE id = $3`

	DELETE_USER = `
				  DELETE 
				  FROM users 
				  WHERE id=$1`
)

type UserRepository interface {
	GetAll(ctx context.Context) ([]entity.User, error)
	GetByID(ctx context.Context, id int) (*entity.User, error)
	Create(ctx context.Context, user *entity.User) error
	Update(ctx context.Context, user *entity.User) (*entity.User, error)
	Delete(ctx context.Context, id int) error
	TakeBook(ctx context.Context, userId int, bookId int) error
	ReturnBook(ctx context.Context, userId int, bookId int) error
	GetByEmail(ctx context.Context, email string) (entity.User, error)
}

type UserRepositoryImpl struct {
	Conn        *pgx.Conn
	RedisClient *redis.Client
}

func NewUserRepository(conn *pgx.Conn, redisClient *redis.Client) UserRepository {
	return &UserRepositoryImpl{Conn: conn, RedisClient: redisClient}
}

func (userRepository *UserRepositoryImpl) GetAll(ctx context.Context) ([]entity.User, error) {

	rows, err := userRepository.Conn.Query(ctx, SELECT_ALL_USERS)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []entity.User
	for rows.Next() {
		var user entity.User
		err = rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Role)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	//TODO: Идея сделать через горутину
	rows, err = userRepository.Conn.Query(ctx, SELECT_ALL_USERS_BOOKS)
	if err != nil {
		return nil, err

	}
	defer rows.Close()

	for rows.Next() {
		var userID, bookID int
		var bookTitle, bookAuthor string
		var bookAvailable bool

		err := rows.Scan(&userID, &bookID, &bookTitle, &bookAuthor, &bookAvailable)
		if err != nil {
			return nil, err
		}

		users[userID-1].Books = append(users[userID-1].Books, &entity.Book{ID: bookID, Title: bookTitle, Author: bookAuthor, Available: bookAvailable})
	}
	return users, nil

}

func (userRepository *UserRepositoryImpl) GetByID(ctx context.Context, id int) (*entity.User, error) {
	// Проверка кеша
	cacheKey := fmt.Sprintf("user:%d", id)
	cachedUser, err := userRepository.RedisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var user entity.User
		err = json.Unmarshal([]byte(cachedUser), &user)
		if err == nil {
			return &user, nil
		}
	}

	user := &entity.User{}

	err = userRepository.Conn.QueryRow(ctx, SELECT_USER_BY_ID, id).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Role)
	if err != nil {
		return nil, err
	}

	//TODO: Идея сделать через горутину
	rows, err := userRepository.Conn.Query(ctx, SELECT_ALL_USER_BOOKS_BY_ID, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		book := &entity.Book{}
		err = rows.Scan(&book.ID, &book.Title, &book.Author, &book.Available)
		if err != nil {
			return nil, err
		}
		user.Books = append(user.Books, book)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	//Сохранение кеша
	userData, err := json.Marshal(user)
	if err == nil {
		userRepository.RedisClient.Set(ctx, cacheKey, userData, 0)
	}
	return user, nil
}

func (userRepository *UserRepositoryImpl) GetByEmail(ctx context.Context, email string) (entity.User, error) {
	user := entity.User{}

	err := userRepository.Conn.QueryRow(ctx, SELECT_USER_BY_EMAIL, email).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Role)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return user, err
	}
	return user, nil
}

func (userRepository *UserRepositoryImpl) Create(ctx context.Context, user *entity.User) error {
	err := userRepository.Conn.QueryRow(ctx, INSERT_USER, user.Name, user.Email, user.Password, user.Role).Scan(&user.ID)
	return err
}

func (userRepository *UserRepositoryImpl) Update(ctx context.Context, user *entity.User) (*entity.User, error) {
	_, err := userRepository.Conn.Exec(ctx, UPDATE_USER, user.Name, user.Email, user.ID)
	if err != nil {
		return nil, err
	}
	err = userRepository.Conn.QueryRow(ctx, SELECT_USER_BY_ID, user.ID).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Role)
	if err != nil {
		return nil, err
	}
	// Удаляем данные из кеша
	cacheKey := fmt.Sprintf("user:%d", user.ID)
	err = userRepository.RedisClient.Del(ctx, cacheKey).Err()
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (userRepository *UserRepositoryImpl) Delete(ctx context.Context, id int) error {
	_, err := userRepository.Conn.Exec(ctx, DELETE_USER, id)
	// Удаление книги с кеша
	bookCacheKey := fmt.Sprintf("book:%d", id)
	err = userRepository.RedisClient.Del(ctx, bookCacheKey).Err()
	if err != nil {
		return err
	}
	return err
}

func (userRepository *UserRepositoryImpl) TakeBook(ctx context.Context, userId int, bookId int) error {

	tx, err := userRepository.Conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
				slog.Error(fmt.Sprintf("tx.Rollback failed: %v", rollbackErr))
			}
		}
	}()

	_, err = tx.Exec(ctx, "INSERT INTO user_books (user_id, book_id) VALUES ($1, $2)", userId, bookId)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, "UPDATE books SET available = FALSE WHERE id = $1", bookId)
	if err != nil {
		return err
	}

	if err = tx.Commit(ctx); err != nil {
		return err
	}
	// Удаляем данные из кеша
	cacheKey := fmt.Sprintf("user:%d", userId)
	err = userRepository.RedisClient.Del(ctx, cacheKey).Err()
	if err != nil {
		return err
	}

	return nil

}

func (userRepository *UserRepositoryImpl) ReturnBook(ctx context.Context, userId int, bookId int) error {
	tx, err := userRepository.Conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
				slog.Error(fmt.Sprintf("tx.Rollback failed: %v", rollbackErr))
			}
		}
	}()

	_, err = tx.Exec(ctx, "DELETE FROM user_books WHERE user_id = $1 AND book_id = $2", userId, bookId)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, "UPDATE books SET available = TRUE WHERE id = $1", bookId)
	if err != nil {
		return err
	}

	//_, err = tx.Exec(ctx, "UPDATE user_books SET return_date = NOW() WHERE user_id = $1 AND book_id = $2", userId, bookId)
	//if err != nil {
	//	return err
	//}

	if err = tx.Commit(ctx); err != nil {
		return err
	}
	// Удаляем данные из кеша
	cacheKey := fmt.Sprintf("user:%d", userId)
	err = userRepository.RedisClient.Del(ctx, cacheKey).Err()
	if err != nil {
		return err
	}
	return nil
}
