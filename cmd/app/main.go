package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/Ablyamitov/simple-rest/internal/app"
	"github.com/Ablyamitov/simple-rest/internal/app/handlers"
	"github.com/Ablyamitov/simple-rest/internal/app/server"
	"github.com/Ablyamitov/simple-rest/internal/app/wrapper"
	"github.com/Ablyamitov/simple-rest/internal/store"
	"github.com/Ablyamitov/simple-rest/internal/store/db"
	"github.com/Ablyamitov/simple-rest/internal/store/db/repository"
	redisconn "github.com/Ablyamitov/simple-rest/internal/store/redis"

	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
)

func main() {

	//conf
	config := app.LoadConfig()
	//migration
	store.ApplyMigrations(config.Migration.Path, config.Migration.URL)

	//postgres
	conn := db.Connect(config.DB.URL)
	defer func(conn *pgx.Conn, ctx context.Context) {
		err := conn.Close(ctx)
		if err != nil {
			wrapper.LogError(fmt.Sprintf("Сlosing database connection: %v", err),
				"main")
		}
	}(conn, context.Background())

	//redis
	redisClient := redisconn.Connect(config.Redis.Addr, config.Redis.Password, config.Redis.DB)
	defer func(redisClient *redis.Client) {
		err := redisClient.Close()
		if err != nil {
			wrapper.LogError(fmt.Sprintf("Сlosing redis connection: %v", err),
				"main")
		}
	}(redisClient)

	userRepository := repository.NewUserRepository(conn, redisClient)
	userHandler := handlers.NewUserHandler(userRepository)
	authHandler := handlers.NewAuthHandler(config.App.Secret, userRepository)

	bookRepository := repository.NewBookRepository(conn, redisClient)
	bookHandler := handlers.NewBookHandler(bookRepository)

	srv := server.NewServer(userHandler, bookHandler, authHandler, config.App.Secret)

	if err := srv.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		wrapper.LogError(fmt.Sprintf("Could not listen on %s:%d: %v\n", config.Server.Host, config.Server.Port, err),
			"main")
	}

}
