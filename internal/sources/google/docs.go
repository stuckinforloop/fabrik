package google

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-viper/mapstructure/v2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/docs/v1"
	"google.golang.org/api/option"

	"github.com/stuckinforloop/fabrik/internal/sources"
)

func init() {
	sources.RegisterDataSource(sources.KindGoogleDocs, &GDocs{})
}

type GDocs struct {
	Config struct {
		DocID string `mapstructure:"doc_id"`
	} `mapstructure:"config"`
	Credentials struct {
		ServiceAccount string `mapstructure:"service_account"`
	} `mapstructure:"credentials"`
}

func (s *GDocs) Open(_ *sources.SourceService, auth map[string]any, cfg map[string]any) error {
	if err := mapstructure.Decode(cfg, &s.Config); err != nil {
		return fmt.Errorf("decode config: %w", err)
	}

	if err := mapstructure.Decode(auth, &s.Credentials); err != nil {
		return fmt.Errorf("decode credentials: %w", err)
	}

	return nil
}

func (s *GDocs) Fetch(ctx context.Context) ([]byte, error) {
	s.Credentials.ServiceAccount = strings.Replace(s.Credentials.ServiceAccount, "\n", "", -1)
	jsonKey := []byte(s.Credentials.ServiceAccount)
	config, err := google.JWTConfigFromJSON(
		jsonKey,
		docs.DocumentsReadonlyScope,
	)
	if err != nil {
		return nil, fmt.Errorf("create jwt config from json key: %w", err)
	}

	client := config.Client(ctx)
	service, err := docs.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("initialize service: %w", err)
	}

	docID := s.Config.DocID
	doc, err := service.Documents.Get(docID).Do()
	if err != nil {
		return nil, fmt.Errorf("get document: %w", err)
	}

	var content strings.Builder
	for _, elem := range doc.Body.Content {
		if elem.Paragraph != nil {
			for _, paraElem := range elem.Paragraph.Elements {
				if paraElem.TextRun != nil {
					content.WriteString(paraElem.TextRun.Content)
				}
			}
		}
	}

	return []byte(content.String()), nil
}
