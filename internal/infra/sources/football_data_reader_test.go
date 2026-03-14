package sources

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestFootballDataReaderReadMatches(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		_, _ = writer.Write([]byte("Div,Date,HomeTeam,AwayTeam,FTHG,FTAG,B365H,B365D,B365A\nE0,18/08/2024, Arsenal , Chelsea,2,1,1.80,3.50,4.20\n"))
	}))
	defer server.Close()

	reader := NewFootballDataReader(server.Client())
	rows, err := reader.ReadMatches(context.Background(), server.URL)
	if err != nil {
		t.Fatalf("ReadMatches returned error: %v", err)
	}

	if len(rows) != 1 {
		t.Fatalf("expected 1 row, got %d", len(rows))
	}

	row := rows[0]
	if row.CompetitionCode != "E0" {
		t.Fatalf("expected competition code E0, got %q", row.CompetitionCode)
	}
	if row.CompetitionName != "Premier League" {
		t.Fatalf("expected competition name Premier League, got %q", row.CompetitionName)
	}
	if row.Country != "England" {
		t.Fatalf("expected country England, got %q", row.Country)
	}
	if row.SeasonLabel != "2024-2025" {
		t.Fatalf("expected season label 2024-2025, got %q", row.SeasonLabel)
	}
	if row.HomeTeamName != "Arsenal" {
		t.Fatalf("expected normalized home team Arsenal, got %q", row.HomeTeamName)
	}
	if row.AwayTeamName != "Chelsea" {
		t.Fatalf("expected normalized away team Chelsea, got %q", row.AwayTeamName)
	}
	if row.HomeGoals == nil || *row.HomeGoals != 2 {
		t.Fatalf("expected home goals 2, got %+v", row.HomeGoals)
	}
	if row.AwayGoals == nil || *row.AwayGoals != 1 {
		t.Fatalf("expected away goals 1, got %+v", row.AwayGoals)
	}
	if row.HomeWinOdds == nil || *row.HomeWinOdds != 1.80 {
		t.Fatalf("expected home odds 1.80, got %+v", row.HomeWinOdds)
	}
	if !row.MatchDate.Equal(time.Date(2024, time.August, 18, 0, 0, 0, 0, time.UTC)) {
		t.Fatalf("unexpected match date: %v", row.MatchDate)
	}
}

func TestFootballDataReaderReadMatchesMalformedRow(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		_, _ = writer.Write([]byte("Div,Date,HomeTeam,AwayTeam\nE0,not-a-date,Arsenal,Chelsea\n"))
	}))
	defer server.Close()

	reader := NewFootballDataReader(server.Client())
	_, err := reader.ReadMatches(context.Background(), server.URL)
	if err == nil {
		t.Fatal("expected parsing error, got nil")
	}

	if !strings.Contains(err.Error(), "invalid match date") {
		t.Fatalf("expected invalid match date error, got %v", err)
	}
}

func TestFootballDataReaderReadMatchesUnexpectedStatus(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusBadGateway)
	}))
	defer server.Close()

	reader := NewFootballDataReader(server.Client())
	_, err := reader.ReadMatches(context.Background(), server.URL)
	if err == nil {
		t.Fatal("expected status error, got nil")
	}

	if !strings.Contains(err.Error(), "unexpected source status") {
		t.Fatalf("expected unexpected status error, got %v", err)
	}
}

func TestFootballDataReaderReadMatchesMissingHeader(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		_, _ = writer.Write([]byte("Div,HomeTeam,AwayTeam\nE0,Arsenal,Chelsea\n"))
	}))
	defer server.Close()

	reader := NewFootballDataReader(server.Client())
	_, err := reader.ReadMatches(context.Background(), server.URL)
	if err == nil {
		t.Fatal("expected missing header error, got nil")
	}

	if !strings.Contains(err.Error(), "missing required header") {
		t.Fatalf("expected missing header error, got %v", err)
	}
}

func TestFootballDataReaderReadMatchesEmptyCSV(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		_, _ = writer.Write([]byte(""))
	}))
	defer server.Close()

	reader := NewFootballDataReader(server.Client())
	rows, err := reader.ReadMatches(context.Background(), server.URL)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if len(rows) != 0 {
		t.Fatalf("expected 0 rows for empty CSV, got %d", len(rows))
	}
}

func TestFootballDataReaderReadMatchesFallbackOddsColumns(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		_, _ = writer.Write([]byte("Div,Date,HomeTeam,AwayTeam,FTHG,FTAG,AvgH,AvgD,AvgA\nE0,18/08/2024,Arsenal,Chelsea,2,1,1.91,3.40,4.00\n"))
	}))
	defer server.Close()

	reader := NewFootballDataReader(server.Client())
	rows, err := reader.ReadMatches(context.Background(), server.URL)
	if err != nil {
		t.Fatalf("ReadMatches returned error: %v", err)
	}
	if len(rows) != 1 {
		t.Fatalf("expected 1 row, got %d", len(rows))
	}

	row := rows[0]
	if row.HomeWinOdds == nil || *row.HomeWinOdds != 1.91 {
		t.Fatalf("expected home odds from AvgH 1.91, got %+v", row.HomeWinOdds)
	}
	if row.DrawOdds == nil || *row.DrawOdds != 3.40 {
		t.Fatalf("expected draw odds from AvgD 3.40, got %+v", row.DrawOdds)
	}
	if row.AwayWinOdds == nil || *row.AwayWinOdds != 4.00 {
		t.Fatalf("expected away odds from AvgA 4.00, got %+v", row.AwayWinOdds)
	}
}

func TestNormalizeSeasonLabel(t *testing.T) {
	t.Parallel()

	date := time.Date(2025, time.May, 12, 0, 0, 0, 0, time.UTC)
	if label := normalizeSeasonLabel(date); label != "2024-2025" {
		t.Fatalf("expected 2024-2025, got %q", label)
	}
}
