package server

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/serhappy/rssagg/internal/app"
	"github.com/serhappy/rssagg/internal/handler"
	"github.com/serhappy/rssagg/internal/middleware"
)

func NewServer(app *app.App) http.Handler {
	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1Router := chi.NewRouter()

	v1Router.Get("/healthz", handler.Healthz)
	v1Router.Get("/error", handler.Error)

	userHandler := handler.NewUserHandler(app)
	v1Router.Route("/users", func(r chi.Router) {
		r.Post("/", userHandler.CreateUser)
		r.With(middleware.Auth(app)).Get("/", userHandler.GetUser)
	})

	feedHandler := handler.NewFeedHandler(app)
	v1Router.Route("/feeds", func(r chi.Router) {
		r.With(middleware.Auth(app)).Post("/", feedHandler.CreateFeed)
		r.Get("/", feedHandler.GetAllFeeds)
	})

	feedFollowHandler := handler.NewFeedFollowHandler(app)
	v1Router.Route("/feed_follows", func(r chi.Router) {
		r.With(middleware.Auth(app)).Post("/", feedFollowHandler.CreateFeedFollow)
		r.With(middleware.Auth(app)).Get("/", feedFollowHandler.GetUserFeedFollows)
		r.With(middleware.Auth(app)).Delete("/{feedFollowID}", feedFollowHandler.DeleteFeedFollow)
	})

	router.Mount("/v1", v1Router)

	return router
}
