package linear

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-viper/mapstructure/v2"

	"github.com/stuckinforloop/fabrik/deps/hclient"
	"github.com/stuckinforloop/fabrik/internal/sources"
)

func init() {
	sources.RegisterDataSource(sources.KindLinear, &Linear{})
}

const (
	baseURL = "https://api.linear.app/graphql"
)

type Linear struct {
	hClient     *hclient.Client
	Credentials struct {
		APIKey string `mapstructure:"api_key"`
	} `mapstructure:"credentials"`
}

func (s *Linear) Open(srv *sources.SourceService, auth map[string]any, _ map[string]any) error {
	if err := mapstructure.Decode(auth, &s.Credentials); err != nil {
		return fmt.Errorf("decode credentials: %w", err)
	}

	s.hClient = srv.HClient

	return nil
}

func (s *Linear) Fetch(ctx context.Context) ([]byte, error) {
	body := []byte(`{ "query": "{ issues { nodes { id title description } } }" }`)
	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": s.Credentials.APIKey,
	}

	resp, err := s.hClient.Do(ctx, baseURL, http.MethodPost, headers, body)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	issueResp := Response{}
	if err := json.NewDecoder(resp.Body).Decode(&issueResp); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	b, err := json.Marshal(issueResp)
	if err != nil {
		return nil, fmt.Errorf("marshal response")
	}

	return b, nil
}

type Response struct {
	Data struct {
		Issues struct {
			Nodes []struct {
				ID          string `json:"id"`
				Title       string `json:"title"`
				Description string `json:"description"`
			} `json:"nodes"`
		} `json:"issues"`
	} `json:"data"`
}
