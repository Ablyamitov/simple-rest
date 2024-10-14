package db

import (
	"context"
	"github.com/Ablyamitov/simple-rest/internal/app/apiserver"
	"github.com/jackc/pgx/v5"
)

func Connect(config *apiserver.Config) (*pgx.Conn, error) {
	conn, err := pgx.Connect(context.Background(), config.DB.URL)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
