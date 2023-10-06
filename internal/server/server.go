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
	"github.com/go-chi/render"
	"github.com/omaily/JWT/config"
)

type apiServer struct {
	conf *config.HTTPServer
}

func NewServer(conf *config.HTTPServer) (*apiServer, error) {
	if conf == nil {
		return nil, errors.New("configuration files are not initialized")
	}
	if conf.Address == "" || conf.Port == "" {
		return nil, errors.New("configuration files address cannot be blank")
	}

	return &apiServer{
		conf: conf,
	}, nil
}

func (s *apiServer) Start() error {

	ser := &http.Server{
		Addr:         s.conf.Port,
		Handler:      s.router(),
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

func (s *apiServer) router() http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Recoverer)
	router.Post("/", s.helloWorld())

	return router
}

func (s *apiServer) helloWorld() http.HandlerFunc {
	return func(write http.ResponseWriter, request *http.Request) {
		render.JSON(write, request, "Hello World")
		return
	}
}
