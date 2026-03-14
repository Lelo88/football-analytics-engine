package usecase

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"football-analytics/internal/domain"
)

func TestIngestionServiceIngestSuccess(t *testing.T) {
	t.Parallel()

	match := domain.IngestionMatch{
		CompetitionName: "Premier League",
		Country:         "England",
		SeasonLabel:     "2024-2025",
		MatchDate:       time.Date(2024, time.August, 18, 0, 0, 0, 0, time.UTC),
		HomeTeamName:    "Arsenal",
		AwayTeamName:    "Chelsea",
	}

	reader := &fakeReader{matches: []domain.IngestionMatch{match}}
	repository := &fakeRepository{}
	service := NewIngestionService(reader, repository)

	err := service.Ingest(context.Background())
	if err != nil {
		t.Fatalf("Ingest returned error: %v", err)
	}

	if reader.fetchCalls != 1 {
		t.Fatalf("expected reader fetchCalls=1, got %d", reader.fetchCalls)
	}
	if repository.saveCalls != 1 {
		t.Fatalf("expected repository saveCalls=1, got %d", repository.saveCalls)
	}
	if len(repository.savedMatches) != 1 {
		t.Fatalf("expected one saved match, got %d", len(repository.savedMatches))
	}
	saved := repository.savedMatches[0]
	if saved.CompetitionName != match.CompetitionName || saved.HomeTeamName != match.HomeTeamName || saved.AwayTeamName != match.AwayTeamName {
		t.Fatalf("unexpected saved match: %+v", saved)
	}
}

func TestIngestionServiceIngestReaderError(t *testing.T) {
	t.Parallel()

	reader := &fakeReader{err: fmt.Errorf("source unavailable")}
	repository := &fakeRepository{}
	service := NewIngestionService(reader, repository)

	err := service.Ingest(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "fetch matches") {
		t.Fatalf("expected fetch error context, got %v", err)
	}
	if repository.saveCalls != 0 {
		t.Fatalf("expected repository not called, got saveCalls=%d", repository.saveCalls)
	}
}

func TestIngestionServiceIngestInvalidData(t *testing.T) {
	t.Parallel()

	reader := &fakeReader{matches: []domain.IngestionMatch{{
		CompetitionName: "",
		SeasonLabel:     "2024-2025",
		MatchDate:       time.Date(2024, time.August, 18, 0, 0, 0, 0, time.UTC),
		HomeTeamName:    "Arsenal",
		AwayTeamName:    "Chelsea",
	}}}
	repository := &fakeRepository{}
	service := NewIngestionService(reader, repository)

	err := service.Ingest(context.Background())
	if err == nil {
		t.Fatal("expected validation error, got nil")
	}
	if !strings.Contains(err.Error(), "validate match") {
		t.Fatalf("expected validation context, got %v", err)
	}
	if repository.saveCalls != 0 {
		t.Fatalf("expected repository not called, got saveCalls=%d", repository.saveCalls)
	}
}

func TestIngestionServiceIngestRejectsSameTeamIgnoringCase(t *testing.T) {
	t.Parallel()

	reader := &fakeReader{matches: []domain.IngestionMatch{{
		CompetitionName: "Premier League",
		SeasonLabel:     "2024-2025",
		MatchDate:       time.Date(2024, time.August, 18, 0, 0, 0, 0, time.UTC),
		HomeTeamName:    "Arsenal",
		AwayTeamName:    "ARSENAL",
	}}}
	repository := &fakeRepository{}
	service := NewIngestionService(reader, repository)

	err := service.Ingest(context.Background())
	if err == nil {
		t.Fatal("expected validation error, got nil")
	}
	if !strings.Contains(err.Error(), "home and away teams must be different") {
		t.Fatalf("expected same-team validation error, got %v", err)
	}
	if repository.saveCalls != 0 {
		t.Fatalf("expected repository not called, got saveCalls=%d", repository.saveCalls)
	}
}

func TestIngestionServiceIngestRepositoryErrorStopsExecution(t *testing.T) {
	t.Parallel()

	matches := []domain.IngestionMatch{
		{
			CompetitionName: "Premier League",
			Country:         "England",
			SeasonLabel:     "2024-2025",
			MatchDate:       time.Date(2024, time.August, 18, 0, 0, 0, 0, time.UTC),
			HomeTeamName:    "Arsenal",
			AwayTeamName:    "Chelsea",
		},
		{
			CompetitionName: "Premier League",
			Country:         "England",
			SeasonLabel:     "2024-2025",
			MatchDate:       time.Date(2024, time.August, 19, 0, 0, 0, 0, time.UTC),
			HomeTeamName:    "Liverpool",
			AwayTeamName:    "Tottenham",
		},
	}

	reader := &fakeReader{matches: matches}
	repository := &fakeRepository{errAtCall: 1, err: fmt.Errorf("db failure")}
	service := NewIngestionService(reader, repository)

	err := service.Ingest(context.Background())
	if err == nil {
		t.Fatal("expected repository error, got nil")
	}
	if !strings.Contains(err.Error(), "save match 0") {
		t.Fatalf("expected save error with index, got %v", err)
	}
	if repository.saveCalls != 1 {
		t.Fatalf("expected execution to stop on first save error, got saveCalls=%d", repository.saveCalls)
	}
}

func TestIngestionServiceIngestEmptyList(t *testing.T) {
	t.Parallel()

	reader := &fakeReader{matches: []domain.IngestionMatch{}}
	repository := &fakeRepository{}
	service := NewIngestionService(reader, repository)

	err := service.Ingest(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if reader.fetchCalls != 1 {
		t.Fatalf("expected reader fetchCalls=1, got %d", reader.fetchCalls)
	}
	if repository.saveCalls != 0 {
		t.Fatalf("expected no repository calls, got %d", repository.saveCalls)
	}
}

type fakeReader struct {
	matches    []domain.IngestionMatch
	err        error
	fetchCalls int
}

func (reader *fakeReader) Fetch(ctx context.Context) ([]domain.IngestionMatch, error) {
	reader.fetchCalls++
	if reader.err != nil {
		return nil, reader.err
	}
	return reader.matches, nil
}

type fakeRepository struct {
	savedMatches []domain.IngestionMatch
	saveCalls    int
	errAtCall    int
	err          error
}

func (repository *fakeRepository) Save(ctx context.Context, match domain.IngestionMatch) error {
	repository.saveCalls++
	if repository.err != nil && repository.saveCalls == repository.errAtCall {
		return repository.err
	}
	repository.savedMatches = append(repository.savedMatches, match)
	return nil
}
