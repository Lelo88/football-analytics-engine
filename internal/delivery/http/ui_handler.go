package http

import (
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"football-analytics/internal/domain"
	"football-analytics/internal/ports"
)

//go:embed templates/*.gohtml
var uiTemplateFS embed.FS

var uiTemplates = template.Must(template.ParseFS(uiTemplateFS, "templates/*.gohtml"))

type uiHomeData struct {
	Teams []teamItem
}

type teamItem struct {
	ID   int64
	Name string
}

type uiTeamSummaryData struct {
	TeamID       int64
	SeasonLabel  string
	LastN        int
	Venue        string
	Threshold    float64
	Form         domain.TeamForm
	OverUnder    domain.OverUnderSummary
	Summaries    []domain.SeasonTeamSummary
	SeasonLabels template.JS
	SeasonPoints template.JS
	SeasonGF     template.JS
	SeasonGA     template.JS
}

func (handler *Handler) redirectToUI(writer http.ResponseWriter, request *http.Request) {
	http.Redirect(writer, request, "/ui", http.StatusFound)
}

func (handler *Handler) getUIHome(writer http.ResponseWriter, request *http.Request) {
	teams, err := handler.useCase.Teams(request.Context())
	if err != nil {
		writeInternalServerError(writer, err)
		return
	}

	viewTeams := make([]teamItem, 0, len(teams))
	for _, team := range teams {
		viewTeams = append(viewTeams, teamItem{ID: team.ID, Name: team.Name})
	}

	writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err = uiTemplates.ExecuteTemplate(writer, "base", uiHomeData{Teams: viewTeams}); err != nil {
		writeInternalServerError(writer, err)
		return
	}
}

func (handler *Handler) getUITeamSummary(writer http.ResponseWriter, request *http.Request) {
	teamID, err := parseRequiredPositiveInt(request.URL.Query().Get("team_id"), "team_id")
	if err != nil {
		handler.renderUIError(writer, http.StatusBadRequest, err.Error())
		return
	}

	lastN, err := parseOptionalInt(request.URL.Query().Get("last_n"), 5)
	if err != nil {
		handler.renderUIError(writer, http.StatusBadRequest, fmt.Sprintf("invalid last_n: %v", err))
		return
	}
	if lastN == 0 {
		lastN = 5
	}

	venue, err := parseVenue(request.URL.Query().Get("venue"))
	if err != nil {
		handler.renderUIError(writer, http.StatusBadRequest, err.Error())
		return
	}

	threshold, err := parseOptionalFloat(request.URL.Query().Get("threshold"), 2.5)
	if err != nil {
		handler.renderUIError(writer, http.StatusBadRequest, fmt.Sprintf("invalid threshold: %v", err))
		return
	}

	seasonLabel := strings.TrimSpace(request.URL.Query().Get("season"))

	if err = handler.ensureTeamExists(request.Context(), teamID); err != nil {
		if err == errTeamNotFound {
			handler.renderUIError(writer, http.StatusNotFound, err.Error())
			return
		}
		handler.renderUIError(writer, http.StatusInternalServerError, err.Error())
		return
	}

	teamFilter := ports.TeamQueryFilter{
		TeamID:      teamID,
		SeasonLabel: seasonLabel,
		LastN:       lastN,
		Venue:       venue,
	}

	overUnderFilter := ports.OverUnderQueryFilter{
		TeamID:      teamID,
		SeasonLabel: seasonLabel,
		LastN:       lastN,
		Venue:       venue,
		Threshold:   threshold,
	}

	seasonFilter := ports.SeasonSummaryFilter{
		TeamID:      teamID,
		SeasonLabel: seasonLabel,
		Venue:       venue,
	}

	form, err := handler.useCase.TeamForm(request.Context(), teamFilter)
	if err != nil {
		handler.renderUIError(writer, http.StatusInternalServerError, err.Error())
		return
	}

	overUnder, err := handler.useCase.OverUnderSummary(request.Context(), overUnderFilter)
	if err != nil {
		handler.renderUIError(writer, http.StatusInternalServerError, err.Error())
		return
	}

	summaries, err := handler.useCase.SeasonSummaries(request.Context(), seasonFilter)
	if err != nil {
		handler.renderUIError(writer, http.StatusInternalServerError, err.Error())
		return
	}

	writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	seasonLabels := make([]string, 0, len(summaries))
	seasonPoints := make([]int, 0, len(summaries))
	seasonGF := make([]int, 0, len(summaries))
	seasonGA := make([]int, 0, len(summaries))
	for _, summary := range summaries {
		seasonLabels = append(seasonLabels, summary.SeasonLabel)
		seasonPoints = append(seasonPoints, summary.Points)
		seasonGF = append(seasonGF, summary.GoalsFor)
		seasonGA = append(seasonGA, summary.GoalsAgainst)
	}

	err = uiTemplates.ExecuteTemplate(writer, "team_summary", uiTeamSummaryData{
		TeamID:       teamID,
		SeasonLabel:  seasonLabel,
		LastN:        lastN,
		Venue:        string(venue),
		Threshold:    threshold,
		Form:         form,
		OverUnder:    overUnder,
		Summaries:    summaries,
		SeasonLabels: marshalForJS(seasonLabels),
		SeasonPoints: marshalForJS(seasonPoints),
		SeasonGF:     marshalForJS(seasonGF),
		SeasonGA:     marshalForJS(seasonGA),
	})
	if err != nil {
		handler.renderUIError(writer, http.StatusInternalServerError, err.Error())
		return
	}
}

func (handler *Handler) renderUIError(writer http.ResponseWriter, statusCode int, message string) {
	writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	writer.WriteHeader(statusCode)
	_ = uiTemplates.ExecuteTemplate(writer, "ui_error", message)
}

func parseRequiredPositiveInt(raw string, fieldName string) (int64, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return 0, fmt.Errorf("%s is required", fieldName)
	}

	value, err := strconv.ParseInt(trimmed, 10, 64)
	if err != nil || value <= 0 {
		return 0, fmt.Errorf("%s must be a positive integer", fieldName)
	}

	return value, nil
}

func marshalForJS(value any) template.JS {
	encoded, err := json.Marshal(value)
	if err != nil {
		return template.JS("[]")
	}

	return template.JS(encoded)
}
