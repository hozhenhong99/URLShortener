package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	"urlShortener/cache"
	"urlShortener/handler"
)

func main() {
	fmt.Println("Starting server...")
	cache.InitCache()
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Get("/", handler.DefaultHandler)
	router.Get("/{url}", handler.RedirectHandler)
	router.Post("/shorten", handler.ShortenHandler)

	err := http.ListenAndServe(":8080", router)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
	fmt.Println("ending")
}
