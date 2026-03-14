package usecase

import (
	"context"
	"fmt"

	"football-analytics/internal/domain"
	"football-analytics/internal/ports"
)

type TeamAnalyticsService struct {
	repository ports.TeamAnalyticsReadRepository
}

func NewTeamAnalyticsService(repository ports.TeamAnalyticsReadRepository) *TeamAnalyticsService {
	return &TeamAnalyticsService{repository: repository}
}

func (service *TeamAnalyticsService) Teams(ctx context.Context) ([]domain.Team, error) {
	result, err := service.repository.ListTeams(ctx)
	if err != nil {
		return nil, fmt.Errorf("list teams: %w", err)
	}

	return result, nil
}

func (service *TeamAnalyticsService) TeamExists(ctx context.Context, teamID int64) (bool, error) {
	if teamID <= 0 {
		return false, fmt.Errorf("team id must be greater than zero")
	}

	exists, err := service.repository.TeamExists(ctx, teamID)
	if err != nil {
		return false, fmt.Errorf("check team existence: %w", err)
	}

	return exists, nil
}

func (service *TeamAnalyticsService) TeamForm(ctx context.Context, filter ports.TeamQueryFilter) (domain.TeamForm, error) {
	if err := validateTeamQueryFilter(filter); err != nil {
		return domain.TeamForm{}, err
	}

	result, err := service.repository.GetTeamForm(ctx, filter)
	if err != nil {
		return domain.TeamForm{}, fmt.Errorf("get team form: %w", err)
	}

	return result, nil
}

func (service *TeamAnalyticsService) GoalsSummary(ctx context.Context, filter ports.TeamQueryFilter) (domain.GoalsSummary, error) {
	if err := validateTeamQueryFilter(filter); err != nil {
		return domain.GoalsSummary{}, err
	}

	result, err := service.repository.GetGoalsSummary(ctx, filter)
	if err != nil {
		return domain.GoalsSummary{}, fmt.Errorf("get goals summary: %w", err)
	}

	return result, nil
}

func (service *TeamAnalyticsService) OverUnderSummary(ctx context.Context, filter ports.OverUnderQueryFilter) (domain.OverUnderSummary, error) {
	if err := validateOverUnderFilter(filter); err != nil {
		return domain.OverUnderSummary{}, err
	}

	result, err := service.repository.GetOverUnderSummary(ctx, filter)
	if err != nil {
		return domain.OverUnderSummary{}, fmt.Errorf("get over/under summary: %w", err)
	}

	return result, nil
}

func (service *TeamAnalyticsService) SeasonSummaries(ctx context.Context, filter ports.SeasonSummaryFilter) ([]domain.SeasonTeamSummary, error) {
	if err := validateSeasonSummaryFilter(filter); err != nil {
		return nil, err
	}

	result, err := service.repository.GetSeasonSummaries(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("get season summaries: %w", err)
	}

	return result, nil
}

func validateTeamQueryFilter(filter ports.TeamQueryFilter) error {
	if filter.TeamID <= 0 {
		return fmt.Errorf("team id must be greater than zero")
	}
	if filter.LastN < 0 {
		return fmt.Errorf("last_n must be zero or greater")
	}
	if err := validateVenue(filter.Venue); err != nil {
		return err
	}

	return nil
}

func validateOverUnderFilter(filter ports.OverUnderQueryFilter) error {
	if filter.TeamID <= 0 {
		return fmt.Errorf("team id must be greater than zero")
	}
	if filter.LastN < 0 {
		return fmt.Errorf("last_n must be zero or greater")
	}
	if filter.Threshold < 0 {
		return fmt.Errorf("threshold must be zero or greater")
	}
	if err := validateVenue(filter.Venue); err != nil {
		return err
	}

	return nil
}

func validateSeasonSummaryFilter(filter ports.SeasonSummaryFilter) error {
	if filter.TeamID <= 0 {
		return fmt.Errorf("team id must be greater than zero")
	}
	if err := validateVenue(filter.Venue); err != nil {
		return err
	}

	return nil
}

func validateVenue(venue domain.MatchVenue) error {
	resolvedVenue := venue
	if resolvedVenue == "" {
		resolvedVenue = domain.MatchVenueAll
	}

	switch resolvedVenue {
	case domain.MatchVenueAll, domain.MatchVenueHome, domain.MatchVenueAway:
		return nil
	default:
		return fmt.Errorf("invalid venue filter %q", venue)
	}
}
