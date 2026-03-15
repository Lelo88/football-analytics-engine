package http

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"football-analytics/internal/domain"
	"football-analytics/internal/ports"
)

type AnalyticsUseCase interface {
	Teams(ctx context.Context) ([]domain.Team, error)
	TeamExists(ctx context.Context, teamID int64) (bool, error)
	TeamForm(ctx context.Context, filter ports.TeamQueryFilter) (domain.TeamForm, error)
	OverUnderSummary(ctx context.Context, filter ports.OverUnderQueryFilter) (domain.OverUnderSummary, error)
	SeasonSummaries(ctx context.Context, filter ports.SeasonSummaryFilter) ([]domain.SeasonTeamSummary, error)
}

type Handler struct {
	useCase AnalyticsUseCase
}

func NewHandler(useCase AnalyticsUseCase) http.Handler {
	handler := &Handler{useCase: useCase}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /", handler.redirectToUI)
	mux.HandleFunc("GET /ui", handler.getUIHome)
	mux.HandleFunc("GET /ui/team-summary", handler.getUITeamSummary)
	mux.HandleFunc("GET /teams", handler.getTeams)
	mux.HandleFunc("GET /teams/{id}/form", handler.getTeamForm)
	mux.HandleFunc("GET /teams/{id}/overunder", handler.getTeamOverUnder)
	mux.HandleFunc("GET /teams/{id}/season-summary", handler.getTeamSeasonSummary)

	return mux
}

func (handler *Handler) getTeams(writer http.ResponseWriter, request *http.Request) {
	teams, err := handler.useCase.Teams(request.Context())
	if err != nil {
		writeInternalServerError(writer, err)
		return
	}

	type teamResponse struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	}

	response := make([]teamResponse, 0, len(teams))
	for _, team := range teams {
		response = append(response, teamResponse{ID: team.ID, Name: team.Name})
	}

	writeJSON(writer, http.StatusOK, map[string]any{"data": response})
}

func (handler *Handler) getTeamForm(writer http.ResponseWriter, request *http.Request) {
	teamID, err := parseTeamID(request)
	if err != nil {
		writeBadRequest(writer, err.Error())
		return
	}

	filter, err := buildTeamQueryFilter(request, teamID)
	if err != nil {
		writeBadRequest(writer, err.Error())
		return
	}

	if err = handler.ensureTeamExists(request.Context(), teamID); err != nil {
		handleDomainError(writer, err)
		return
	}

	form, err := handler.useCase.TeamForm(request.Context(), filter)
	if err != nil {
		handleDomainError(writer, err)
		return
	}

	writeJSON(writer, http.StatusOK, map[string]any{"data": form})
}

func (handler *Handler) getTeamOverUnder(writer http.ResponseWriter, request *http.Request) {
	teamID, err := parseTeamID(request)
	if err != nil {
		writeBadRequest(writer, err.Error())
		return
	}

	filter, err := buildOverUnderFilter(request, teamID)
	if err != nil {
		writeBadRequest(writer, err.Error())
		return
	}

	if err = handler.ensureTeamExists(request.Context(), teamID); err != nil {
		handleDomainError(writer, err)
		return
	}

	summary, err := handler.useCase.OverUnderSummary(request.Context(), filter)
	if err != nil {
		handleDomainError(writer, err)
		return
	}

	writeJSON(writer, http.StatusOK, map[string]any{"data": summary})
}

func (handler *Handler) getTeamSeasonSummary(writer http.ResponseWriter, request *http.Request) {
	teamID, err := parseTeamID(request)
	if err != nil {
		writeBadRequest(writer, err.Error())
		return
	}

	filter, err := buildSeasonSummaryFilter(request, teamID)
	if err != nil {
		writeBadRequest(writer, err.Error())
		return
	}

	if err = handler.ensureTeamExists(request.Context(), teamID); err != nil {
		handleDomainError(writer, err)
		return
	}

	summaries, err := handler.useCase.SeasonSummaries(request.Context(), filter)
	if err != nil {
		handleDomainError(writer, err)
		return
	}

	writeJSON(writer, http.StatusOK, map[string]any{"data": summaries})
}

func (handler *Handler) ensureTeamExists(ctx context.Context, teamID int64) error {
	exists, err := handler.useCase.TeamExists(ctx, teamID)
	if err != nil {
		return err
	}
	if !exists {
		return errTeamNotFound
	}

	return nil
}

var errTeamNotFound = errors.New("team not found")

func parseTeamID(request *http.Request) (int64, error) {
	teamIDValue := request.PathValue("id")
	if strings.TrimSpace(teamIDValue) == "" {
		return 0, fmt.Errorf("team id is required")
	}

	teamID, err := strconv.ParseInt(teamIDValue, 10, 64)
	if err != nil || teamID <= 0 {
		return 0, fmt.Errorf("team id must be a positive integer")
	}

	return teamID, nil
}

func buildTeamQueryFilter(request *http.Request, teamID int64) (ports.TeamQueryFilter, error) {
	venue, err := parseVenue(request.URL.Query().Get("venue"))
	if err != nil {
		return ports.TeamQueryFilter{}, err
	}

	lastN, err := parseOptionalInt(request.URL.Query().Get("last_n"), 0)
	if err != nil {
		return ports.TeamQueryFilter{}, fmt.Errorf("invalid last_n: %w", err)
	}

	return ports.TeamQueryFilter{
		TeamID:      teamID,
		SeasonLabel: strings.TrimSpace(request.URL.Query().Get("season")),
		LastN:       lastN,
		Venue:       venue,
	}, nil
}

func buildOverUnderFilter(request *http.Request, teamID int64) (ports.OverUnderQueryFilter, error) {
	venue, err := parseVenue(request.URL.Query().Get("venue"))
	if err != nil {
		return ports.OverUnderQueryFilter{}, err
	}

	lastN, err := parseOptionalInt(request.URL.Query().Get("last_n"), 0)
	if err != nil {
		return ports.OverUnderQueryFilter{}, fmt.Errorf("invalid last_n: %w", err)
	}

	threshold, err := parseOptionalFloat(request.URL.Query().Get("threshold"), 2.5)
	if err != nil {
		return ports.OverUnderQueryFilter{}, fmt.Errorf("invalid threshold: %w", err)
	}

	return ports.OverUnderQueryFilter{
		TeamID:      teamID,
		SeasonLabel: strings.TrimSpace(request.URL.Query().Get("season")),
		LastN:       lastN,
		Venue:       venue,
		Threshold:   threshold,
	}, nil
}

func buildSeasonSummaryFilter(request *http.Request, teamID int64) (ports.SeasonSummaryFilter, error) {
	venue, err := parseVenue(request.URL.Query().Get("venue"))
	if err != nil {
		return ports.SeasonSummaryFilter{}, err
	}

	return ports.SeasonSummaryFilter{
		TeamID:      teamID,
		SeasonLabel: strings.TrimSpace(request.URL.Query().Get("season")),
		Venue:       venue,
	}, nil
}

func parseVenue(raw string) (domain.MatchVenue, error) {
	value := strings.TrimSpace(strings.ToLower(raw))
	if value == "" {
		return domain.MatchVenueAll, nil
	}

	venue := domain.MatchVenue(value)
	switch venue {
	case domain.MatchVenueAll, domain.MatchVenueHome, domain.MatchVenueAway:
		return venue, nil
	default:
		return "", fmt.Errorf("invalid venue, expected one of: all, home, away")
	}
}

func parseOptionalInt(raw string, fallback int) (int, error) {
	value := strings.TrimSpace(raw)
	if value == "" {
		return fallback, nil
	}

	parsed, err := strconv.Atoi(value)
	if err != nil || parsed < 0 {
		return 0, fmt.Errorf("must be a non-negative integer")
	}

	return parsed, nil
}

func parseOptionalFloat(raw string, fallback float64) (float64, error) {
	value := strings.TrimSpace(raw)
	if value == "" {
		return fallback, nil
	}

	parsed, err := strconv.ParseFloat(value, 64)
	if err != nil || parsed < 0 {
		return 0, fmt.Errorf("must be a non-negative number")
	}

	return parsed, nil
}

func handleDomainError(writer http.ResponseWriter, err error) {
	if errors.Is(err, errTeamNotFound) {
		writeError(writer, http.StatusNotFound, "not_found", err.Error())
		return
	}

	if strings.Contains(err.Error(), "must be") || strings.Contains(err.Error(), "invalid") {
		writeBadRequest(writer, err.Error())
		return
	}

	writeInternalServerError(writer, err)
}

func writeBadRequest(writer http.ResponseWriter, message string) {
	writeError(writer, http.StatusBadRequest, "bad_request", message)
}

func writeInternalServerError(writer http.ResponseWriter, err error) {
	writeError(writer, http.StatusInternalServerError, "internal_error", err.Error())
}

func writeError(writer http.ResponseWriter, statusCode int, code string, message string) {
	writeJSON(writer, statusCode, map[string]any{
		"error": map[string]string{
			"code":    code,
			"message": message,
		},
	})
}

func writeJSON(writer http.ResponseWriter, statusCode int, payload any) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(statusCode)
	_ = json.NewEncoder(writer).Encode(payload)
}
