package ports

import (
	"context"
	"time"

	"football-analytics/internal/domain"
)

type SourceMatchRow struct {
	CompetitionCode string
	CompetitionName string
	Country         string
	SeasonLabel     string
	MatchDate       time.Time
	HomeTeamName    string
	AwayTeamName    string
	HomeGoals       *int
	AwayGoals       *int
	HomeWinOdds     *float64
	DrawOdds        *float64
	AwayWinOdds     *float64
}

type SourceReader interface {
	ReadMatches(ctx context.Context, sourceURL string) ([]SourceMatchRow, error)
}

type CompetitionRepository interface {
	CreateOrGet(ctx context.Context, name string, country string) (domain.Competition, error)
}

type SeasonRepository interface {
	CreateOrGet(ctx context.Context, competitionID int64, label string) (domain.Season, error)
}

type TeamRepository interface {
	CreateOrGet(ctx context.Context, name string) (domain.Team, error)
}

type MatchUpsertParams struct {
	CompetitionID int64
	SeasonID      int64
	MatchDate     time.Time
	HomeTeamID    int64
	AwayTeamID    int64
	HomeGoals     *int
	AwayGoals     *int
}

type MatchUpsertResult struct {
	Match    domain.Match
	Inserted bool
	Updated  bool
}

type MatchOddsUpsertParams struct {
	MatchID int64
	HomeWin *float64
	Draw    *float64
	AwayWin *float64
}

type MatchRepository interface {
	UpsertMatch(ctx context.Context, params MatchUpsertParams) (MatchUpsertResult, error)
	UpsertMatchOdds(ctx context.Context, params MatchOddsUpsertParams) error
}

type IngestionSaveResult struct {
	Inserted bool
	Updated  bool
}

type IngestionRunRepository interface {
	Create(ctx context.Context, source string) (domain.IngestionRun, error)
	MarkSuccess(ctx context.Context, runID int64, rowsProcessed int, rowsInserted int, rowsUpdated int) error
	MarkFailed(ctx context.Context, runID int64, rowsProcessed int, rowsInserted int, rowsUpdated int, errorMessage string) error
}
