CREATE TABLE competitions (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    country VARCHAR(100),
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE seasons (
    id BIGSERIAL PRIMARY KEY,
    competition_id BIGINT NOT NULL REFERENCES competitions(id),
    label VARCHAR(20) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT seasons_competition_id_label_key UNIQUE (competition_id, label)
);

CREATE TABLE teams (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(150) NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE matches (
    id BIGSERIAL PRIMARY KEY,
    competition_id BIGINT NOT NULL REFERENCES competitions(id),
    season_id BIGINT NOT NULL REFERENCES seasons(id),
    match_date TIMESTAMP NOT NULL,
    home_team_id BIGINT NOT NULL REFERENCES teams(id),
    away_team_id BIGINT NOT NULL REFERENCES teams(id),
    home_goals INTEGER,
    away_goals INTEGER,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT matches_home_team_away_team_check CHECK (home_team_id <> away_team_id),
    CONSTRAINT matches_logical_identity_key UNIQUE (
        competition_id,
        season_id,
        match_date,
        home_team_id,
        away_team_id
    )
);

CREATE TABLE match_odds (
    id BIGSERIAL PRIMARY KEY,
    match_id BIGINT NOT NULL UNIQUE REFERENCES matches(id) ON DELETE CASCADE,
    home_win NUMERIC(10,4),
    draw NUMERIC(10,4),
    away_win NUMERIC(10,4)
);

CREATE TABLE ingestion_runs (
    id BIGSERIAL PRIMARY KEY,
    source VARCHAR(100) NOT NULL,
    started_at TIMESTAMP NOT NULL DEFAULT NOW(),
    finished_at TIMESTAMP,
    status VARCHAR(50) NOT NULL,
    rows_processed INTEGER NOT NULL DEFAULT 0,
    rows_inserted INTEGER NOT NULL DEFAULT 0,
    rows_updated INTEGER NOT NULL DEFAULT 0,
    error_message TEXT
);

CREATE INDEX idx_matches_match_date ON matches (match_date);
CREATE INDEX idx_matches_season_id ON matches (season_id);
CREATE INDEX idx_matches_home_team_id ON matches (home_team_id);
CREATE INDEX idx_matches_away_team_id ON matches (away_team_id);
CREATE INDEX idx_matches_competition_id ON matches (competition_id);
CREATE INDEX idx_seasons_competition_id ON seasons (competition_id);