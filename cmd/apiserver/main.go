package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Ablyamitov/simple-rest/internal/app/apiserver"
	"github.com/Ablyamitov/simple-rest/internal/app/apiserver/handlers"
	"github.com/Ablyamitov/simple-rest/internal/store/db"
	"github.com/Ablyamitov/simple-rest/internal/store/db/entity"
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

	userRepository := repository.NewUserRepository(conn)
	userHandler := handlers.NewUserHandler(userRepository)

	bookRepository := repository.NewBookRepository(conn)

	r := chi.NewRouter()
	log.Info("Starting server...")

	r.Get("/users", userHandler.GetAll)
	r.Get("/users/{id}", userHandler.GetById)
	r.Post("/users/add", userHandler.Create)
	r.Patch("/users/update", userHandler.Update)
	r.Delete("/users/{id}", userHandler.Delete)

	r.Post("/books", func(w http.ResponseWriter, r *http.Request) {
		var book entity.Book
		if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err := bookRepository.Create(context.Background(), &book)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)

		if err := json.NewEncoder(w).Encode(book); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	if err := http.ListenAndServe(fmt.Sprintf("%s:%d", config.Server.Host, config.Server.Port), r); err != nil {
		log.Error("Server failed to start: ", err)
		os.Exit(1)
	}

}

func setupLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, nil))
}
