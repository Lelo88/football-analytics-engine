package http

import (
	"context"
	"encoding/json"
	"fmt"
	stdhttp "net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"football-analytics/internal/domain"
	"football-analytics/internal/ports"
)

func TestGetTeams(t *testing.T) {
	t.Parallel()

	handler := NewHandler(&fakeAnalyticsUseCase{
		teams: []domain.Team{{ID: 1, Name: "Arsenal"}, {ID: 2, Name: "Chelsea"}},
	})

	request := httptest.NewRequest(stdhttp.MethodGet, "/teams", nil)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)

	if response.Code != stdhttp.StatusOK {
		t.Fatalf("expected status 200, got %d", response.Code)
	}

	var payload map[string][]map[string]any
	if err := json.Unmarshal(response.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if len(payload["data"]) != 2 {
		t.Fatalf("expected 2 teams, got %+v", payload)
	}
}

func TestGetTeamFormInvalidTeamID(t *testing.T) {
	t.Parallel()

	handler := NewHandler(&fakeAnalyticsUseCase{})

	request := httptest.NewRequest(stdhttp.MethodGet, "/teams/abc/form", nil)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)

	if response.Code != stdhttp.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", response.Code)
	}
}

func TestGetTeamFormNotFound(t *testing.T) {
	t.Parallel()

	handler := NewHandler(&fakeAnalyticsUseCase{teamExists: false})

	request := httptest.NewRequest(stdhttp.MethodGet, "/teams/99/form", nil)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)

	if response.Code != stdhttp.StatusNotFound {
		t.Fatalf("expected status 404, got %d", response.Code)
	}
}

func TestGetTeamFormInvalidVenue(t *testing.T) {
	t.Parallel()

	handler := NewHandler(&fakeAnalyticsUseCase{teamExists: true})

	request := httptest.NewRequest(stdhttp.MethodGet, "/teams/1/form?venue=invalid", nil)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)

	if response.Code != stdhttp.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", response.Code)
	}
}

func TestGetTeamOverUnderInvalidThreshold(t *testing.T) {
	t.Parallel()

	handler := NewHandler(&fakeAnalyticsUseCase{teamExists: true})

	request := httptest.NewRequest(stdhttp.MethodGet, "/teams/1/overunder?threshold=abc", nil)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)

	if response.Code != stdhttp.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", response.Code)
	}
}

func TestGetTeamFormUnexpectedFailure(t *testing.T) {
	t.Parallel()

	handler := NewHandler(&fakeAnalyticsUseCase{teamExists: true, teamFormErr: fmt.Errorf("db down")})

	request := httptest.NewRequest(stdhttp.MethodGet, "/teams/1/form", nil)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)

	if response.Code != stdhttp.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", response.Code)
	}
}

func TestGetUIHome(t *testing.T) {
	t.Parallel()

	handler := NewHandler(&fakeAnalyticsUseCase{
		teams: []domain.Team{{ID: 1, Name: "Arsenal"}},
	})

	request := httptest.NewRequest(stdhttp.MethodGet, "/ui", nil)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)

	if response.Code != stdhttp.StatusOK {
		t.Fatalf("expected status 200, got %d", response.Code)
	}

	if !strings.Contains(response.Body.String(), "Football Analytics UI Lite") {
		t.Fatalf("expected UI title in response body")
	}
}

func TestGetUITeamSummary(t *testing.T) {
	t.Parallel()

	handler := NewHandler(&fakeAnalyticsUseCase{
		teamExists: true,
		teamForm: domain.TeamForm{
			MatchesPlayed: 5,
			Wins:          3,
			Draws:         1,
			Losses:        1,
			Points:        10,
			GoalsFor:      9,
			GoalsAgainst:  4,
		},
		overUnder: domain.OverUnderSummary{
			MatchesPlayed:   5,
			Threshold:       2.5,
			OverCount:       3,
			UnderEqualCount: 2,
		},
		seasonSummaries: []domain.SeasonTeamSummary{{
			SeasonLabel:   "2024-2025",
			MatchesPlayed: 38,
			Points:        72,
			GoalsFor:      65,
			GoalsAgainst:  40,
		}},
	})

	request := httptest.NewRequest(stdhttp.MethodGet, "/ui/team-summary?team_id=1&last_n=5&venue=all&threshold=2.5", nil)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)

	if response.Code != stdhttp.StatusOK {
		t.Fatalf("expected status 200, got %d", response.Code)
	}

	if !strings.Contains(response.Body.String(), "Resumen del equipo #1") {
		t.Fatalf("expected team summary title in response body")
	}

	if !strings.Contains(response.Body.String(), "season-summary-chart") {
		t.Fatalf("expected chart canvas in response body")
	}

	if !strings.Contains(response.Body.String(), "new Chart") {
		t.Fatalf("expected chart script in response body")
	}
}

func TestGetUITeamSummaryInvalidTeamID(t *testing.T) {
	t.Parallel()

	handler := NewHandler(&fakeAnalyticsUseCase{})

	request := httptest.NewRequest(stdhttp.MethodGet, "/ui/team-summary?team_id=abc", nil)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)

	if response.Code != stdhttp.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", response.Code)
	}

	if !strings.Contains(response.Body.String(), "team_id") {
		t.Fatalf("expected validation message about team_id")
	}
}

type fakeAnalyticsUseCase struct {
	teams              []domain.Team
	teamsErr           error
	teamExists         bool
	teamExistsErr      error
	teamForm           domain.TeamForm
	teamFormErr        error
	overUnder          domain.OverUnderSummary
	overUnderErr       error
	seasonSummaries    []domain.SeasonTeamSummary
	seasonSummariesErr error
}

func (fakeUseCase *fakeAnalyticsUseCase) Teams(ctx context.Context) ([]domain.Team, error) {
	if fakeUseCase.teamsErr != nil {
		return nil, fakeUseCase.teamsErr
	}
	return fakeUseCase.teams, nil
}

func (fakeUseCase *fakeAnalyticsUseCase) TeamExists(ctx context.Context, teamID int64) (bool, error) {
	if fakeUseCase.teamExistsErr != nil {
		return false, fakeUseCase.teamExistsErr
	}
	return fakeUseCase.teamExists, nil
}

func (fakeUseCase *fakeAnalyticsUseCase) TeamForm(ctx context.Context, filter ports.TeamQueryFilter) (domain.TeamForm, error) {
	if fakeUseCase.teamFormErr != nil {
		return domain.TeamForm{}, fakeUseCase.teamFormErr
	}
	return fakeUseCase.teamForm, nil
}

func (fakeUseCase *fakeAnalyticsUseCase) OverUnderSummary(ctx context.Context, filter ports.OverUnderQueryFilter) (domain.OverUnderSummary, error) {
	if fakeUseCase.overUnderErr != nil {
		return domain.OverUnderSummary{}, fakeUseCase.overUnderErr
	}
	return fakeUseCase.overUnder, nil
}

func (fakeUseCase *fakeAnalyticsUseCase) SeasonSummaries(ctx context.Context, filter ports.SeasonSummaryFilter) ([]domain.SeasonTeamSummary, error) {
	if fakeUseCase.seasonSummariesErr != nil {
		return nil, fakeUseCase.seasonSummariesErr
	}
	return fakeUseCase.seasonSummaries, nil
}
