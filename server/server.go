package server

import (
	"context"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/stuckinforloop/fabrik/db"
	"github.com/stuckinforloop/fabrik/deps/hclient"
	"github.com/stuckinforloop/fabrik/deps/id"
	"github.com/stuckinforloop/fabrik/deps/logger"
	"github.com/stuckinforloop/fabrik/deps/timeutils"
	"github.com/stuckinforloop/fabrik/internal/sources/bundle"
	"github.com/stuckinforloop/fabrik/server/api"
)

type Server struct {
	logger *zap.Logger
	http   *http.Server
	api    *api.API
	db     *db.DB
	port   string
}

func New(port string) *Server {
	env := os.Getenv("environment")
	logger := logger.New(env)

	var nowFunc timeutils.TimeNow = func() time.Time {
		return time.Now()
	}

	idSource := id.New(
		rand.New(rand.NewSource(time.Now().UnixNano())),
		nowFunc,
		false,
	)
	db, err := db.New(context.Background(), db.Config{
		MasterDSN: os.Getenv("master_dsn"),
		ReaderDSN: os.Getenv("reader_dsn"),
		MaxConns:  10,
	})
	if err != nil {
		logger.Error("create db", zap.Error(err))
	}

	hclient := hclient.NewClient(10)
	if err != nil {
		logger.Error("create hclient", zap.Error(err))
	}

	api := api.New(logger, hclient, db, idSource, nowFunc)

	bundle.Import()

	if port == "" {
		port = "8080"
	}

	return &Server{
		logger: logger,
		api:    api,
		db:     db,
		port:   port,
	}
}

func (s *Server) Start() {
	mux := s.api.Mux()
	s.api.RegisterRoutes()

	s.http = &http.Server{
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 10,
		Handler:      mux,
		Addr:         ":" + s.port,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		s.logger.Info("start server", zap.String("port", s.port))
		if err := s.http.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Error("listen and serve", zap.Error(err))
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	s.logger.Info("shutting down server gracefully...")
	if err := s.Shutdown(shutdownCtx); err != nil {
		s.logger.Error("server shutdown", zap.Error(err))
	}
}

func (s *Server) Shutdown(ctx context.Context) error {
	if err := s.http.Shutdown(ctx); err != nil {
		s.logger.Error("http server shutdown", zap.Error(err))
		return err
	}

	if err := s.db.Close(); err != nil {
		s.logger.Error("database shutdown", zap.Error(err))
		return err
	}

	s.logger.Info("server shutdown complete")
	return nil
}
