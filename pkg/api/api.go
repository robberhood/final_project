package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/robberhood/final_project/pkg/api/handlers"
)

func Init(r chi.Router) {
	r.Get("/api/nextdate", handlers.NextDayHandler)
	r.Post("/api/task", handlers.TaskHandler)
	r.Get("/api/tasks", handlers.TasksHandler)
	r.Get("/api/task", handlers.TaskGetHandler)
	r.Put("/api/task", handlers.TaskUPDHandler)
	r.Post("/api/task/done", handlers.TaskCompleteHandler)
	r.Delete("/api/task", handlers.TaskDELHandler)
	r.Post("/api/signin", handlers.SignInHandler)
}
