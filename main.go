// Package main implements a client for Person service.
package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"

	"github.com/finest08/PubSubPublisher/handler/v1"
)

type Person struct {
	FirstName  string
	LastName   string
	Email      string
	Occupation string
	Age        string
}

func main() {
	r := chi.NewRouter()
	r.Use(
		middleware.Logger,
		middleware.StripSlashes,
		cors.Handler(cors.Options{
			AllowedOrigins:   []string{"https://*", "http://*"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "QUERY"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: true,
			MaxAge:           300,
			// Debug:            true,
		}),
	)

	p := &handler.PersonServer{}
	r.Route("/person", func(r chi.Router) {
		r.Post("/", p.Person)
	})

	// start server
	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Print(":" + os.Getenv("PORT"))
		fmt.Print(err)
	}
}


