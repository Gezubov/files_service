package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Gezubov/file_service/config"
	"github.com/Gezubov/file_service/internal/controller"
	"github.com/Gezubov/file_service/internal/infrastructure/db"
	"github.com/Gezubov/file_service/internal/middlewares"
	"github.com/Gezubov/file_service/internal/repository"
	"github.com/Gezubov/file_service/internal/service"
	"github.com/go-chi/chi"
	_ "github.com/lib/pq"
)

func main() {

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	slog.Info("Loading configuration...")
	config.Load()

	slog.Info("Initializing database...")
	db.InitDB()
	database := db.GetDB()

	fileRepo := repository.NewFileRepository(database)
	fileService := service.NewFileService(fileRepo)
	fileController := controller.NewFileController(fileService)

	r := SetupRoutes(fileController)

	port := config.GetConfig().Server.Port
	serverAddr := ":" + port
	srv := &http.Server{
		Addr:    serverAddr,
		Handler: r,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		slog.Info("Server started", "port", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Server crashed", "error", err)
		}
	}()

	<-stop
	slog.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("Server forced to shutdown", "error", err)
	} else {
		slog.Info("Server exited properly")
	}

	db.CloseDB()
}

func SetupRoutes(fileController *controller.FileController) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middlewares.CorsMiddleware())

	r.Route("/files", func(r chi.Router) {
		r.Post("/", fileController.CreateFile)
		r.Post("/upload", fileController.UploadFile)
		r.Get("/", fileController.GetFiles)
		r.Get("/{id}", fileController.GetFile)
		r.Delete("/{id}", fileController.DeleteFile)
	})
	return r
}
