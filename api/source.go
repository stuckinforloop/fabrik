package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/stuckinforloop/fabrik/internal/sources"
)

func (a *API) NewSource(_ http.ResponseWriter, r *http.Request) *Response {
	srv, ok := r.Context().Value(ContextKeySourceService).(*sources.SourceService)
	if !ok {
		a.logger.Error("get source service from context")
		return &Response{
			StatusCode: http.StatusInternalServerError,
			Success:    false,
			Err: &ResponseError{
				Code:    "internal_server_error",
				Message: "failed to get source service",
			},
		}
	}

	var payload sources.Config
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		a.logger.Error("failed to decode payload", zap.Error(err))
		return &Response{
			StatusCode: http.StatusBadRequest,
			Success:    false,
			Err: &ResponseError{
				Code:    "bad_request",
				Message: "failed to decode payload",
			},
		}
	}

	if !payload.Kind.Valid() {
		return &Response{
			StatusCode: http.StatusBadRequest,
			Success:    false,
			Err: &ResponseError{
				Code:    "bad_request",
				Message: "invalid source kind",
			},
		}
	}

	if err := srv.CreateSource(r.Context(), &payload); err != nil {
		a.logger.Error("failed to create source", zap.Error(err))

		return &Response{
			StatusCode: http.StatusInternalServerError,
			Success:    false,
			Err: &ResponseError{
				Code:    "internal_server_error",
				Message: "failed to create source",
			},
		}
	}

	return &Response{
		StatusCode: http.StatusOK,
		Success:    true,
		Data: map[string]any{
			"id": payload.ID,
		},
	}
}

func (a *API) Fetch(_ http.ResponseWriter, r *http.Request) *Response {
	srv, ok := r.Context().Value(ContextKeySourceService).(*sources.SourceService)
	if !ok {
		a.logger.Error("get source service from context")
		return &Response{
			StatusCode: http.StatusInternalServerError,
			Success:    false,
			Err: &ResponseError{
				Code:    "internal_server_error",
				Message: "failed to get source service",
			},
		}
	}
	k := chi.URLParam(r, "kind")
	kind := sources.Kind(k)
	if !kind.Valid() {
		return &Response{
			StatusCode: http.StatusBadRequest,
			Success:    false,
			Err: &ResponseError{
				Code:    "bad_request",
				Message: "invalid source kind",
			},
		}
	}
	id := chi.URLParam(r, "id")

	dataSource, err := sources.GetDataSource(kind)
	if err != nil {
		a.logger.Error("failed to get data source", zap.Error(err))
		return &Response{
			StatusCode: http.StatusInternalServerError,
			Success:    false,
			Err: &ResponseError{
				Code:    "internal_server_error",
				Message: "failed to get data source",
			},
		}
	}

	source, err := srv.GetSource(r.Context(), id, kind)
	if err != nil {
		a.logger.Error("failed to get source", zap.Error(err))
		return &Response{
			StatusCode: http.StatusInternalServerError,
			Success:    false,
			Err: &ResponseError{
				Code:    "internal_server_error",
				Message: "failed to get source",
			},
		}
	}

	if err := dataSource.Open(srv, source.Credentials, source.Config); err != nil {
		a.logger.Error("failed to open data source", zap.Error(err))
		return &Response{
			StatusCode: http.StatusInternalServerError,
			Success:    false,
			Err: &ResponseError{
				Code:    "internal_server_error",
				Message: "failed to open data source",
			},
		}
	}

	body, err := dataSource.Fetch(r.Context())
	if err != nil {
		a.logger.Error("failed to fetch data source", zap.Error(err))
		return &Response{
			StatusCode: http.StatusInternalServerError,
			Success:    false,
			Err: &ResponseError{
				Code:    "internal_server_error",
				Message: "failed to fetch data source",
			},
		}
	}

	return &Response{
		StatusCode: http.StatusOK,
		Success:    true,
		Data:       string(body),
	}
}
