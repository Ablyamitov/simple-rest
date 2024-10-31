package db

import (
	"context"
	"fmt"
	"github.com/Ablyamitov/simple-rest/internal/app/wrapper"
	"github.com/jackc/pgx/v5"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
)

func Connect(connectionUrl string) *pgx.Conn {

	conn, err := pgx.Connect(context.Background(), connectionUrl)
	if err != nil {
		wrapper.LogError(fmt.Sprintf("Could not connect to the database: %v", err),
			"main")
		os.Exit(1)
	}

	return conn
}

func ConnectGorm(connectionUrl string) *gorm.DB {
	dsn := "host=localhost user=postgres dbname=library password=12345678 sslmode=disable"
	conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Could not connect gorm to the db: %v", err)
	}
	return conn
}
