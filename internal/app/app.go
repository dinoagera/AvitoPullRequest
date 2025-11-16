package app

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dinoagera/AvitoPullRequest/config"
	"github.com/dinoagera/AvitoPullRequest/internal/http/handler"
	"github.com/dinoagera/AvitoPullRequest/internal/repository/postgres"
	"github.com/dinoagera/AvitoPullRequest/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Run(cfg *config.Config, l *slog.Logger) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	pool, err := pgxpool.New(ctx, cfg.StoragePath)
	if err != nil {
		log.Fatalf("can't connect to postgresql: %v", err)
	}
	defer pool.Close()
	teamRepo := postgres.NewTeamRepository(pool)
	userRepo := postgres.NewUserRepository(pool)
	prRepo := postgres.NewPRRepository(pool)
	statsRepo := postgres.NewStatsRepository(pool)

	teamSvc := service.NewTeamService(l, teamRepo)
	userSvc := service.NewUserService(l, userRepo, prRepo)
	statsSvc := service.NewStatsService(l, statsRepo)

	prSvc := service.NewPRService(l, userRepo, teamRepo, prRepo)
	teamHandler := handler.NewTeamHandler(l, teamSvc)
	userHandler := handler.NewUserHandler(l, userSvc)
	prHandler := handler.NewPRHandler(l, prSvc)
	statsHandler := handler.NewStatsHandler(l, statsSvc)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Post("/team/add", teamHandler.AddTeam)
	r.Get("/team/get", teamHandler.GetTeam)
	r.Post("/users/setIsActive", userHandler.SetActive)
	r.Get("/users/getReview", userHandler.GetReview)
	r.Post("/pullRequest/create", prHandler.CreatePR)
	r.Post("/pullRequest/merge", prHandler.MergePR)
	r.Post("/pullRequest/reassign", prHandler.ReassignReviewer)
	r.Get("/stats/reviewers", statsHandler.GetReviewerStats)
	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}
	go func() {
		fmt.Println("server started on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server failed: %v", err)
			os.Exit(1)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	fmt.Println("shutting down server...")
	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("server shutdown failed: %v", err)
	}
	l.Info("server exited")
}
