package sources

import (
	"time"
)

// Config represents a source configuration.
type Config struct {
	ID            string         `json:"id" mapstructure:"id"`
	Name          string         `json:"name" mapstructure:"name"`
	Kind          Kind           `json:"kind" mapstructure:"kind"`
	Config        map[string]any `json:"config" mapstructure:"config"`
	Credentials   map[string]any `json:"credentials" mapstructure:"credentials"`
	Status        Status         `json:"status" mapstructure:"status"`
	SyncFrequency string         `json:"sync_frequency" mapstructure:"sync_frequency"`
	CreatedAt     time.Time      `json:"created_at" mapstructure:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at" mapstructure:"updated_at"`
	SyncedAt      *time.Time     `json:"synced_at" mapstructure:"synced_at"`
	DeletedAt     *time.Time     `json:"deleted_at" mapstructure:"deleted_at"`
}

// Status represents the status of a source.
type Status string

const (
	StatusActive   Status = "active"
	StatusInactive Status = "inactive"
	StatusPending  Status = "pending"
	StatusError    Status = "error"
)

// Kind represents a type of source
type Kind string

const (
	KindGoogleDocs Kind = "googledocs"
	KindSlack      Kind = "slack"
	KindLinear     Kind = "linear"
)

var validSources = map[Kind]struct{}{
	KindGoogleDocs: {},
	KindSlack:      {},
	KindLinear:     {},
}

func (k Kind) Valid() bool {
	_, ok := validSources[k]
	return ok
}
