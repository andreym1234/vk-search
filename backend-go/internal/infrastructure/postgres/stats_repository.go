package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"vk-search/internal/domain"
)

type statsRepository struct {
	pool *pgxpool.Pool
}

func NewStatsRepository(pool *pgxpool.Pool) domain.StatsRepository {
	return &statsRepository{
		pool: pool,
	}
}

func (r *statsRepository) GetGeneralStats(ctx context.Context) (*domain.Stats, error) {
	stats := &domain.Stats{}

	// 1. Собираем все счетчики из PostgreSQL (включая чанки)
	queryCounts := `
		SELECT 
			(SELECT COUNT(*) FROM sources),
			(SELECT COUNT(*) FROM documents),
			(SELECT COUNT(*) FROM chunks),
			(SELECT COUNT(*) FROM search_logs)
	`
	err := r.pool.QueryRow(ctx, queryCounts).Scan(
		&stats.SourcesCount,
		&stats.DocumentsCount,
		&stats.ChunksCount,
		&stats.SearchLogsCount,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get counts: %w", err)
	}

	// 2. Получаем данные последнего запуска индексации
	queryLastRun := `
		SELECT status, documents_count, chunks_count, finished_at 
		FROM index_runs 
		ORDER BY finished_at DESC 
		LIMIT 1
	`
	err = r.pool.QueryRow(ctx, queryLastRun).Scan(
		&stats.LastIndexRun.Status,
		&stats.LastIndexRun.DocumentsCount,
		&stats.LastIndexRun.ChunksCount,
		&stats.LastIndexRun.FinishedAt, // Теперь используем sql.NullString
	)

	if err != nil {
		// Если данных об индексации еще нет (таблица пуста), не считаем это ошибкой
		if errors.Is(err, pgx.ErrNoRows) {
			return stats, nil
		}
		return nil, fmt.Errorf("failed to get last index run: %w", err)
	}

	return stats, nil
}