package db

import (
	"context"
	"fmt"
	"github.com/Ablyamitov/simple-rest/internal/app/wrapper"
	"github.com/jackc/pgx/v5"
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
