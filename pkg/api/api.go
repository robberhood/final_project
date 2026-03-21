package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/robberhood/final_project/pkg/api/handlers"
)

func Init(r chi.Router) {
	r.Get("/api/nextdate", handlers.NextDayHandler)
}
