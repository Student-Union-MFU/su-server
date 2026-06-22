package main

import (
	"log/slog"
	"net/http"
	"os"
	"su-server/config"
	"su-server/internal/handler"
	"su-server/internal/repository"
	"su-server/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

func main() {

	godotenv.Load()

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
    	slog.Error("DB connection failed:", "err", err)
	} else {
		slog.Info("DB CONNECTED")
	}	

	eventRepository := repository.NewEventRepository(db)
	eventService := service.NewEventService(eventRepository)
	eventHandler := handler.NewEventHandler(eventService)

	userRepository := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepository)
	userHandler := handler.NewUserHandler(userService)

	jwtService := service.NewJWTService()

	oauthService := service.NewOAuthService(userService)
	oauthHandler := handler.NewOAuthHandler(oauthService, jwtService)

	stepRepository := repository.NewStepsRepository(db)
	stepService := service.NewStepsService(stepRepository)
	stepHandler := handler.NewStepsHandler(stepService)

	leaderboardRepository := repository.NewLeaderboardRepository(db)
	leaderboardService := service.NewLeaderboardService(leaderboardRepository)
	leaderboardHandler := handler.NewLeaderboardHandler(leaderboardService)

    r.Get("/", func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "SU Backend running"}`))
    })

	r.Route("/su-server", func(r chi.Router) {
		r.Route("/events", func(r chi.Router) {
			r.Get("/", eventHandler.GetAllEvents)
			r.Get("/{id}", eventHandler.GetOneEvents)
			r.Post("/", eventHandler.CreateOneEvent)
			r.Put("/{id}", eventHandler.UpdateOneEvent)
			r.Delete("/{id}", eventHandler.DeleteOneEvents)
		})
		r.Route("/users", func(r chi.Router) {
			r.Get("/{id}", userHandler.GetUserByID)
			r.Get("/email/{email}", userHandler.GetUserByEmail)
			r.Post("/insert", userHandler.InsertUser)
			r.Post("/upsert", userHandler.UpsertUser)
			r.Patch("/{id}", userHandler.UpdateUser)
			r.Delete("/{id}", eventHandler.DeleteOneEvents)
		})
		r.Route("/auth", func(r chi.Router) {
			r.Get("/google", oauthHandler.GoogleLogin)
			r.Get("/google/callback", oauthHandler.GoogleCallback)
			r.Post("/google/verify", oauthHandler.GoogleVerify)
		})
		r.Route("/steps", func(r chi.Router) {
			r.Get("/{userID}", stepHandler.GetStepsByUserID)
			r.Get("/{userID}/range?from=2026-06-01&to=2026-08-01", stepHandler.GetStepsByDateRange)
			r.Post("/sync", stepHandler.SyncSteps)
			r.Post("/sync/bulk", stepHandler.SyncManySteps)
		})
		r.Route("/leaderboard", func(r chi.Router) {
			r.Get("/", leaderboardHandler.GetLeaderboard)
			r.Get("/{userID}", leaderboardHandler.GetUserRank)
			r.Post("/update", leaderboardHandler.UpdateEntry)
			r.Post("/reset", leaderboardHandler.Reset)
		})
	})

	serverPort := os.Getenv("PORT")

	if serverPort == "" { serverPort = "8080"}
    
	slog.Info("Server running :", "port", serverPort )
	
	if err := http.ListenAndServe(":" + serverPort, r); err != nil {
        slog.Error("SERVER RUN FAILED", "err", err)
    }
    
	http.ListenAndServe(":8080", r)
}
