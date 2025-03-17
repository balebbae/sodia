// export PATH=$PATH:$(go env GOPATH)/bin

package main

import (
	"time"

	"github.com/balebbae/sodia/internal/db"
	"github.com/balebbae/sodia/internal/env"
	"github.com/balebbae/sodia/internal/store"
	"go.uber.org/zap"
)

const version = "0.0.2"

//	@title			Sodia API
//	@description	API for Social Media app Sodia.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@BasePath					/v1
//
//	@securityDefinitions.apiKey	ApiKeyAuth
//	@in							header
//	@name						Authorization
//	@description
func main() {
	cfg := config{
		addr: env.GetString("ADDR", ":8080"),
		apiURL: env.GetString("EXTERNAL_URL", "localhost:8080"),
		db: dbConfig{
			addr: env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost:5432/socialnetwork?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime: env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		env: env.GetString("ENV", "devlopment"),
		mail: mailConfig{
			exp: time.Hour * 24 * 3, // 3 days
		},
	}
	
	// Logger
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	// Database
	db, err := db.New(
		cfg.db.addr,
		cfg.db.maxOpenConns,
		cfg.db.maxIdleConns,
		cfg.db.maxIdleTime,
	)
	if err != nil {
		logger.Fatal(err)
	}

	defer db.Close()
	logger.Info("db connection established")

	store := store.NewStorage(db)

	app := &application{
		config: cfg,
		store: store,
		logger: logger,
	}

	mux := app.mount()

	logger.Fatal(app.run(mux))
}