package main

import (
	"fmt"

	"net/http"
	"time"

	"github.com/balebbae/sodia/docs" // This is rquired to generate swagger docs
	"github.com/balebbae/sodia/internal/mailer"
	"github.com/balebbae/sodia/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
)

type application struct {
	config config
	store store.Storage
	logger *zap.SugaredLogger
	mailer mailer.Client
}

type config struct {
	addr string
	db dbConfig
	env string
	apiURL string
	mail mailConfig
	frontendURL string
}

type mailConfig struct {
	sendGrid sendGridConfig
	fromEmail string
	exp time.Duration
}

type sendGridConfig struct {
	apiKey string
}

type dbConfig struct {
	addr string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime string
}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
  	r.Use(middleware.RealIP)
  	r.Use(middleware.Logger)
  	r.Use(middleware.Recoverer)

	r.Use(middleware.Timeout(60 * time.Second))

    // 1) Enable CORS
    r.Use(cors.Handler(cors.Options{
        // Use your frontend URL here. This might be "http://localhost:5173".
        AllowedOrigins:   []string{app.config.frontendURL},

        // If you need multiple, simply add them:
        // AllowedOrigins: []string{"http://localhost:3000", "http://localhost:5173"},

        // You can also set a wildcard: []string{"*"}

        AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
        ExposedHeaders:   []string{"Link"},
        AllowCredentials: false,
        MaxAge:           300, // 300 seconds = 5 minutes
    }))

	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", app.healthCheckHandler)
		
		docsURL := fmt.Sprintf("%s/swagger/doc.json", app.config.addr)
		r.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL(docsURL)))

		r.Route("/posts", func(r chi.Router) {
			r.Post("/", app.createPostHandler) // POST /v1/Posts
			r.Route("/{postID}", func(r chi.Router) { // WE will need postID more later
				r.Use(app.postsContextMiddleware)
				r.Get("/", app.getPostHandler)
				r.Patch("/", app.updatePostHandler)
				r.Delete("/", app.deletePostHandler)

				//Comments
				r.Post("/comments", app.createCommentHandler)
			})
		})
		r.Route("/users", func(r chi.Router) {
			r.Put("/activate/{token}", app.activateUserHandler)

			r.Route("/{userID}", func(r chi.Router) { 
				r.Use(app.userContextMiddleware)
				r.Get("/", app.getUserHandler)
				r.Put("/follow", app.followUserHandler)
				r.Put("/unfollow", app.unfollowUserHandler)
			})

			r.Group(func(r chi.Router){
				r.Get("/feed", app.getUserFeedHandler)
			})
		})

		// Public Routes
		r.Route("/authentication", func(r chi.Router) {
			r.Post("/user", app.registerUserHandler)
		})
	})

	return r
}

func (app *application) run(mux http.Handler) error {
	// Docs
	docs.SwaggerInfo.Version = version
	docs.SwaggerInfo.Host = app.config.apiURL
	docs.SwaggerInfo.BasePath = "/v1"
	
	server := &http.Server{
		Addr: app.config.addr,
		Handler: mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout: time.Second * 10,
		IdleTimeout: time.Minute,
	}

	app.logger.Infow("server has started", "addr", app.config.addr, "env", app.config.env)

	return server.ListenAndServe()
}

// create user(user struct, *db.DB){
// }