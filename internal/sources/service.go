package sources

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/stuckinforloop/fabrik/db"
	"github.com/stuckinforloop/fabrik/deps/hclient"
	"github.com/stuckinforloop/fabrik/deps/id"
	"github.com/stuckinforloop/fabrik/deps/timeutils"
)

type SourceService struct {
	symKey  string
	db      *db.DB
	HClient *hclient.Client
	ID      *id.Source
	Logger  *zap.Logger
	NowFunc timeutils.TimeNow
}

func NewSourceService(
	symKey string,
	db *db.DB,
	hclient *hclient.Client,
	id *id.Source,
	logger *zap.Logger,
	nowFunc timeutils.TimeNow,
) *SourceService {
	return &SourceService{
		symKey:  symKey,
		db:      db,
		HClient: hclient,
		ID:      id,
		Logger:  logger,
		NowFunc: nowFunc,
	}
}

func (s *SourceService) CreateSource(ctx context.Context, config *Config) error {
	if !config.Kind.Valid() {
		return fmt.Errorf("invalid source kind: %s", config.Kind)
	}

	id, err := s.ID.Generate()
	if err != nil {
		return fmt.Errorf("generate id: %w", err)
	}

	creds, err := json.Marshal(config.Credentials)
	if err != nil {
		return fmt.Errorf("marshal credentials: %w", err)
	}

	config.ID = id
	config.CreatedAt = s.NowFunc()
	config.UpdatedAt = s.NowFunc()

	query := `
		INSERT INTO sources (id, name, kind, config, credentials, sync_frequency, created_at, updated_at)
		VALUES ($1, $2, $3, $4, pgp_sym_encrypt_bytea($5, $6), $7, $8, $9)
	`

	if err := s.db.RW().Exec(ctx, query,
		config.ID,
		config.Name,
		config.Kind,
		config.Config,
		creds,
		s.symKey,
		config.SyncFrequency,
		config.CreatedAt,
		config.UpdatedAt,
	); err != nil {
		return fmt.Errorf("insert source: %w", err)
	}

	return nil
}

func (s *SourceService) GetSource(ctx context.Context, id string, kind Kind) (*Config, error) {
	config := &Config{}

	var creds []byte
	var syncFrequency time.Duration

	query := `
		SELECT id, name, kind, config, pgp_sym_decrypt_bytea(credentials, $1), sync_frequency, created_at, updated_at
		FROM sources
		WHERE id = $2
		AND kind = $3
	`
	if err := s.db.RW().QueryRow(ctx, query, s.symKey, id, kind).Scan(
		&config.ID,
		&config.Name,
		&config.Kind,
		&config.Config,
		&creds,
		&syncFrequency,
		&config.CreatedAt,
		&config.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("get source: %w", err)
	}

	if err := json.Unmarshal(creds, &config.Credentials); err != nil {
		return nil, fmt.Errorf("unmarshal credentials: %w", err)
	}

	config.SyncFrequency = syncFrequency.String()

	return config, nil
}
