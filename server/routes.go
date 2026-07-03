package server

import (
	"Task-Management/handler"
	"Task-Management/middleware"
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

type Server struct {
	chi.Router
	server *http.Server
}

const (
	readTimeout       = 5 * time.Minute
	readHeaderTimeout = 30 * time.Second
	writeTimeout      = 5 * time.Minute
)

func SetupRoutes() *Server {
	r := chi.NewRouter()

	r.Route("/v1", func(r chi.Router) {
		r.Post("/register", handler.RegisterUser)
		r.Post("/login", handler.LoginUser)

		r.Route("/users", func(r chi.Router) {
			r.Use(middleware.Authenticate)

			r.Get("/me", handler.GetUser)
			r.Post("/logout", handler.LogoutUser)
			r.Delete("/", handler.DeleteUser)
		})

		r.Route("/todos", func(r chi.Router) {
			r.Use(middleware.Authenticate)

			r.Post("/", handler.CreateTodo)
			r.Get("/", handler.GetTodos)
			r.Delete("/", handler.DeleteAllTodos)

			r.Route("/{todoId}", func(r chi.Router) {
				r.Get("/", handler.GetTodoById)
				r.Put("/", handler.UpdateTodo)
				r.Put("/marktodoascompleted", handler.MarkTodoAsCompleted)
				r.Delete("/", handler.DeleteTodo)
			})
		})
	})

	return &Server{
		Router: r,
	}
}

func (svc *Server) Run(port string) error {
	svc.server = &http.Server{
		Addr:              port,
		Handler:           svc.Router,
		ReadTimeout:       readTimeout,
		ReadHeaderTimeout: readHeaderTimeout,
		WriteTimeout:      writeTimeout,
	}
	return svc.server.ListenAndServe()
}

func (svc *Server) Shutdown(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return svc.server.Shutdown(ctx)
}
