package domain

import "time"

const (
	IngestionRunStatusStarted = "started"
	IngestionRunStatusSuccess = "success"
	IngestionRunStatusFailed  = "failed"
)

type Competition struct {
	ID        int64
	Name      string
	Country   string
	CreatedAt time.Time
}

type Season struct {
	ID            int64
	CompetitionID int64
	Label         string
	CreatedAt     time.Time
}

type Team struct {
	ID        int64
	Name      string
	CreatedAt time.Time
}

type Match struct {
	ID            int64
	CompetitionID int64
	SeasonID      int64
	MatchDate     time.Time
	HomeTeamID    int64
	AwayTeamID    int64
	HomeGoals     *int
	AwayGoals     *int
	CreatedAt     time.Time
}

type MatchOdds struct {
	ID      int64
	MatchID int64
	HomeWin *float64
	Draw    *float64
	AwayWin *float64
}

type IngestionMatch struct {
	CompetitionName string
	Country         string
	SeasonLabel     string
	MatchDate       time.Time
	HomeTeamName    string
	AwayTeamName    string
	HomeGoals       *int
	AwayGoals       *int
	HomeWinOdds     *float64
	DrawOdds        *float64
	AwayWinOdds     *float64
}

type IngestionRun struct {
	ID            int64
	Source        string
	StartedAt     time.Time
	FinishedAt    *time.Time
	Status        string
	RowsProcessed int
	RowsInserted  int
	RowsUpdated   int
	ErrorMessage  *string
}
