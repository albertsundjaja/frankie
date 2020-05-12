package main

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"net/http"
	"github.com/albertsundjaja/frankie/isGood"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Post("/isgood", isGood.IsGoodHandler)
	http.ListenAndServe(":8888", r)
}



