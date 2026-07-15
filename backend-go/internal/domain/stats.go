package domain

import (
	"context"
	"database/sql"
)

type IndexRun struct {
	Status         string
	DocumentsCount int
	ChunksCount    int
	FinishedAt     sql.NullString // Используем sql.NullString для корректной обработки NULL из БД
}

type Stats struct {
	SourcesCount    int
	DocumentsCount  int
	ChunksCount     int
	SearchLogsCount int
	LastIndexRun    IndexRun
}

type StatsRepository interface {
	// Метод теперь возвращает полную статистику, агрегированную из разных таблиц Postgres
	GetGeneralStats(ctx context.Context) (*Stats, error)
}

type StatsUseCase interface {
	GetStats(ctx context.Context) (*Stats, error)
}