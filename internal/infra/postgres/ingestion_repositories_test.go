package postgres

import (
	"context"
	"database/sql"
	"regexp"
	"testing"
	"time"

	"football-analytics/internal/domain"
	"football-analytics/internal/ports"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestCompetitionRepositoryCreateOrGetReturnsExisting(t *testing.T) {
	t.Parallel()

	database, mock := newMockDB(t)
	defer database.Close()

	createdAt := time.Date(2026, time.March, 13, 10, 0, 0, 0, time.UTC)
	rows := sqlmock.NewRows([]string{"id", "name", "country", "created_at"}).AddRow(int64(1), "Premier League", "England", createdAt)
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO competitions (name, country)
		 VALUES ($1, NULLIF($2, ''))
		 ON CONFLICT (name, (COALESCE(country, '')))
		 DO UPDATE SET name = EXCLUDED.name
		 RETURNING id, name, COALESCE(country, ''), created_at`)).WithArgs("Premier League", "England").WillReturnRows(rows)

	repository := NewCompetitionRepository(database)
	competition, err := repository.CreateOrGet(context.Background(), "Premier League", "England")
	if err != nil {
		t.Fatalf("CreateOrGet returned error: %v", err)
	}
	if competition.ID != 1 || competition.Name != "Premier League" || competition.Country != "England" {
		t.Fatalf("unexpected competition: %+v", competition)
	}

	assertNoMockErrors(t, mock)
}

func TestCompetitionRepositoryCreateOrGetWithEmptyCountry(t *testing.T) {
	t.Parallel()

	database, mock := newMockDB(t)
	defer database.Close()

	createdAt := time.Date(2026, time.March, 13, 10, 0, 0, 0, time.UTC)
	rows := sqlmock.NewRows([]string{"id", "name", "country", "created_at"}).AddRow(int64(2), "Premier League", "", createdAt)
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO competitions (name, country)
		 VALUES ($1, NULLIF($2, ''))
		 ON CONFLICT (name, (COALESCE(country, '')))
		 DO UPDATE SET name = EXCLUDED.name
		 RETURNING id, name, COALESCE(country, ''), created_at`)).WithArgs("Premier League", "").WillReturnRows(rows)

	repository := NewCompetitionRepository(database)
	competition, err := repository.CreateOrGet(context.Background(), "Premier League", "")
	if err != nil {
		t.Fatalf("CreateOrGet returned error: %v", err)
	}
	if competition.ID != 2 {
		t.Fatalf("expected competition id 2, got %d", competition.ID)
	}
	if competition.Country != "" {
		t.Fatalf("expected normalized empty country, got %q", competition.Country)
	}

	assertNoMockErrors(t, mock)
}

func TestSeasonRepositoryCreateOrGet(t *testing.T) {
	t.Parallel()

	database, mock := newMockDB(t)
	defer database.Close()

	createdAt := time.Date(2026, time.March, 13, 10, 0, 0, 0, time.UTC)
	rows := sqlmock.NewRows([]string{"id", "competition_id", "label", "created_at"}).AddRow(int64(3), int64(1), "2024-2025", createdAt)
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO seasons (competition_id, label)
		 VALUES ($1, $2)
		 ON CONFLICT ON CONSTRAINT seasons_competition_id_label_key
		 DO UPDATE SET label = EXCLUDED.label
		 RETURNING id, competition_id, label, created_at`)).WithArgs(int64(1), "2024-2025").WillReturnRows(rows)

	repository := NewSeasonRepository(database)
	season, err := repository.CreateOrGet(context.Background(), 1, "2024-2025")
	if err != nil {
		t.Fatalf("CreateOrGet returned error: %v", err)
	}
	if season.ID != 3 || season.Label != "2024-2025" {
		t.Fatalf("unexpected season: %+v", season)
	}

	assertNoMockErrors(t, mock)
}

func TestTeamRepositoryCreateOrGet(t *testing.T) {
	t.Parallel()

	database, mock := newMockDB(t)
	defer database.Close()

	createdAt := time.Date(2026, time.March, 13, 10, 0, 0, 0, time.UTC)
	rows := sqlmock.NewRows([]string{"id", "name", "created_at"}).AddRow(int64(4), "Arsenal", createdAt)
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO teams (name)
		 VALUES ($1)
		 ON CONFLICT (name)
		 DO UPDATE SET name = EXCLUDED.name
		 RETURNING id, name, created_at`)).WithArgs("Arsenal").WillReturnRows(rows)

	repository := NewTeamRepository(database)
	team, err := repository.CreateOrGet(context.Background(), "Arsenal")
	if err != nil {
		t.Fatalf("CreateOrGet returned error: %v", err)
	}
	if team.ID != 4 || team.Name != "Arsenal" {
		t.Fatalf("unexpected team: %+v", team)
	}

	assertNoMockErrors(t, mock)
}

func TestMatchRepositoryUpsertMatchInserted(t *testing.T) {
	t.Parallel()

	database, mock := newMockDB(t)
	defer database.Close()

	createdAt := time.Date(2026, time.March, 13, 10, 0, 0, 0, time.UTC)
	matchDate := time.Date(2024, time.August, 18, 0, 0, 0, 0, time.UTC)
	rows := sqlmock.NewRows([]string{"id", "competition_id", "season_id", "match_date", "home_team_id", "away_team_id", "home_goals", "away_goals", "created_at", "inserted", "updated"}).AddRow(int64(5), int64(1), int64(2), matchDate, int64(10), int64(20), int64(2), int64(1), createdAt, true, false)
	mock.ExpectQuery(regexp.QuoteMeta(`WITH upserted AS (
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
		LIMIT 1`)).WithArgs(int64(1), int64(2), matchDate, int64(10), int64(20), intPtrValue(2), intPtrValue(1)).WillReturnRows(rows)

	repository := NewMatchRepository(database)
	result, err := repository.UpsertMatch(context.Background(), ports.MatchUpsertParams{CompetitionID: 1, SeasonID: 2, MatchDate: matchDate, HomeTeamID: 10, AwayTeamID: 20, HomeGoals: intPtrValue(2), AwayGoals: intPtrValue(1)})
	if err != nil {
		t.Fatalf("UpsertMatch returned error: %v", err)
	}
	if !result.Inserted || result.Updated {
		t.Fatalf("unexpected upsert flags: %+v", result)
	}
	if result.Match.HomeGoals == nil || *result.Match.HomeGoals != 2 {
		t.Fatalf("expected home goals 2, got %+v", result.Match.HomeGoals)
	}

	assertNoMockErrors(t, mock)
}

func TestMatchRepositoryUpsertMatchOdds(t *testing.T) {
	t.Parallel()

	database, mock := newMockDB(t)
	defer database.Close()

	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO match_odds (match_id, home_win, draw, away_win)
		 VALUES ($1, $2, $3, $4)
		 ON CONFLICT (match_id)
		 DO UPDATE SET
			home_win = EXCLUDED.home_win,
			draw = EXCLUDED.draw,
			away_win = EXCLUDED.away_win`)).WithArgs(int64(5), floatPtrValue(1.8), floatPtrValue(3.5), floatPtrValue(4.2)).WillReturnResult(sqlmock.NewResult(1, 1))

	repository := NewMatchRepository(database)
	err := repository.UpsertMatchOdds(context.Background(), ports.MatchOddsUpsertParams{MatchID: 5, HomeWin: floatPtrValue(1.8), Draw: floatPtrValue(3.5), AwayWin: floatPtrValue(4.2)})
	if err != nil {
		t.Fatalf("UpsertMatchOdds returned error: %v", err)
	}

	assertNoMockErrors(t, mock)
}

func TestIngestionRunRepositoryLifecycle(t *testing.T) {
	t.Parallel()

	database, mock := newMockDB(t)
	defer database.Close()

	startedAt := time.Date(2026, time.March, 13, 10, 0, 0, 0, time.UTC)
	rows := sqlmock.NewRows([]string{"id", "source", "started_at", "status", "rows_processed", "rows_inserted", "rows_updated"}).AddRow(int64(7), "https://example.test/source.csv", startedAt, domain.IngestionRunStatusStarted, 0, 0, 0)
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO ingestion_runs (source, status)
		 VALUES ($1, $2)
		 RETURNING id, source, started_at, status, rows_processed, rows_inserted, rows_updated`)).WithArgs("https://example.test/source.csv", domain.IngestionRunStatusStarted).WillReturnRows(rows)
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE ingestion_runs
		 SET finished_at = NOW(),
		     status = $2,
		     rows_processed = $3,
		     rows_inserted = $4,
		     rows_updated = $5,
		     error_message = NULL
		 WHERE id = $1`)).WithArgs(int64(7), domain.IngestionRunStatusSuccess, 10, 8, 2).WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE ingestion_runs
		 SET finished_at = NOW(),
		     status = $2,
		     rows_processed = $3,
		     rows_inserted = $4,
		     rows_updated = $5,
		     error_message = $6
		 WHERE id = $1`)).WithArgs(int64(7), domain.IngestionRunStatusFailed, 10, 8, 2, "boom").WillReturnResult(sqlmock.NewResult(0, 1))

	repository := NewIngestionRunRepository(database)
	run, err := repository.Create(context.Background(), "https://example.test/source.csv")
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	if run.ID != 7 {
		t.Fatalf("unexpected run: %+v", run)
	}
	if err := repository.MarkSuccess(context.Background(), 7, 10, 8, 2); err != nil {
		t.Fatalf("MarkSuccess returned error: %v", err)
	}
	if err := repository.MarkFailed(context.Background(), 7, 10, 8, 2, "boom"); err != nil {
		t.Fatalf("MarkFailed returned error: %v", err)
	}

	assertNoMockErrors(t, mock)
}

func TestNullableInt64ToInt(t *testing.T) {
	t.Parallel()

	if value := nullableInt64ToInt(sql.NullInt64{}); value != nil {
		t.Fatalf("expected nil for invalid null int, got %v", *value)
	}

	value := nullableInt64ToInt(sql.NullInt64{Int64: 3, Valid: true})
	if value == nil || *value != 3 {
		t.Fatalf("expected 3, got %+v", value)
	}
}

func newMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	t.Helper()

	database, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("create sqlmock: %v", err)
	}

	return database, mock
}

func assertNoMockErrors(t *testing.T, mock sqlmock.Sqlmock) {
	t.Helper()

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet sql expectations: %v", err)
	}
}

func intPtrValue(value int) *int {
	return &value
}

func floatPtrValue(value float64) *float64 {
	return &value
}
