package opensearch

import (
	"github.com/opensearch-project/opensearch-go/v2"
	"vk-search/internal/infrastructure/config"
)

func NewOpenSearchClient(cfg *config.Config) (*opensearch.Client, error) {
	return opensearch.NewClient(opensearch.Config{
		Addresses: []string{cfg.GetOpenSearchURL()},
	})
}