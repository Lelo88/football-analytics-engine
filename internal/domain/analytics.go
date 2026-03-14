package domain

type MatchVenue string

const (
	MatchVenueAll  MatchVenue = "all"
	MatchVenueHome MatchVenue = "home"
	MatchVenueAway MatchVenue = "away"
)

type TeamForm struct {
	MatchesPlayed int
	Wins          int
	Draws         int
	Losses        int
	Points        int
	GoalsFor      int
	GoalsAgainst  int
}

type GoalsSummary struct {
	MatchesPlayed int
	GoalsFor      int
	GoalsAgainst  int
}

type OverUnderSummary struct {
	Threshold       float64
	MatchesPlayed   int
	OverCount       int
	UnderEqualCount int
}

type SeasonTeamSummary struct {
	SeasonLabel   string
	MatchesPlayed int
	Points        int
	GoalsFor      int
	GoalsAgainst  int
}
