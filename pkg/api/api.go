package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/robberhood/final_project/pkg/api/handlers"
	"github.com/robberhood/final_project/pkg/api/service"
)

func Init(r chi.Router) {
	r.Post("/api/signin", handlers.SignInHandler)
	r.Get("/api/nextdate", handlers.NextDayHandler)

	r.Group(func(r chi.Router) {
		r.Use(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				service.Auth(func(w http.ResponseWriter, r *http.Request) {
					next.ServeHTTP(w, r)
				})(w, r)
			})
		})

		r.Post("/api/task", handlers.TaskHandler)
		r.Put("/api/task", handlers.TaskUPDHandler)
		r.Delete("/api/task", handlers.TaskDELHandler)
		r.Get("/api/tasks", handlers.TasksHandler)
		r.Post("/api/task/done", handlers.TaskCompleteHandler)
	})
}
