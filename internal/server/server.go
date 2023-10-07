package server

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/omaily/JWT/config"
	"github.com/omaily/JWT/internal/server/middlewareRedef"
	"github.com/omaily/JWT/internal/storage"
)

type apiServer struct {
	conf    *config.HTTPServer
	storage *storage.Storage
}

func NewServer(conf *config.HTTPServer, instance *storage.Storage) (*apiServer, error) {
	if conf == nil {
		return nil, errors.New("configuration files are not initialized")
	}
	if conf.Address == "" || conf.Port == "" {
		return nil, errors.New("configuration files address cannot be blank")
	}

	return &apiServer{
		conf:    conf,
		storage: instance,
	}, nil
}

func (s *apiServer) Start(logger *slog.Logger) error {

	ser := &http.Server{
		Addr:         s.conf.Port,
		Handler:      s.router(logger),
		ReadTimeout:  s.conf.Timeout * time.Second,
		WriteTimeout: s.conf.Timeout * 2 * time.Second,
		IdleTimeout:  s.conf.IdleTimeout * time.Second,
	}

	go func() {
		slog.Info("starting server", slog.String("addres", s.conf.Address))
		if err := ser.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("не маслает", slog.StringValue(err.Error()))
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-stop

	slog.Info("stopping server due to syscall")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return ser.Shutdown(ctx)
}

func (s *apiServer) router(logger *slog.Logger) http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Recoverer)
	router.Use(middlewareRedef.New(logger))

	router.Post("/", s.helloWorld())
	router.Post("/api/auth/createAccount", s.authorized(s.storage.CreateAccount))
	router.Post("/api/auth/login", s.authorized(s.storage.LoginAccount))

	return router
}
