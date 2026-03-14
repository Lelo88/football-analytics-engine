package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"football-analytics/internal/domain"
	"football-analytics/internal/ports"
)

type TeamAnalyticsReadRepository struct {
	db *sql.DB
}

func NewTeamAnalyticsReadRepository(db *sql.DB) *TeamAnalyticsReadRepository {
	return &TeamAnalyticsReadRepository{db: db}
}

func (repository *TeamAnalyticsReadRepository) GetTeamForm(ctx context.Context, filter ports.TeamQueryFilter) (domain.TeamForm, error) {
	baseQuery, args := buildTeamMatchesBaseQuery(filter.TeamID, normalizeVenue(filter.Venue), filter.SeasonLabel)
	query := wrapWithLastN(baseQuery, filter.LastN)

	form := domain.TeamForm{}
	err := repository.db.QueryRowContext(ctx, `
WITH filtered_matches AS (`+query+`)
SELECT
	COUNT(*) AS matches_played,
	COALESCE(SUM(CASE WHEN goals_for > goals_against THEN 1 ELSE 0 END), 0) AS wins,
	COALESCE(SUM(CASE WHEN goals_for = goals_against THEN 1 ELSE 0 END), 0) AS draws,
	COALESCE(SUM(CASE WHEN goals_for < goals_against THEN 1 ELSE 0 END), 0) AS losses,
	COALESCE(SUM(
		CASE
			WHEN goals_for > goals_against THEN 3
			WHEN goals_for = goals_against THEN 1
			ELSE 0
		END
	), 0) AS points,
	COALESCE(SUM(goals_for), 0) AS goals_for,
	COALESCE(SUM(goals_against), 0) AS goals_against
FROM filtered_matches`, args...).Scan(
		&form.MatchesPlayed,
		&form.Wins,
		&form.Draws,
		&form.Losses,
		&form.Points,
		&form.GoalsFor,
		&form.GoalsAgainst,
	)
	if err != nil {
		return domain.TeamForm{}, fmt.Errorf("query team form: %w", err)
	}

	return form, nil
}

func (repository *TeamAnalyticsReadRepository) GetGoalsSummary(ctx context.Context, filter ports.TeamQueryFilter) (domain.GoalsSummary, error) {
	baseQuery, args := buildTeamMatchesBaseQuery(filter.TeamID, normalizeVenue(filter.Venue), filter.SeasonLabel)
	query := wrapWithLastN(baseQuery, filter.LastN)

	summary := domain.GoalsSummary{}
	err := repository.db.QueryRowContext(ctx, `
WITH filtered_matches AS (`+query+`)
SELECT
	COUNT(*) AS matches_played,
	COALESCE(SUM(goals_for), 0) AS goals_for,
	COALESCE(SUM(goals_against), 0) AS goals_against
FROM filtered_matches`, args...).Scan(
		&summary.MatchesPlayed,
		&summary.GoalsFor,
		&summary.GoalsAgainst,
	)
	if err != nil {
		return domain.GoalsSummary{}, fmt.Errorf("query goals summary: %w", err)
	}

	return summary, nil
}

func (repository *TeamAnalyticsReadRepository) GetOverUnderSummary(ctx context.Context, filter ports.OverUnderQueryFilter) (domain.OverUnderSummary, error) {
	baseQuery, args := buildTeamMatchesBaseQuery(filter.TeamID, normalizeVenue(filter.Venue), filter.SeasonLabel)
	args = append(args, filter.Threshold)
	thresholdPosition := len(args)
	query := wrapWithLastN(baseQuery, filter.LastN)

	summary := domain.OverUnderSummary{}
	err := repository.db.QueryRowContext(ctx, `
WITH filtered_matches AS (`+query+`)
SELECT
	COUNT(*) AS matches_played,
	COALESCE(SUM(CASE WHEN (goals_for + goals_against) > $`+fmt.Sprintf("%d", thresholdPosition)+` THEN 1 ELSE 0 END), 0) AS over_count,
	COALESCE(SUM(CASE WHEN (goals_for + goals_against) <= $`+fmt.Sprintf("%d", thresholdPosition)+` THEN 1 ELSE 0 END), 0) AS under_equal_count
FROM filtered_matches`, args...).Scan(
		&summary.MatchesPlayed,
		&summary.OverCount,
		&summary.UnderEqualCount,
	)
	if err != nil {
		return domain.OverUnderSummary{}, fmt.Errorf("query over/under summary: %w", err)
	}

	summary.Threshold = filter.Threshold

	return summary, nil
}

func (repository *TeamAnalyticsReadRepository) GetSeasonSummaries(ctx context.Context, filter ports.SeasonSummaryFilter) ([]domain.SeasonTeamSummary, error) {
	baseQuery, args := buildTeamMatchesBaseQuery(filter.TeamID, normalizeVenue(filter.Venue), filter.SeasonLabel)

	rows, err := repository.db.QueryContext(ctx, `
WITH filtered_matches AS (`+baseQuery+`)
SELECT
	season_label,
	COUNT(*) AS matches_played,
	COALESCE(SUM(
		CASE
			WHEN goals_for > goals_against THEN 3
			WHEN goals_for = goals_against THEN 1
			ELSE 0
		END
	), 0) AS points,
	COALESCE(SUM(goals_for), 0) AS goals_for,
	COALESCE(SUM(goals_against), 0) AS goals_against
FROM filtered_matches
GROUP BY season_label
ORDER BY season_label`, args...)
	if err != nil {
		return nil, fmt.Errorf("query season summaries: %w", err)
	}
	defer rows.Close()

	summaries := make([]domain.SeasonTeamSummary, 0)
	for rows.Next() {
		summary := domain.SeasonTeamSummary{}
		err = rows.Scan(
			&summary.SeasonLabel,
			&summary.MatchesPlayed,
			&summary.Points,
			&summary.GoalsFor,
			&summary.GoalsAgainst,
		)
		if err != nil {
			return nil, fmt.Errorf("scan season summary: %w", err)
		}

		summaries = append(summaries, summary)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate season summaries: %w", err)
	}

	return summaries, nil
}

func buildTeamMatchesBaseQuery(teamID int64, venue domain.MatchVenue, seasonLabel string) (string, []any) {
	args := []any{teamID}

	builder := strings.Builder{}
	builder.WriteString(`
SELECT
	s.label AS season_label,
	m.match_date,
	CASE WHEN m.home_team_id = $1 THEN m.home_goals ELSE m.away_goals END AS goals_for,
	CASE WHEN m.home_team_id = $1 THEN m.away_goals ELSE m.home_goals END AS goals_against
FROM matches m
INNER JOIN seasons s ON s.id = m.season_id
WHERE m.home_goals IS NOT NULL
	AND m.away_goals IS NOT NULL
`)

	switch venue {
	case domain.MatchVenueHome:
		builder.WriteString("\tAND m.home_team_id = $1\n")
	case domain.MatchVenueAway:
		builder.WriteString("\tAND m.away_team_id = $1\n")
	default:
		builder.WriteString("\tAND (m.home_team_id = $1 OR m.away_team_id = $1)\n")
	}

	if seasonLabel != "" {
		args = append(args, seasonLabel)
		builder.WriteString("\tAND s.label = $")
		builder.WriteString(fmt.Sprintf("%d", len(args)))
		builder.WriteString("\n")
	}

	builder.WriteString("ORDER BY m.match_date DESC")

	return builder.String(), args
}

func wrapWithLastN(baseQuery string, lastN int) string {
	if lastN <= 0 {
		return baseQuery
	}

	return baseQuery + "\nLIMIT " + fmt.Sprintf("%d", lastN)
}

func normalizeVenue(venue domain.MatchVenue) domain.MatchVenue {
	if venue == "" {
		return domain.MatchVenueAll
	}

	return venue
}
