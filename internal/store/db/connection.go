package db

import (
	"context"
	"fmt"
	"github.com/Ablyamitov/simple-rest/internal/app/wrapper"
	"github.com/jackc/pgx/v5"
	"os"
	"sync"
)

var (
	conn     *pgx.Conn
	connOnce sync.Once
)

func Connect(connectionUrl string) *pgx.Conn {
	var err error
	connOnce.Do(func() {
		conn, err = pgx.Connect(context.Background(), connectionUrl)
		if err != nil {
			wrapper.LogError(fmt.Sprintf("Could not connect to the database: %v", err),
				"main")
			os.Exit(1)
		}
		defer func(conn *pgx.Conn, ctx context.Context) {
			err := conn.Close(ctx)
			if err != nil {
				wrapper.LogError(fmt.Sprintf("Error closing database connection: %v", err),
					"main")
			}
		}(conn, context.Background())
	})

	return conn
}
