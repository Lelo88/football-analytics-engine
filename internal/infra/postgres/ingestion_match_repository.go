package postgres

import (
	"context"
	"database/sql"

	"football-analytics/internal/domain"
	"football-analytics/internal/ports"
)

type IngestionMatchRepository struct {
	competitionRepository *CompetitionRepository
	seasonRepository      *SeasonRepository
	teamRepository        *TeamRepository
	matchRepository       *MatchRepository
}

func NewIngestionMatchRepository(db *sql.DB) *IngestionMatchRepository {
	return &IngestionMatchRepository{
		competitionRepository: NewCompetitionRepository(db),
		seasonRepository:      NewSeasonRepository(db),
		teamRepository:        NewTeamRepository(db),
		matchRepository:       NewMatchRepository(db),
	}
}

func (repository *IngestionMatchRepository) Save(ctx context.Context, match domain.IngestionMatch) (ports.IngestionSaveResult, error) {
	competition, err := repository.competitionRepository.CreateOrGet(ctx, match.CompetitionName, match.Country)
	if err != nil {
		return ports.IngestionSaveResult{}, err
	}

	season, err := repository.seasonRepository.CreateOrGet(ctx, competition.ID, match.SeasonLabel)
	if err != nil {
		return ports.IngestionSaveResult{}, err
	}

	homeTeam, err := repository.teamRepository.CreateOrGet(ctx, match.HomeTeamName)
	if err != nil {
		return ports.IngestionSaveResult{}, err
	}

	awayTeam, err := repository.teamRepository.CreateOrGet(ctx, match.AwayTeamName)
	if err != nil {
		return ports.IngestionSaveResult{}, err
	}

	upsertResult, err := repository.matchRepository.UpsertMatch(ctx, ports.MatchUpsertParams{
		CompetitionID: competition.ID,
		SeasonID:      season.ID,
		MatchDate:     match.MatchDate,
		HomeTeamID:    homeTeam.ID,
		AwayTeamID:    awayTeam.ID,
		HomeGoals:     match.HomeGoals,
		AwayGoals:     match.AwayGoals,
	})
	if err != nil {
		return ports.IngestionSaveResult{}, err
	}

	if match.HomeWinOdds != nil || match.DrawOdds != nil || match.AwayWinOdds != nil {
		err = repository.matchRepository.UpsertMatchOdds(ctx, ports.MatchOddsUpsertParams{
			MatchID: upsertResult.Match.ID,
			HomeWin: match.HomeWinOdds,
			Draw:    match.DrawOdds,
			AwayWin: match.AwayWinOdds,
		})
		if err != nil {
			return ports.IngestionSaveResult{}, err
		}
	}

	return ports.IngestionSaveResult{Inserted: upsertResult.Inserted, Updated: upsertResult.Updated}, nil
}
