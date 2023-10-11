package rest

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gerich/s3test/internal/ports"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type Server struct {
	log    *zap.Logger
	app    ports.App
	cfg    *Config
	server *http.Server
}

type Config struct {
	Port          int
	ReadTimeout   time.Duration
	WriteTimeout  time.Duration
	JWTSecret     string
	AllowedUsers  []string
	MaxFileSizeMB int64
}

func NewServer(app ports.App, cfg *Config, log *zap.Logger) *Server {
	s := &Server{app: app, cfg: cfg, log: log.Named("http")}
	s.server = &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      s.router(),
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}

	return s
}

func (s *Server) Run() {
	go func() {
		s.log.Info("http server are starting", zap.Int("port", s.cfg.Port))
		if err := s.server.ListenAndServe(); err != http.ErrServerClosed && err != nil {
			s.log.Fatal("cant start http server", zap.Error(err))
		}
	}()
}

func (s *Server) Stop(ctx context.Context) {
	s.server.SetKeepAlivesEnabled(false)
	if err := s.server.Shutdown(ctx); err != nil {
		s.log.Error("shutdown error on http server: %s", zap.Error(err))
	}
}

func (s *Server) router() http.Handler {
	r := chi.NewRouter()

	r.Use(s.recover, s.jwt())

	r.Route("/health", func(router chi.Router) {
		router.Get("/ready", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(200)
			if _, err := w.Write([]byte("{status: ok}")); err != nil {
				s.log.Error("error while http request", zap.Error(err))
			}
		})
	})

	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/upload", s.upload)
		r.Get("/download/{id}", s.download)
		r.Get("/list", s.list)
	})

	return r
}
