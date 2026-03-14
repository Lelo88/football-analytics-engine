package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"football-analytics/internal/domain"
	"football-analytics/internal/ports"
)

type CompetitionRepository struct {
	db *sql.DB
}

func NewCompetitionRepository(db *sql.DB) *CompetitionRepository {
	return &CompetitionRepository{db: db}
}

func (repository *CompetitionRepository) CreateOrGet(ctx context.Context, name string, country string) (domain.Competition, error) {
	competition := domain.Competition{}
	err := repository.db.QueryRowContext(
		ctx,
		`INSERT INTO competitions (name, country)
		 VALUES ($1, NULLIF($2, ''))
		 ON CONFLICT (name, (COALESCE(country, '')))
		 DO UPDATE SET name = EXCLUDED.name
		 RETURNING id, name, COALESCE(country, ''), created_at`,
		name,
		country,
	).Scan(&competition.ID, &competition.Name, &competition.Country, &competition.CreatedAt)
	if err != nil {
		return domain.Competition{}, fmt.Errorf("upsert competition: %w", err)
	}

	return competition, nil
}

type SeasonRepository struct {
	db *sql.DB
}

func NewSeasonRepository(db *sql.DB) *SeasonRepository {
	return &SeasonRepository{db: db}
}

func (repository *SeasonRepository) CreateOrGet(ctx context.Context, competitionID int64, label string) (domain.Season, error) {
	season := domain.Season{}
	err := repository.db.QueryRowContext(
		ctx,
		`INSERT INTO seasons (competition_id, label)
		 VALUES ($1, $2)
		 ON CONFLICT ON CONSTRAINT seasons_competition_id_label_key
		 DO UPDATE SET label = EXCLUDED.label
		 RETURNING id, competition_id, label, created_at`,
		competitionID,
		label,
	).Scan(&season.ID, &season.CompetitionID, &season.Label, &season.CreatedAt)
	if err != nil {
		return domain.Season{}, fmt.Errorf("upsert season: %w", err)
	}

	return season, nil
}

type TeamRepository struct {
	db *sql.DB
}

func NewTeamRepository(db *sql.DB) *TeamRepository {
	return &TeamRepository{db: db}
}

func (repository *TeamRepository) CreateOrGet(ctx context.Context, name string) (domain.Team, error) {
	team := domain.Team{}
	err := repository.db.QueryRowContext(
		ctx,
		`INSERT INTO teams (name)
		 VALUES ($1)
		 ON CONFLICT (name)
		 DO UPDATE SET name = EXCLUDED.name
		 RETURNING id, name, created_at`,
		name,
	).Scan(&team.ID, &team.Name, &team.CreatedAt)
	if err != nil {
		return domain.Team{}, fmt.Errorf("upsert team: %w", err)
	}

	return team, nil
}

type MatchRepository struct {
	db *sql.DB
}

func NewMatchRepository(db *sql.DB) *MatchRepository {
	return &MatchRepository{db: db}
}

func (repository *MatchRepository) UpsertMatch(ctx context.Context, params ports.MatchUpsertParams) (ports.MatchUpsertResult, error) {
	result := ports.MatchUpsertResult{}
	var homeGoals sql.NullInt64
	var awayGoals sql.NullInt64

	err := repository.db.QueryRowContext(
		ctx,
		`WITH upserted AS (
			INSERT INTO matches (
				competition_id,
				season_id,
				match_date,
				home_team_id,
				away_team_id,
				home_goals,
				away_goals
			) VALUES ($1, $2, $3, $4, $5, $6, $7)
			ON CONFLICT ON CONSTRAINT matches_logical_identity_key
			DO UPDATE SET
				home_goals = EXCLUDED.home_goals,
				away_goals = EXCLUDED.away_goals
			WHERE matches.home_goals IS DISTINCT FROM EXCLUDED.home_goals
			   OR matches.away_goals IS DISTINCT FROM EXCLUDED.away_goals
			RETURNING id, competition_id, season_id, match_date, home_team_id, away_team_id, home_goals, away_goals, created_at, (xmax = 0) AS inserted, (xmax <> 0) AS updated
		)
		SELECT id, competition_id, season_id, match_date, home_team_id, away_team_id, home_goals, away_goals, created_at, inserted, updated
		FROM upserted
		UNION ALL
		SELECT id, competition_id, season_id, match_date, home_team_id, away_team_id, home_goals, away_goals, created_at, FALSE AS inserted, FALSE AS updated
		FROM matches
		WHERE competition_id = $1
		  AND season_id = $2
		  AND match_date = $3
		  AND home_team_id = $4
		  AND away_team_id = $5
		  AND NOT EXISTS (SELECT 1 FROM upserted)
		LIMIT 1`,
		params.CompetitionID,
		params.SeasonID,
		params.MatchDate,
		params.HomeTeamID,
		params.AwayTeamID,
		params.HomeGoals,
		params.AwayGoals,
	).Scan(
		&result.Match.ID,
		&result.Match.CompetitionID,
		&result.Match.SeasonID,
		&result.Match.MatchDate,
		&result.Match.HomeTeamID,
		&result.Match.AwayTeamID,
		&homeGoals,
		&awayGoals,
		&result.Match.CreatedAt,
		&result.Inserted,
		&result.Updated,
	)
	if err != nil {
		return ports.MatchUpsertResult{}, fmt.Errorf("upsert match: %w", err)
	}

	result.Match.HomeGoals = nullableInt64ToInt(homeGoals)
	result.Match.AwayGoals = nullableInt64ToInt(awayGoals)

	return result, nil
}

func (repository *MatchRepository) UpsertMatchOdds(ctx context.Context, params ports.MatchOddsUpsertParams) error {
	_, err := repository.db.ExecContext(
		ctx,
		`INSERT INTO match_odds (match_id, home_win, draw, away_win)
		 VALUES ($1, $2, $3, $4)
		 ON CONFLICT (match_id)
		 DO UPDATE SET
			home_win = EXCLUDED.home_win,
			draw = EXCLUDED.draw,
			away_win = EXCLUDED.away_win`,
		params.MatchID,
		params.HomeWin,
		params.Draw,
		params.AwayWin,
	)
	if err != nil {
		return fmt.Errorf("upsert match odds: %w", err)
	}

	return nil
}

type IngestionRunRepository struct {
	db *sql.DB
}

func NewIngestionRunRepository(db *sql.DB) *IngestionRunRepository {
	return &IngestionRunRepository{db: db}
}

func (repository *IngestionRunRepository) Create(ctx context.Context, source string) (domain.IngestionRun, error) {
	run := domain.IngestionRun{}
	err := repository.db.QueryRowContext(
		ctx,
		`INSERT INTO ingestion_runs (source, status)
		 VALUES ($1, $2)
		 RETURNING id, source, started_at, status, rows_processed, rows_inserted, rows_updated`,
		source,
		domain.IngestionRunStatusStarted,
	).Scan(
		&run.ID,
		&run.Source,
		&run.StartedAt,
		&run.Status,
		&run.RowsProcessed,
		&run.RowsInserted,
		&run.RowsUpdated,
	)
	if err != nil {
		return domain.IngestionRun{}, fmt.Errorf("create ingestion run: %w", err)
	}

	return run, nil
}

func (repository *IngestionRunRepository) MarkSuccess(ctx context.Context, runID int64, rowsProcessed int, rowsInserted int, rowsUpdated int) error {
	_, err := repository.db.ExecContext(
		ctx,
		`UPDATE ingestion_runs
		 SET finished_at = NOW(),
		     status = $2,
		     rows_processed = $3,
		     rows_inserted = $4,
		     rows_updated = $5,
		     error_message = NULL
		 WHERE id = $1`,
		runID,
		domain.IngestionRunStatusSuccess,
		rowsProcessed,
		rowsInserted,
		rowsUpdated,
	)
	if err != nil {
		return fmt.Errorf("mark ingestion run success: %w", err)
	}

	return nil
}

func (repository *IngestionRunRepository) MarkFailed(ctx context.Context, runID int64, rowsProcessed int, rowsInserted int, rowsUpdated int, errorMessage string) error {
	_, err := repository.db.ExecContext(
		ctx,
		`UPDATE ingestion_runs
		 SET finished_at = NOW(),
		     status = $2,
		     rows_processed = $3,
		     rows_inserted = $4,
		     rows_updated = $5,
		     error_message = $6
		 WHERE id = $1`,
		runID,
		domain.IngestionRunStatusFailed,
		rowsProcessed,
		rowsInserted,
		rowsUpdated,
		errorMessage,
	)
	if err != nil {
		return fmt.Errorf("mark ingestion run failed: %w", err)
	}

	return nil
}

func nullableInt64ToInt(value sql.NullInt64) *int {
	if !value.Valid {
		return nil
	}

	converted := int(value.Int64)
	return &converted
}
