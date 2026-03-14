package postgres

import (
	"context"
	"regexp"
	"strings"
	"testing"
	"time"

	"football-analytics/internal/domain"
	"football-analytics/internal/ports"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestTeamAnalyticsReadRepositoryListTeams(t *testing.T) {
	t.Parallel()

	database, mock := newMockDB(t)
	defer database.Close()

	createdAt := time.Date(2026, time.March, 14, 10, 0, 0, 0, time.UTC)
	rows := sqlmock.NewRows([]string{"id", "name", "created_at"}).
		AddRow(int64(1), "Arsenal", createdAt).
		AddRow(int64(2), "Chelsea", createdAt)

	mock.ExpectQuery(queryPattern("SELECT id, name, created_at", "FROM teams", "ORDER BY name")).WillReturnRows(rows)

	repository := NewTeamAnalyticsReadRepository(database)
	teams, err := repository.ListTeams(context.Background())
	if err != nil {
		t.Fatalf("ListTeams returned error: %v", err)
	}

	if len(teams) != 2 || teams[0].Name != "Arsenal" || teams[1].ID != 2 {
		t.Fatalf("unexpected teams: %+v", teams)
	}

	assertNoMockErrors(t, mock)
}

func TestTeamAnalyticsReadRepositoryTeamExists(t *testing.T) {
	t.Parallel()

	database, mock := newMockDB(t)
	defer database.Close()

	rows := sqlmock.NewRows([]string{"exists"}).AddRow(true)
	mock.ExpectQuery(queryPattern("SELECT EXISTS(", "FROM teams", "WHERE id = $1", ")")).
		WithArgs(int64(10)).
		WillReturnRows(rows)

	repository := NewTeamAnalyticsReadRepository(database)
	exists, err := repository.TeamExists(context.Background(), 10)
	if err != nil {
		t.Fatalf("TeamExists returned error: %v", err)
	}
	if !exists {
		t.Fatal("expected team to exist")
	}

	assertNoMockErrors(t, mock)
}

func TestTeamAnalyticsReadRepositoryGetTeamForm(t *testing.T) {
	t.Parallel()

	database, mock := newMockDB(t)
	defer database.Close()

	rows := sqlmock.NewRows([]string{"matches_played", "wins", "draws", "losses", "points", "goals_for", "goals_against"}).
		AddRow(5, 3, 1, 1, 10, 9, 4)

	mock.ExpectQuery(queryPattern(
		"WITH filtered_matches AS (",
		"AND (m.home_team_id = $1 OR m.away_team_id = $1)",
		"SELECT",
		"COUNT(*) AS matches_played",
		"AS wins",
		"AS draws",
		"AS losses",
		"AS points",
		"AS goals_for",
		"AS goals_against",
		"FROM filtered_matches",
	)).WithArgs(int64(10)).WillReturnRows(rows)

	repository := NewTeamAnalyticsReadRepository(database)
	result, err := repository.GetTeamForm(context.Background(), ports.TeamQueryFilter{TeamID: 10, Venue: domain.MatchVenueAll})
	if err != nil {
		t.Fatalf("GetTeamForm returned error: %v", err)
	}

	if result.Points != 10 || result.Wins != 3 || result.MatchesPlayed != 5 {
		t.Fatalf("unexpected team form: %+v", result)
	}

	assertNoMockErrors(t, mock)
}

func TestTeamAnalyticsReadRepositoryGetGoalsSummaryHomeFilterAndLastN(t *testing.T) {
	t.Parallel()

	database, mock := newMockDB(t)
	defer database.Close()

	rows := sqlmock.NewRows([]string{"matches_played", "goals_for", "goals_against"}).
		AddRow(3, 7, 2)

	mock.ExpectQuery(queryPattern(
		"WITH filtered_matches AS (",
		"AND m.home_team_id = $1",
		"ORDER BY m.match_date DESC",
		"LIMIT 3",
		"COUNT(*) AS matches_played",
		"SUM(goals_for)",
		"SUM(goals_against)",
		"FROM filtered_matches",
	)).WithArgs(int64(8)).WillReturnRows(rows)

	repository := NewTeamAnalyticsReadRepository(database)
	result, err := repository.GetGoalsSummary(context.Background(), ports.TeamQueryFilter{TeamID: 8, Venue: domain.MatchVenueHome, LastN: 3})
	if err != nil {
		t.Fatalf("GetGoalsSummary returned error: %v", err)
	}

	if result.MatchesPlayed != 3 || result.GoalsFor != 7 || result.GoalsAgainst != 2 {
		t.Fatalf("unexpected goals summary: %+v", result)
	}

	assertNoMockErrors(t, mock)
}

func TestTeamAnalyticsReadRepositoryGetOverUnderSummary(t *testing.T) {
	t.Parallel()

	database, mock := newMockDB(t)
	defer database.Close()

	rows := sqlmock.NewRows([]string{"matches_played", "over_count", "under_equal_count"}).
		AddRow(4, 3, 1)

	mock.ExpectQuery(queryPattern(
		"WITH filtered_matches AS (",
		"AND m.away_team_id = $1",
		"AND s.label = $2",
		"(goals_for + goals_against) > $3",
		"AS over_count",
		"(goals_for + goals_against) <= $3",
		"AS under_equal_count",
		"FROM filtered_matches",
	)).WithArgs(int64(12), "2024-2025", 2.5).WillReturnRows(rows)

	repository := NewTeamAnalyticsReadRepository(database)
	result, err := repository.GetOverUnderSummary(context.Background(), ports.OverUnderQueryFilter{
		TeamID:      12,
		Venue:       domain.MatchVenueAway,
		SeasonLabel: "2024-2025",
		Threshold:   2.5,
	})
	if err != nil {
		t.Fatalf("GetOverUnderSummary returned error: %v", err)
	}

	if result.MatchesPlayed != 4 || result.OverCount != 3 || result.UnderEqualCount != 1 || result.Threshold != 2.5 {
		t.Fatalf("unexpected over/under summary: %+v", result)
	}

	assertNoMockErrors(t, mock)
}

func TestTeamAnalyticsReadRepositoryGetSeasonSummaries(t *testing.T) {
	t.Parallel()

	database, mock := newMockDB(t)
	defer database.Close()

	rows := sqlmock.NewRows([]string{"season_label", "matches_played", "points", "goals_for", "goals_against"}).
		AddRow("2023-2024", 38, 76, 81, 32).
		AddRow("2024-2025", 20, 44, 46, 19)

	mock.ExpectQuery(queryPattern(
		"WITH filtered_matches AS (",
		"AND (m.home_team_id = $1 OR m.away_team_id = $1)",
		"SELECT",
		"season_label",
		"COUNT(*) AS matches_played",
		"AS points",
		"SUM(goals_for)",
		"SUM(goals_against)",
		"GROUP BY season_label",
		"ORDER BY season_label",
	)).WithArgs(int64(6)).WillReturnRows(rows)

	repository := NewTeamAnalyticsReadRepository(database)
	result, err := repository.GetSeasonSummaries(context.Background(), ports.SeasonSummaryFilter{TeamID: 6, Venue: domain.MatchVenueAll})
	if err != nil {
		t.Fatalf("GetSeasonSummaries returned error: %v", err)
	}

	if len(result) != 2 {
		t.Fatalf("expected 2 season summaries, got %d", len(result))
	}
	if result[0].SeasonLabel != "2023-2024" || result[1].Points != 44 {
		t.Fatalf("unexpected season summaries: %+v", result)
	}

	assertNoMockErrors(t, mock)
}

func TestBuildTeamMatchesBaseQueryVenueModes(t *testing.T) {
	t.Parallel()

	homeQuery, _ := buildTeamMatchesBaseQuery(1, domain.MatchVenueHome, "")
	if !strings.Contains(homeQuery, "AND m.home_team_id = $1") {
		t.Fatalf("expected home-only filter in query: %s", homeQuery)
	}

	awayQuery, _ := buildTeamMatchesBaseQuery(1, domain.MatchVenueAway, "")
	if !strings.Contains(awayQuery, "AND m.away_team_id = $1") {
		t.Fatalf("expected away-only filter in query: %s", awayQuery)
	}

	allQuery, _ := buildTeamMatchesBaseQuery(1, domain.MatchVenueAll, "")
	if !strings.Contains(allQuery, "AND (m.home_team_id = $1 OR m.away_team_id = $1)") {
		t.Fatalf("expected all-matches filter in query: %s", allQuery)
	}
}

func queryPattern(parts ...string) string {
	escapedParts := make([]string, 0, len(parts))
	for _, part := range parts {
		escapedParts = append(escapedParts, regexp.QuoteMeta(part))
	}

	return "(?s)" + strings.Join(escapedParts, ".*")
}
