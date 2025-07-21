-- CDL Champs 2025 Data
-- Assumes schema.sql has already been run

-- Insert Champs Teams
INSERT INTO teams (name, abbreviation, city, is_active) VALUES
('LA Thieves', 'LAT', 'Los Angeles', true),
('ATL FaZe', 'ATL', 'Atlanta', true),
('TOR Ultra', 'TOR', 'Toronto', true),
('VAN Surge', 'VAN', 'Vancouver', true),
('MIA Heretics', 'MIA', 'Miami', true),
('CAR Royal Ravens', 'CAR', 'Charlotte', true),
('OpTic TEX', 'TEX', 'Dallas', true),
('BOS Breach', 'BOS', 'Boston', true)
ON CONFLICT DO NOTHING;

-- Insert Champs Coaches
CREATE TABLE IF NOT EXISTS coaches (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    team_id INTEGER REFERENCES teams(id),
    season_id INTEGER REFERENCES seasons(id)
);

INSERT INTO coaches (name, team_id, season_id) VALUES
('ShAnE', (SELECT id FROM teams WHERE name = 'LA Thieves'), 6),
('Sender', (SELECT id FROM teams WHERE name = 'LA Thieves'), 6),
('Crowder', (SELECT id FROM teams WHERE name = 'ATL FaZe'), 6),
('RJ', (SELECT id FROM teams WHERE name = 'ATL FaZe'), 6),
('Flux', (SELECT id FROM teams WHERE name = 'TOR Ultra'), 6),
('Joee', (SELECT id FROM teams WHERE name = 'TOR Ultra'), 6),
('Rambo', (SELECT id FROM teams WHERE name = 'VAN Surge'), 6),
('MethodZ', (SELECT id FROM teams WHERE name = 'MIA Heretics'), 6),
('Sikotik', (SELECT id FROM teams WHERE name = 'MIA Heretics'), 6),
('Saintt', (SELECT id FROM teams WHERE name = 'CAR Royal Ravens'), 6),
('Karma', (SELECT id FROM teams WHERE name = 'OpTic TEX'), 6),
('JP Krez', (SELECT id FROM teams WHERE name = 'OpTic TEX'), 6),
('Seany', (SELECT id FROM teams WHERE name = 'BOS Breach'), 6),
('Magxck', (SELECT id FROM teams WHERE name = 'BOS Breach'), 6);

-- Insert Champs Players (excluding coaches)
INSERT INTO players (gamertag, is_active) VALUES
-- LA Thieves
('Ghosty', true),
('Scrap', true),
('HyDra', true),
('Envoy', true),
-- ATL FaZe
('aBeZy', true),
('Cellium', true),
('Simp', true),
('Drazah', true),
-- TOR Ultra
('CleanX', true),
('Insight', true),
('Beans', true),
('JoeDeceives', true),
-- VAN Surge
('Abuzah', true),
('04', true),
('Nastie', true),
('Neptune', true),
-- MIA Heretics
('MettalZ', true),
('ReeaL', true),
('RenKoR', true),
('SupeR', true),
-- CAR Royal Ravens
('Gwinn', true),
('TJHaLy', true),
('SlasheR', true),
('Wrecks', true),
-- OpTic TEX
('Shotzzy', true),
('Dashy', true),
('Huke', true),
('Mercules', true),
-- BOS Breach
('Snoopy', true),
('Cammy', true),
('Owakening', true),
('Purj', true)
ON CONFLICT DO NOTHING;

-- Assign players to teams for 2025 season (season_id = 6)
INSERT INTO team_rosters (team_id, player_id, season_id, is_starter, start_date) VALUES
-- LA Thieves
((SELECT id FROM teams WHERE name = 'LA Thieves'), (SELECT id FROM players WHERE gamertag = 'Ghosty'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'LA Thieves'), (SELECT id FROM players WHERE gamertag = 'Scrap'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'LA Thieves'), (SELECT id FROM players WHERE gamertag = 'HyDra'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'LA Thieves'), (SELECT id FROM players WHERE gamertag = 'Envoy'), 6, true, '2024-11-01'),
-- ATL FaZe
((SELECT id FROM teams WHERE name = 'ATL FaZe'), (SELECT id FROM players WHERE gamertag = 'aBeZy'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'ATL FaZe'), (SELECT id FROM players WHERE gamertag = 'Cellium'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'ATL FaZe'), (SELECT id FROM players WHERE gamertag = 'Simp'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'ATL FaZe'), (SELECT id FROM players WHERE gamertag = 'Drazah'), 6, true, '2024-11-01'),
-- TOR Ultra
((SELECT id FROM teams WHERE name = 'TOR Ultra'), (SELECT id FROM players WHERE gamertag = 'CleanX'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'TOR Ultra'), (SELECT id FROM players WHERE gamertag = 'Insight'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'TOR Ultra'), (SELECT id FROM players WHERE gamertag = 'Beans'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'TOR Ultra'), (SELECT id FROM players WHERE gamertag = 'JoeDeceives'), 6, true, '2024-11-01'),
-- VAN Surge
((SELECT id FROM teams WHERE name = 'VAN Surge'), (SELECT id FROM players WHERE gamertag = 'Abuzah'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'VAN Surge'), (SELECT id FROM players WHERE gamertag = '04'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'VAN Surge'), (SELECT id FROM players WHERE gamertag = 'Nastie'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'VAN Surge'), (SELECT id FROM players WHERE gamertag = 'Neptune'), 6, true, '2024-11-01'),
-- MIA Heretics
((SELECT id FROM teams WHERE name = 'MIA Heretics'), (SELECT id FROM players WHERE gamertag = 'MettalZ'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'MIA Heretics'), (SELECT id FROM players WHERE gamertag = 'ReeaL'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'MIA Heretics'), (SELECT id FROM players WHERE gamertag = 'RenKoR'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'MIA Heretics'), (SELECT id FROM players WHERE gamertag = 'SupeR'), 6, true, '2024-11-01'),
-- CAR Royal Ravens
((SELECT id FROM teams WHERE name = 'CAR Royal Ravens'), (SELECT id FROM players WHERE gamertag = 'Gwinn'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'CAR Royal Ravens'), (SELECT id FROM players WHERE gamertag = 'TJHaLy'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'CAR Royal Ravens'), (SELECT id FROM players WHERE gamertag = 'SlasheR'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'CAR Royal Ravens'), (SELECT id FROM players WHERE gamertag = 'Wrecks'), 6, true, '2024-11-01'),
-- OpTic TEX
((SELECT id FROM teams WHERE name = 'OpTic TEX'), (SELECT id FROM players WHERE gamertag = 'Shotzzy'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'OpTic TEX'), (SELECT id FROM players WHERE gamertag = 'Dashy'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'OpTic TEX'), (SELECT id FROM players WHERE gamertag = 'Huke'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'OpTic TEX'), (SELECT id FROM players WHERE gamertag = 'Mercules'), 6, true, '2024-11-01'),
-- BOS Breach
((SELECT id FROM teams WHERE name = 'BOS Breach'), (SELECT id FROM players WHERE gamertag = 'Snoopy'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'BOS Breach'), (SELECT id FROM players WHERE gamertag = 'Cammy'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'BOS Breach'), (SELECT id FROM players WHERE gamertag = 'Owakening'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'BOS Breach'), (SELECT id FROM players WHERE gamertag = 'Purj'), 6, true, '2024-11-01');

-- Insert cumulative player stats for Champs (players only, using OVERALL K, D, KD from your table)
INSERT INTO player_tournament_stats (player_id, team_id, tournament_id, total_kills, total_deaths, total_assists, total_damage, kd_ratio, kda_ratio) VALUES
-- LA Thieves
((SELECT id FROM players WHERE gamertag = 'Ghosty'), (SELECT id FROM teams WHERE name = 'LA Thieves'), 5, 143, 158, 0, 0, 0.905, 0.905),
((SELECT id FROM players WHERE gamertag = 'Scrap'), (SELECT id FROM teams WHERE name = 'LA Thieves'), 5, 158, 165, 0, 0, 0.958, 0.958),
((SELECT id FROM players WHERE gamertag = 'HyDra'), (SELECT id FROM teams WHERE name = 'LA Thieves'), 5, 183, 159, 0, 0, 1.151, 1.151),
((SELECT id FROM players WHERE gamertag = 'Envoy'), (SELECT id FROM teams WHERE name = 'LA Thieves'), 5, 169, 169, 0, 0, 1.000, 1.000),
-- ATL FaZe
((SELECT id FROM players WHERE gamertag = 'aBeZy'), (SELECT id FROM teams WHERE name = 'ATL FaZe'), 5, 162, 157, 0, 0, 1.032, 1.032),
((SELECT id FROM players WHERE gamertag = 'Cellium'), (SELECT id FROM teams WHERE name = 'ATL FaZe'), 5, 151, 132, 0, 0, 1.144, 1.144),
((SELECT id FROM players WHERE gamertag = 'Simp'), (SELECT id FROM teams WHERE name = 'ATL FaZe'), 5, 131, 161, 0, 0, 0.814, 0.814),
((SELECT id FROM players WHERE gamertag = 'Drazah'), (SELECT id FROM teams WHERE name = 'ATL FaZe'), 5, 141, 170, 0, 0, 0.829, 0.829),
-- TOR Ultra
((SELECT id FROM players WHERE gamertag = 'CleanX'), (SELECT id FROM teams WHERE name = 'TOR Ultra'), 5, 195, 209, 0, 0, 0.933, 0.933),
((SELECT id FROM players WHERE gamertag = 'Insight'), (SELECT id FROM teams WHERE name = 'TOR Ultra'), 5, 151, 190, 0, 0, 0.795, 0.795),
((SELECT id FROM players WHERE gamertag = 'Beans'), (SELECT id FROM teams WHERE name = 'TOR Ultra'), 5, 184, 190, 0, 0, 0.968, 0.968),
((SELECT id FROM players WHERE gamertag = 'JoeDeceives'), (SELECT id FROM teams WHERE name = 'TOR Ultra'), 5, 196, 190, 0, 0, 1.032, 1.032),
-- VAN Surge
((SELECT id FROM players WHERE gamertag = 'Abuzah'), (SELECT id FROM teams WHERE name = 'VAN Surge'), 5, 517, 527, 0, 0, 0.981, 0.981),
((SELECT id FROM players WHERE gamertag = '04'), (SELECT id FROM teams WHERE name = 'VAN Surge'), 5, 543, 563, 0, 0, 0.965, 0.965),
((SELECT id FROM players WHERE gamertag = 'Nastie'), (SELECT id FROM teams WHERE name = 'VAN Surge'), 5, 494, 493, 0, 0, 1.002, 1.002),
((SELECT id FROM players WHERE gamertag = 'Neptune'), (SELECT id FROM teams WHERE name = 'VAN Surge'), 5, 569, 551, 0, 0, 1.033, 1.033),
-- MIA Heretics
((SELECT id FROM players WHERE gamertag = 'MettalZ'), (SELECT id FROM teams WHERE name = 'MIA Heretics'), 5, 291, 329, 0, 0, 0.885, 0.885),
((SELECT id FROM players WHERE gamertag = 'ReeaL'), (SELECT id FROM teams WHERE name = 'MIA Heretics'), 5, 326, 295, 0, 0, 1.105, 1.105),
((SELECT id FROM players WHERE gamertag = 'RenKoR'), (SELECT id FROM teams WHERE name = 'MIA Heretics'), 5, 355, 288, 0, 0, 1.233, 1.233),
((SELECT id FROM players WHERE gamertag = 'SupeR'), (SELECT id FROM teams WHERE name = 'MIA Heretics'), 5, 315, 271, 0, 0, 1.162, 1.162),
-- CAR Royal Ravens
((SELECT id FROM players WHERE gamertag = 'Gwinn'), (SELECT id FROM teams WHERE name = 'CAR Royal Ravens'), 5, 166, 142, 0, 0, 1.169, 1.169),
((SELECT id FROM players WHERE gamertag = 'TJHaLy'), (SELECT id FROM teams WHERE name = 'CAR Royal Ravens'), 5, 127, 156, 0, 0, 0.814, 0.814),
((SELECT id FROM players WHERE gamertag = 'SlasheR'), (SELECT id FROM teams WHERE name = 'CAR Royal Ravens'), 5, 113, 135, 0, 0, 0.837, 0.837),
((SELECT id FROM players WHERE gamertag = 'Wrecks'), (SELECT id FROM teams WHERE name = 'CAR Royal Ravens'), 5, 107, 145, 0, 0, 0.738, 0.738),
-- OpTic TEX
((SELECT id FROM players WHERE gamertag = 'Shotzzy'), (SELECT id FROM teams WHERE name = 'OpTic TEX'), 5, 338, 273, 0, 0, 1.238, 1.238),
((SELECT id FROM players WHERE gamertag = 'Dashy'), (SELECT id FROM teams WHERE name = 'OpTic TEX'), 5, 306, 248, 0, 0, 1.234, 1.234),
((SELECT id FROM players WHERE gamertag = 'Huke'), (SELECT id FROM teams WHERE name = 'OpTic TEX'), 5, 313, 290, 0, 0, 1.079, 1.079),
((SELECT id FROM players WHERE gamertag = 'Mercules'), (SELECT id FROM teams WHERE name = 'OpTic TEX'), 5, 316, 270, 0, 0, 1.170, 1.170),
-- BOS Breach
((SELECT id FROM players WHERE gamertag = 'Snoopy'), (SELECT id FROM teams WHERE name = 'BOS Breach'), 5, 325, 336, 0, 0, 0.967, 0.967),
((SELECT id FROM players WHERE gamertag = 'Cammy'), (SELECT id FROM teams WHERE name = 'BOS Breach'), 5, 299, 338, 0, 0, 0.885, 0.885),
((SELECT id FROM players WHERE gamertag = 'Owakening'), (SELECT id FROM teams WHERE name = 'BOS Breach'), 5, 285, 338, 0, 0, 0.843, 0.843),
((SELECT id FROM players WHERE gamertag = 'Purj'), (SELECT id FROM teams WHERE name = 'BOS Breach'), 5, 315, 339, 0, 0, 0.929, 0.929)
ON CONFLICT DO NOTHING; 