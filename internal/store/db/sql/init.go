package sql

import (
	"context"

	"github.com/jackc/pgx/v5"
	"log"
)

func InitializeDatabase(conn *pgx.Conn) {
	createUsersTable := `
    CREATE TABLE IF NOT EXISTS users (
        id SERIAL PRIMARY KEY,
        name VARCHAR(100) NOT NULL,
        email VARCHAR(100) UNIQUE NOT NULL
    );
    `

	createBooksTable := `
    CREATE TABLE IF NOT EXISTS books (
        id SERIAL PRIMARY KEY,
        title VARCHAR(255) NOT NULL,
        author VARCHAR(255) NOT NULL,
        available BOOLEAN DEFAULT TRUE
    );
    `

	createUserBooksTable := `
    CREATE TABLE IF NOT EXISTS user_books (
        user_id INT REFERENCES users(id),
        book_id INT REFERENCES books(id),
        taken_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        return_date TIMESTAMP,
        PRIMARY KEY (user_id, book_id)
    );
    `

	_, err := conn.Exec(context.Background(), createUsersTable)
	if err != nil {
		log.Fatalf("Error creating users table: %v", err)
	}

	_, err = conn.Exec(context.Background(), createBooksTable)
	if err != nil {
		log.Fatalf("Error creating books table: %v", err)
	}

	_, err = conn.Exec(context.Background(), createUserBooksTable)
	if err != nil {
		log.Fatalf("Error creating user_books table: %v", err)
	}
}
