package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/robberhood/final_project/pkg/api"
)

func Start() {
	r := chi.NewRouter()

	api.Init(r)
	r.Handle("/*", http.FileServer(http.Dir("web")))

	http.ListenAndServe(":7540", r)
}
