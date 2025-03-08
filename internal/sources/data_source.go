package sources

import (
	"context"
	"fmt"
	"sync"
)

type DataSource interface {
	Open(srv *SourceService, auth map[string]any, cfg map[string]any) error
	Fetch(ctx context.Context) ([]byte, error)
}

var (
	dataSourcesMu sync.RWMutex
	dataSources   = make(map[Kind]DataSource)
)

func RegisterDataSource(k Kind, dataSource DataSource) {
	dataSourcesMu.Lock()
	defer dataSourcesMu.Unlock()

	if dataSource == nil {
		panic("data source cannot be nil")
	}

	if _, ok := dataSources[k]; ok {
		panic("data source already registered")
	}

	dataSources[k] = dataSource
}

func GetDataSource(k Kind) (DataSource, error) {
	dataSourcesMu.RLock()
	defer dataSourcesMu.RUnlock()

	dataSource, ok := dataSources[k]
	if !ok {
		return nil, fmt.Errorf("unknown data source: %s", k)
	}

	if dataSource == nil {
		return nil, fmt.Errorf("data source not found: %s", k)
	}

	return dataSource, nil
}
