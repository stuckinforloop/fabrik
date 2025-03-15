package api

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"

	"github.com/stuckinforloop/fabrik/internal/sources"
)

type contextKey string

const (
	ContextKeyLogger        contextKey = "logger"
	ContextKeyRequestID     contextKey = "request_id"
	ContextKeySourceService contextKey = "dao"
)

func (a *API) WithLogger() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			reqID := a.id.MustULID()
			logger := a.logger.With(
				zap.String("request_id", reqID),
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
			)

			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			now := time.Now()

			defer func() {
				logger.Info("handled",
					zap.Int("status_code", ww.Status()),
					zap.Duration("duration", time.Since(now)),
				)
			}()

			ctx := r.Context()
			ctx = context.WithValue(ctx, ContextKeyRequestID, reqID)
			ctx = context.WithValue(ctx, ContextKeyLogger, logger)
			next.ServeHTTP(ww, r.WithContext(ctx))
		}

		return http.HandlerFunc(fn)
	}
}

func (a *API) WithSourceService() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			srv := sources.NewSourceService(a.db, a.hclient, a.id, a.logger, a.nowFunc)
			ctx = context.WithValue(ctx, ContextKeySourceService, srv)
			next.ServeHTTP(w, r.WithContext(ctx))
		}

		return http.HandlerFunc(fn)
	}
}
