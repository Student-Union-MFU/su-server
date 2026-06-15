package main

import (
	"fmt"
	"log"
	"net/http"
	"su-server/config"
	"su-server/internal/handler"
	"su-server/internal/repository"
	"su-server/internal/service"

	"github.com/go-chi/cors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
    r := chi.NewRouter()
    
	r.Use(middleware.Logger)
    r.Use(middleware.Recoverer)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://localhost:3001", "https://yourdomain.com"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	db, err := config.ConnectDB()
	if err != nil {
    	log.Fatal("DB connection failed:", err)
}

	eventRepository := repository.NewEventRepository(db)
	eventService := service.NewEventService(eventRepository)
	eventHandler := handler.NewEventHandler(eventService)

    r.Get("/", func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "SU Backend running"}`))
    })

	r.Route("/su-server", func(r chi.Router) {
		r.Route("/events", func(r chi.Router) {
			r.Get("/", eventHandler.GetAllEvents)
			r.Get("/{id}", eventHandler.GetOneEvents)
			r.Post("/", eventHandler.CreateOneEvent)
			// r.Put("/{id}", eventHandler.)
			r.Delete("/{id}", eventHandler.DeleteOneEvents)
		})
	})

    fmt.Println("Server running on :8080")
    if err := http.ListenAndServe(":8080", r); err != nil {
        log.Fatal(err)
    }
    http.ListenAndServe(":8080", r)
}
