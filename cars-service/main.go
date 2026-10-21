package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/LeUrok/DS-lab2/cars-service/internal/api"
	"github.com/LeUrok/DS-lab2/cars-service/internal/auth"
	"github.com/LeUrok/DS-lab2/cars-service/internal/config"
	"github.com/LeUrok/DS-lab2/cars-service/internal/db"
	"github.com/LeUrok/DS-lab2/cars-service/internal/repository"
	"github.com/LeUrok/DS-lab2/cars-service/internal/service"
)

func main() {
	cfg := config.Load()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	issuer := "https://" + cfg.AuthDomain + "/"
	jwksURL := issuer + ".well-known/jwks.json"
	mw, err := auth.NewMiddleware(context.Background(), issuer, jwksURL, cfg.AuthAudience)
	if err != nil {
		log.Fatalf("failed to init auth middleware: %v", err)
	}

	pool, err := db.New(ctx, cfg.DBURL())
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	defer pool.Close()

	repo := repository.NewCarRepository(pool)
	svc := service.NewCarService(repo)
	h := api.NewCarHandler(svc)
	router := api.NewRouter(h, mw)

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	go func() {
		log.Printf("cars-service listening on :%s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("shutdown error: %v", err)
	}
	log.Println("cars-service stopped")
}
