package main

import (
	stdhttp "net/http"

	"go.uber.org/zap"

	"github.com/yourname/energy-usage-api/internal/config"
	"github.com/yourname/energy-usage-api/internal/db"
	"github.com/yourname/energy-usage-api/internal/domain"
	apihttp "github.com/yourname/energy-usage-api/internal/http"
	"github.com/yourname/energy-usage-api/internal/http/handlers"
)

func main() {
	cfg := config.FromEnv()
	log, _ := zap.NewProduction()
	defer log.Sync()

	pg, err := db.Connect(cfg.DatabaseURL)
	if err != nil { log.Fatal("db connect", zap.Error(err)) }

	repo := domain.NewRepository(pg)
	svc  := domain.NewService(repo)

	ah := handlers.NewAuthHandler(svc, cfg.JWTSecret)
	mh := &handlers.MetersHandler{Svc: svc}
	rh := &handlers.ReadingsHandler{Svc: svc}
	router := apihttp.NewRouter(ah, mh, rh, cfg.JWTSecret)

	log.Info("server starting", zap.String("addr", cfg.Addr))
	if err := stdhttp.ListenAndServe(cfg.Addr, router); err != nil {
		log.Fatal("server error", zap.Error(err))
	}
}
