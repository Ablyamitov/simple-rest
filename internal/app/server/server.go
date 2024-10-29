package server

import (
	"context"
	"log"
	"net/http"

	"github.com/Ablyamitov/simple-rest/internal/app/handlers"
	"github.com/Ablyamitov/simple-rest/internal/app/middlewares"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Server interface {
	Start() error
	Stop(ctx context.Context) error
}

type HttpServer struct {
	server *http.Server
	router *chi.Mux
}

func NewServer(userHandler *handlers.UserHandler, bookHandler *handlers.BookHandler, authHandler *handlers.AuthHandler, host string, port int) Server {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	routeUsers(r, userHandler)
	routeBooks(r, bookHandler)
	routeAuth(r, authHandler)

	r.Get("/swagger/*", httpSwagger.WrapHandler)
	r.Get("/swagger/doc.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./gen/api/openapi.yaml")
	})
	srv := &http.Server{
		//Addr:    fmt.Sprintf("%s:%d", host, port),
		Addr:    ":8080",
		Handler: r,
	}

	return &HttpServer{server: srv, router: r}
}

func routeUsers(r chi.Router, userHandler *handlers.UserHandler) {
	//users
	r.Route("/users", func(r chi.Router) {
		r.Use(middlewares.IsAuthorized)

		r.Get("/", userHandler.GetAll)            //Get All Users
		r.Get("/{id}", userHandler.GetById)       //Get User by id
		r.Post("/add", userHandler.Create)        //Create User
		r.Patch("/update", userHandler.Update)    //Update User
		r.Delete("/{id}", userHandler.Delete)     //Delete User
		r.Post("/take", userHandler.TakeBook)     //Take book to User
		r.Post("/return", userHandler.ReturnBook) //Return book from User
	})
}

func routeBooks(r chi.Router, bookHandler *handlers.BookHandler) {
	//books
	r.Route("/books", func(r chi.Router) {
		r.Use(middlewares.IsAuthorized)

		r.Get("/", bookHandler.GetAll)         //Get All Books
		r.Get("/{id}", bookHandler.GetById)    //Get Book by id
		r.Post("/add", bookHandler.Create)     //Create Book
		r.Patch("/update", bookHandler.Update) //Update Book
		r.Delete("/{id}", bookHandler.Delete)  //Delete Book
	})

}

func routeAuth(r chi.Router, authHandler *handlers.AuthHandler) {
	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", authHandler.Register)    //User register
		r.Post("/login", authHandler.Login)          //User login
		r.Post("/check-auth", authHandler.CheckAuth) //check auth
	})
}

func (s *HttpServer) Start() error {
	log.Printf("Starting server on %s", s.server.Addr)
	return s.server.ListenAndServe()
}

func (s *HttpServer) Stop(ctx context.Context) error {
	log.Printf("Stopping server on %s", s.server.Addr)
	return s.server.Shutdown(ctx)
}
