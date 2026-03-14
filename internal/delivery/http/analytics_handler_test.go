package http

import (
	"context"
	"encoding/json"
	"fmt"
	stdhttp "net/http"
	"net/http/httptest"
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
