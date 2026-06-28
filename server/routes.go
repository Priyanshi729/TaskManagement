package server

import (
	"Task-Management/handler"
	"Task-Management/middleware"

	"github.com/go-chi/chi/v5"
)

func SetupRoutes() *chi.Mux {
	r := chi.NewRouter()

	r.Route("/v1", func(r chi.Router) {
		r.Post("/register", handler.RegisterUser)
		r.Post("/login", handler.LoginUser)

		r.Route("/users", func(r chi.Router) {
			r.Use(middleware.Authenticate)

			r.Get("/", handler.GetUser)
			r.Post("/logout", handler.LogoutUser)
			r.Delete("/", handler.DeleteUser)
		})
		r.Route("/todos", func(r chi.Router) {
			r.Post("/", handler.CreateTodo)
			r.Get("/", handler.GetTodos)

			r.Route("/{todoId}", func(r chi.Router) {
				r.Get("/", handler.GetTodoById)
				r.Put("/", handler.UpdateTodo)
				r.Put("/marktodoascompleted", handler.MarkTodoAsCompleted)
				r.Delete("/", handler.DeleteTodo)
			})
		})
	})

	return r
}
