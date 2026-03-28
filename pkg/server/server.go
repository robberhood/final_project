package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/robberhood/final_project/pkg/api"
)

func Start(port string) {
	r := chi.NewRouter()

	api.Init(r)
	r.Handle("/*", http.FileServer(http.Dir("web")))

	http.ListenAndServe(port, r)
}
