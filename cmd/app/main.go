package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"github.com/Ablyamitov/simple-rest/internal/app"
	"github.com/Ablyamitov/simple-rest/internal/app/handlers"
	"github.com/Ablyamitov/simple-rest/internal/app/server"
	"github.com/Ablyamitov/simple-rest/internal/app/wrapper"
	"github.com/Ablyamitov/simple-rest/internal/store"
	"github.com/Ablyamitov/simple-rest/internal/store/db"
	"github.com/Ablyamitov/simple-rest/internal/store/db/repository"
	redisconn "github.com/Ablyamitov/simple-rest/internal/store/redis"
)

func main() {

	config := app.Config
	migrationPath, err := filepath.Abs(config.Migration.Path)
	if err != nil {
		log.Fatalf("Failed to get absolute path: %v", err)
	}
	store.ApplyMigrations(migrationPath, config.Migration.URL)

	conn := db.Connect(config.DB.URL)
	//if err != nil {
	//	wrapper.LogError(fmt.Sprintf("Could not connect to the database: %v", err),
	//		"main")
	//	os.Exit(1)
	//}
	//defer func(conn *pgx.Conn, ctx context.Context) {
	//	err := conn.Close(ctx)
	//	if err != nil {
	//		wrapper.LogError(fmt.Sprintf("Error closing database connection: %v", err),
	//			"main")
	//	}
	//}(conn, context.Background())

	redisClient := redisconn.Connect(config)
	//defer func(redisClient *redis.Client) {
	//	err := redisClient.Close()
	//	if err != nil {
	//		wrapper.LogError(fmt.Sprintf("Error closing redis connection: %v", err),
	//			"main")
	//	}
	//}(redisClient)

	userRepository := repository.NewUserRepository(conn, redisClient)
	userHandler := handlers.NewUserHandler(userRepository)
	authHandler := handlers.NewAuthHandler(config.App.Secret, userRepository)

	bookRepository := repository.NewBookRepository(conn, redisClient)
	bookHandler := handlers.NewBookHandler(bookRepository)

	srv := server.NewServer(userHandler, bookHandler, authHandler, config.Server.Host, config.Server.Port)

	if err := srv.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		wrapper.LogError(fmt.Sprintf("Could not listen on %s:%d: %v\n", config.Server.Host, config.Server.Port, err),
			"main")
	}

}
