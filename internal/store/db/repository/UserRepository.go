package repository

import (
	"context"
	"github.com/Ablyamitov/simple-rest/internal/store/db/entity"
	"github.com/jackc/pgx/v5"
)

type UserRepository struct {
	Conn *pgx.Conn
}

func NewUserRepository(conn *pgx.Conn) *UserRepository {
	return &UserRepository{Conn: conn}
}
func (userRepository *UserRepository) GetAll(ctx context.Context) ([]entity.User, error) {
	rows, err := userRepository.Conn.Query(ctx, "SELECT id, name, email FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []entity.User
	for rows.Next() {
		var user entity.User
		err = rows.Scan(&user.ID, &user.Name, &user.Email)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil

}

func (userRepository *UserRepository) GetByID(ctx context.Context, id int) (*entity.User, error) {
	user := &entity.User{}
	err := userRepository.Conn.QueryRow(ctx, "SELECT id, name, email FROM users WHERE id=$1", id).Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (userRepository *UserRepository) Create(ctx context.Context, user *entity.User) error {
	//_, err := userRepository.Conn.Exec(ctx, "INSERT INTO users (name, email) VALUES ($1, $2)", user.Name, user.Email)
	err := userRepository.Conn.QueryRow(ctx,
		"INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id",
		user.Name, user.Email).Scan(&user.ID)
	return err
}

func (userRepository *UserRepository) Update(ctx context.Context, user *entity.User) (*entity.User, error) {
	_, err := userRepository.Conn.Exec(ctx, "UPDATE users set name = $1, email=$2 WHERE id = $3", user.Name, user.Email, user.ID)
	if err != nil {
		return nil, err
	}
	err = userRepository.Conn.QueryRow(ctx, "SELECT id, name, email FROM users WHERE id = $1", user.ID).Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (userRepository *UserRepository) Delete(ctx context.Context, id int) error {
	_, err := userRepository.Conn.Exec(ctx, "DELETE FROM users WHERE id=$1", id)
	return err
}
