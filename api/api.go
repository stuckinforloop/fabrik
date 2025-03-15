package api

import (
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/stuckinforloop/fabrik/db"
	"github.com/stuckinforloop/fabrik/deps/hclient"
	"github.com/stuckinforloop/fabrik/deps/id"
	"github.com/stuckinforloop/fabrik/deps/timeutils"
)

type API struct {
	mux     *chi.Mux
	symKey  string
	logger  *zap.Logger
	hclient *hclient.Client
	db      *db.DB
	id      *id.Source
	nowFunc timeutils.TimeNow
}

func New(
	symKey string,
	logger *zap.Logger,
	hclient *hclient.Client,
	db *db.DB,
	id *id.Source,
	nowFunc timeutils.TimeNow,
) *API {
	mux := chi.NewMux()
	return &API{
		mux:     mux,
		symKey:  symKey,
		logger:  logger,
		hclient: hclient,
		db:      db,
		id:      id,
		nowFunc: nowFunc,
	}
}

func (a *API) Mux() *chi.Mux {
	return a.mux
}
