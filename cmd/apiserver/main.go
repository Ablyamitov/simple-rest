package main

import (
	"context"
	"fmt"
	"github.com/Ablyamitov/simple-rest/internal/app/apiserver"
	"github.com/Ablyamitov/simple-rest/internal/app/apiserver/handlers"
	"github.com/Ablyamitov/simple-rest/internal/store/db"
	"github.com/Ablyamitov/simple-rest/internal/store/db/repository"
	"github.com/Ablyamitov/simple-rest/internal/store/db/sql"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	log := setupLogger()

	config := apiserver.LoadConfig()
	conn, err := db.Connect(config)
	if err != nil {
		log.Error("Could not connect to the database: %v", err)
		os.Exit(1)
	}
	defer func(conn *pgx.Conn, ctx context.Context) {
		err := conn.Close(ctx)
		if err != nil {
			log.Error("Error closing database connection: %v", err)
		}
	}(conn, context.Background())

	sql.InitializeDatabase(conn)

	r := chi.NewRouter()
	log.Info("Starting server...")

	//users
	userRepository := repository.NewUserRepository(conn)
	userHandler := handlers.NewUserHandler(userRepository)
	r.Get("/users", userHandler.GetAll)
	r.Get("/users/{id}", userHandler.GetById)
	r.Post("/users/add", userHandler.Create)
	r.Patch("/users/update", userHandler.Update)
	r.Delete("/users/{id}", userHandler.Delete)

	//books
	bookRepository := repository.NewBookRepository(conn)
	bookHandler := handlers.NewBookHandler(bookRepository)
	r.Get("/books", bookHandler.GetAll)
	r.Get("/books/{id}", bookHandler.GetById)
	r.Post("/books/add", bookHandler.Create)
	r.Patch("/books/update", bookHandler.Update)
	r.Delete("/books/{id}", bookHandler.Delete)

	if err := http.ListenAndServe(fmt.Sprintf("%s:%d", config.Server.Host, config.Server.Port), r); err != nil {
		log.Error("Server failed to start: ", err)
		os.Exit(1)
	}

}

func setupLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, nil))
}
