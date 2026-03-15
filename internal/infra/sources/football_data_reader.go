package sources

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"football-analytics/internal/domain"
	"football-analytics/internal/ports"
)

type FootballDataReader struct {
	client    *http.Client
	sourceURL string
}

func NewFootballDataReader(client *http.Client) *FootballDataReader {
	if client == nil {
		client = http.DefaultClient
	}

	return &FootballDataReader{client: client}
}

func NewFootballDataSource(client *http.Client, sourceURL string) *FootballDataReader {
	reader := NewFootballDataReader(client)
	reader.sourceURL = sourceURL
	return reader
}

func (reader *FootballDataReader) Fetch(ctx context.Context) ([]domain.IngestionMatch, error) {
	if strings.TrimSpace(reader.sourceURL) == "" {
		return nil, fmt.Errorf("source url is required")
	}

	rows, err := reader.ReadMatches(ctx, reader.sourceURL)
	if err != nil {
		return nil, err
	}

	matches := make([]domain.IngestionMatch, 0, len(rows))
	for _, row := range rows {
		matches = append(matches, domain.IngestionMatch{
			CompetitionName: row.CompetitionName,
			Country:         row.Country,
			SeasonLabel:     row.SeasonLabel,
			MatchDate:       row.MatchDate,
			HomeTeamName:    row.HomeTeamName,
			AwayTeamName:    row.AwayTeamName,
			HomeGoals:       row.HomeGoals,
			AwayGoals:       row.AwayGoals,
			HomeWinOdds:     row.HomeWinOdds,
			DrawOdds:        row.DrawOdds,
			AwayWinOdds:     row.AwayWinOdds,
		})
	}

	return matches, nil
}

func (reader *FootballDataReader) ReadMatches(ctx context.Context, sourceURL string) ([]ports.SourceMatchRow, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, sourceURL, nil)
	if err != nil {
		return nil, fmt.Errorf("create source request: %w", err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36")

	resp, err := reader.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch source csv: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected source status: %s", resp.Status)
	}

	csvReader := csv.NewReader(resp.Body)
	csvReader.FieldsPerRecord = -1
	csvReader.ReuseRecord = false

	headers, err := csvReader.Read()
	if err != nil {
		if err == io.EOF {
			return []ports.SourceMatchRow{}, nil
		}

		return nil, fmt.Errorf("read csv header: %w", err)
	}

	indexByHeader, err := buildHeaderIndex(headers)
	if err != nil {
		return nil, err
	}

	rows := make([]ports.SourceMatchRow, 0)
	lineNumber := 1
	for {
		record, err := csvReader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}

			return nil, fmt.Errorf("read csv row %d: %w", lineNumber+1, err)
		}

		lineNumber++
		row, err := parseSourceMatchRow(record, indexByHeader)
		if err != nil {
			return nil, fmt.Errorf("parse csv row %d: %w", lineNumber, err)
		}

		rows = append(rows, row)
	}

	return rows, nil
}

func buildHeaderIndex(headers []string) (map[string]int, error) {
	indexByHeader := make(map[string]int, len(headers))
	for index, header := range headers {
		normalized := normalizeHeaderName(header)
		if normalized == "" {
			continue
		}

		indexByHeader[normalized] = index
	}

	requiredHeaders := []string{"Div", "Date", "HomeTeam", "AwayTeam"}
	for _, header := range requiredHeaders {
		if _, ok := indexByHeader[header]; !ok {
			return nil, fmt.Errorf("missing required header %q", header)
		}
	}

	return indexByHeader, nil
}

func normalizeHeaderName(header string) string {
	trimmed := strings.TrimSpace(header)
	trimmed = strings.TrimPrefix(trimmed, "\ufeff")
	return strings.TrimSpace(trimmed)
}

func parseSourceMatchRow(record []string, indexByHeader map[string]int) (ports.SourceMatchRow, error) {
	competitionCode := normalizeCompetitionCode(fieldValue(record, indexByHeader, "Div"))
	if competitionCode == "" {
		return ports.SourceMatchRow{}, fmt.Errorf("missing competition code")
	}

	competitionName, country := normalizeCompetition(competitionCode)
	matchDate, err := parseMatchDate(fieldValue(record, indexByHeader, "Date"))
	if err != nil {
		return ports.SourceMatchRow{}, err
	}

	homeTeamName := normalizeTeamName(fieldValue(record, indexByHeader, "HomeTeam"))
	if homeTeamName == "" {
		return ports.SourceMatchRow{}, fmt.Errorf("missing home team")
	}

	awayTeamName := normalizeTeamName(fieldValue(record, indexByHeader, "AwayTeam"))
	if awayTeamName == "" {
		return ports.SourceMatchRow{}, fmt.Errorf("missing away team")
	}

	homeGoals, err := parseOptionalInt(fieldValue(record, indexByHeader, "FTHG"))
	if err != nil {
		return ports.SourceMatchRow{}, fmt.Errorf("parse FTHG: %w", err)
	}

	awayGoals, err := parseOptionalInt(fieldValue(record, indexByHeader, "FTAG"))
	if err != nil {
		return ports.SourceMatchRow{}, fmt.Errorf("parse FTAG: %w", err)
	}

	homeWinOdds, err := parseOptionalFloat(firstAvailableValue(record, indexByHeader, []string{"B365H", "AvgH", "PSH"}))
	if err != nil {
		return ports.SourceMatchRow{}, fmt.Errorf("parse home odds: %w", err)
	}

	drawOdds, err := parseOptionalFloat(firstAvailableValue(record, indexByHeader, []string{"B365D", "AvgD", "PSD"}))
	if err != nil {
		return ports.SourceMatchRow{}, fmt.Errorf("parse draw odds: %w", err)
	}

	awayWinOdds, err := parseOptionalFloat(firstAvailableValue(record, indexByHeader, []string{"B365A", "AvgA", "PSA"}))
	if err != nil {
		return ports.SourceMatchRow{}, fmt.Errorf("parse away odds: %w", err)
	}

	return ports.SourceMatchRow{
		CompetitionCode: competitionCode,
		CompetitionName: competitionName,
		Country:         country,
		SeasonLabel:     normalizeSeasonLabel(matchDate),
		MatchDate:       matchDate,
		HomeTeamName:    homeTeamName,
		AwayTeamName:    awayTeamName,
		HomeGoals:       homeGoals,
		AwayGoals:       awayGoals,
		HomeWinOdds:     homeWinOdds,
		DrawOdds:        drawOdds,
		AwayWinOdds:     awayWinOdds,
	}, nil
}

func fieldValue(record []string, indexByHeader map[string]int, header string) string {
	index, ok := indexByHeader[header]
	if !ok || index >= len(record) {
		return ""
	}

	return strings.TrimSpace(record[index])
}

func firstAvailableValue(record []string, indexByHeader map[string]int, headers []string) string {
	for _, header := range headers {
		value := fieldValue(record, indexByHeader, header)
		if value != "" {
			return value
		}
	}

	return ""
}

func parseMatchDate(raw string) (time.Time, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return time.Time{}, fmt.Errorf("missing match date")
	}

	layouts := []string{
		"02/01/2006",
		"02/01/06",
		"2006-01-02",
		"2006-01-02 15:04:05",
	}

	for _, layout := range layouts {
		parsed, err := time.Parse(layout, trimmed)
		if err == nil {
			return parsed, nil
		}
	}

	return time.Time{}, fmt.Errorf("invalid match date %q", raw)
}

func parseOptionalInt(raw string) (*int, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return nil, nil
	}

	parsed, err := strconv.Atoi(trimmed)
	if err != nil {
		return nil, fmt.Errorf("invalid integer %q", raw)
	}

	return &parsed, nil
}

func parseOptionalFloat(raw string) (*float64, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return nil, nil
	}

	parsed, err := strconv.ParseFloat(trimmed, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid decimal %q", raw)
	}

	return &parsed, nil
}

func normalizeCompetitionCode(raw string) string {
	return strings.ToUpper(strings.TrimSpace(raw))
}

func normalizeCompetition(code string) (string, string) {
	switch code {
	case "E0":
		return "Premier League", "England"
	default:
		return code, ""
	}
}

func normalizeSeasonLabel(matchDate time.Time) string {
	startYear := matchDate.Year()
	if matchDate.Month() < time.July {
		startYear--
	}

	return fmt.Sprintf("%d-%d", startYear, startYear+1)
}

func normalizeTeamName(raw string) string {
	parts := strings.Fields(raw)
	return strings.Join(parts, " ")
}
