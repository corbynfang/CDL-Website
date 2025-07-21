-- CDL Stats Database Schema
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Seasons table
CREATE TABLE IF NOT EXISTS seasons (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    game_title VARCHAR(100) NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE,
    is_active BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Teams table
CREATE TABLE IF NOT EXISTS teams (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    abbreviation VARCHAR(10) NOT NULL,
    city VARCHAR(100),
    logo_url TEXT,
    primary_color VARCHAR(7),
    secondary_color VARCHAR(7),
    founded_date DATE,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Players table
CREATE TABLE IF NOT EXISTS players (
    id SERIAL PRIMARY KEY,
    gamertag VARCHAR(100) NOT NULL,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    country VARCHAR(3),
    birthdate DATE,
    role VARCHAR(50),
    is_active BOOLEAN DEFAULT true,
    liquipedia_url TEXT,
    twitter_handle VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Team rosters
CREATE TABLE IF NOT EXISTS team_rosters (
    id SERIAL PRIMARY KEY,
    team_id INTEGER REFERENCES teams(id),
    player_id INTEGER REFERENCES players(id),
    season_id INTEGER REFERENCES seasons(id),
    role VARCHAR(50),
    start_date DATE NOT NULL,
    end_date DATE,
    is_starter BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Tournaments table
CREATE TABLE IF NOT EXISTS tournaments (
    id SERIAL PRIMARY KEY,
    season_id INTEGER REFERENCES seasons(id),
    name VARCHAR(200) NOT NULL,
    tournament_type VARCHAR(50),
    start_date DATE NOT NULL,
    end_date DATE,
    prize_pool DECIMAL(12,2),
    location VARCHAR(100),
    tournament_format VARCHAR(50),
    liquipedia_url TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Matches table
CREATE TABLE IF NOT EXISTS matches (
    id SERIAL PRIMARY KEY,
    tournament_id INTEGER REFERENCES tournaments(id),
    team1_id INTEGER REFERENCES teams(id),
    team2_id INTEGER REFERENCES teams(id),
    match_date TIMESTAMP,
    match_type VARCHAR(50),
    format VARCHAR(20),
    team1_score INTEGER DEFAULT 0,
    team2_score INTEGER DEFAULT 0,
    winner_id INTEGER REFERENCES teams(id),
    duration_minutes INTEGER,
    vod_url TEXT,
    liquipedia_url TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Cumulative player stats per tournament
CREATE TABLE IF NOT EXISTS player_tournament_stats (
    id SERIAL PRIMARY KEY,
    player_id INTEGER REFERENCES players(id),
    team_id INTEGER REFERENCES teams(id),
    tournament_id INTEGER REFERENCES tournaments(id),
    total_kills INTEGER,
    total_deaths INTEGER,
    total_assists INTEGER,
    total_damage INTEGER,
    kd_ratio FLOAT,
    kda_ratio FLOAT
);

-- Player Match Statistics table
CREATE TABLE IF NOT EXISTS player_match_stats (
    id SERIAL PRIMARY KEY,
    match_id INTEGER REFERENCES matches(id),
    player_id INTEGER REFERENCES players(id),
    team_id INTEGER REFERENCES teams(id),
    maps_played INTEGER DEFAULT 0,
    total_kills INTEGER DEFAULT 0,
    total_deaths INTEGER DEFAULT 0,
    total_assists INTEGER DEFAULT 0,
    total_damage INTEGER DEFAULT 0,
    kd_ratio DECIMAL(4,2) DEFAULT 0,
    kda_ratio DECIMAL(4,2) DEFAULT 0,
    adr DECIMAL(6,2) DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Team Tournament Statistics table
CREATE TABLE IF NOT EXISTS team_tournament_stats (
    id SERIAL PRIMARY KEY,
    tournament_id INTEGER REFERENCES tournaments(id),
    team_id INTEGER REFERENCES teams(id),
    placement INTEGER,
    matches_played INTEGER DEFAULT 0,
    matches_won INTEGER DEFAULT 0,
    matches_lost INTEGER DEFAULT 0,
    maps_played INTEGER DEFAULT 0,
    maps_won INTEGER DEFAULT 0,
    maps_lost INTEGER DEFAULT 0,
    prize_money DECIMAL(10,2) DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Coaches table
CREATE TABLE IF NOT EXISTS coaches (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    team_id INTEGER REFERENCES teams(id),
    season_id INTEGER REFERENCES seasons(id)
);

-- Insert 2025 season (needed for foreign keys)
INSERT INTO seasons (id, name, game_title, start_date, end_date, is_active) VALUES
(6, 'CDL 2025', 'Black Ops 6', '2025-01-30', NULL, true);
