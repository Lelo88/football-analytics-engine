package ports

import (
	"context"

	"football-analytics/internal/domain"
)

type TeamQueryFilter struct {
	TeamID      int64
	SeasonLabel string
	LastN       int
	Venue       domain.MatchVenue
}

type OverUnderQueryFilter struct {
	TeamID      int64
	SeasonLabel string
	LastN       int
	Venue       domain.MatchVenue
	Threshold   float64
}

type SeasonSummaryFilter struct {
	TeamID      int64
	SeasonLabel string
	Venue       domain.MatchVenue
}

type TeamAnalyticsReadRepository interface {
	GetTeamForm(ctx context.Context, filter TeamQueryFilter) (domain.TeamForm, error)
	GetGoalsSummary(ctx context.Context, filter TeamQueryFilter) (domain.GoalsSummary, error)
	GetOverUnderSummary(ctx context.Context, filter OverUnderQueryFilter) (domain.OverUnderSummary, error)
	GetSeasonSummaries(ctx context.Context, filter SeasonSummaryFilter) ([]domain.SeasonTeamSummary, error)
}
