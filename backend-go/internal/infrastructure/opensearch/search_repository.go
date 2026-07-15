package opensearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/opensearch-project/opensearch-go/v2"
	"vk-search/internal/domain"
	"vk-search/internal/infrastructure/config"
)

type openSearchRepository struct {
	osClient *opensearch.Client
	dbPool   *pgxpool.Pool
	index    string
}

func NewOpenSearchRepository(osClient *opensearch.Client, pool *pgxpool.Pool, cfg *config.Config) domain.SearchRepository {
	return &openSearchRepository{
		osClient: osClient,
		dbPool:   pool,
		index:    cfg.GetOpenSearchIndex(),
	}
}

type osResponse struct {
	Hits struct {
		Hits []struct {
			Source domain.Post `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}

func (r *openSearchRepository) Search(ctx context.Context, query string, limit int) ([]domain.Post, error) {
	searchQuery := map[string]interface{}{
		"size": limit,
		"query": map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":  query,
				"fields": []string{"title", "content"},
			},
		},
	}
	if query == "" {
		searchQuery["query"] = map[string]interface{}{"match_all": map[string]interface{}{}}
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(searchQuery); err != nil {
		return nil, err
	}

	res, err := r.osClient.Search(
		r.osClient.Search.WithContext(ctx),
		r.osClient.Search.WithIndex(r.index),
		r.osClient.Search.WithBody(&buf),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("opensearch error: %s", res.String())
	}

	var rawResp osResponse
	if err := json.NewDecoder(res.Body).Decode(&rawResp); err != nil {
		return nil, err
	}

	posts := make([]domain.Post, 0, len(rawResp.Hits.Hits))
	for _, hit := range rawResp.Hits.Hits {
		posts = append(posts, hit.Source)
	}

	return posts, nil
}

func (r *openSearchRepository) SaveLog(ctx context.Context, log *domain.SearchLog) error {
	query := `
		INSERT INTO search_logs (user_id, query, mode, limit_value, result_count)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at;
	`
	return r.dbPool.QueryRow(ctx, query,
		log.UserID,
		log.Query,
		log.Mode,
		log.LimitValue,
		log.ResultCount,
	).Scan(&log.ID, &log.CreatedAt)
}