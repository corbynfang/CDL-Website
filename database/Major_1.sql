-- Major 1 Data for CDL 2025
-- Assumes schema.sql has already been run

-- Insert Major 1 Teams
INSERT INTO teams (name, abbreviation, city, is_active) VALUES
('OpTic TEX', 'TEX', 'Dallas', true),
('TOR Ultra', 'TOR', 'Toronto', true),
('BOS Breach', 'BOS', 'Boston', true),
('CAR Royal Ravens', 'CAR', 'Charlotte', true),
('LA Thieves', 'LAT', 'Los Angeles', true),
('ATL FaZe', 'ATL', 'Atlanta', true),
('VAN Surge', 'VAN', 'Vancouver', true),
('MIA Heretics', 'MIA', 'Miami', true),
('LA Guerrillas M8', 'LAG', 'Los Angeles', true),
('MIN RØKKR', 'MIN', 'Minneapolis', true),
('Cloud9 NY', 'NYC', 'New York', true),
('LV Falcons', 'LVF', 'Las Vegas', true)
ON CONFLICT DO NOTHING;

-- Insert Major 1 Coaches
CREATE TABLE IF NOT EXISTS coaches (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    team_id INTEGER REFERENCES teams(id),
    season_id INTEGER REFERENCES seasons(id)
);

INSERT INTO coaches (name, team_id, season_id) VALUES
('JP Krez', (SELECT id FROM teams WHERE name = 'OpTic TEX'), 6),
('Karma', (SELECT id FROM teams WHERE name = 'OpTic TEX'), 6),
('Flux', (SELECT id FROM teams WHERE name = 'TOR Ultra'), 6),
('Joee', (SELECT id FROM teams WHERE name = 'TOR Ultra'), 6),
('Seany', (SELECT id FROM teams WHERE name = 'BOS Breach'), 6),
('Magxck', (SELECT id FROM teams WHERE name = 'BOS Breach'), 6),
('Saintt', (SELECT id FROM teams WHERE name = 'CAR Royal Ravens'), 6),
('ShAnE', (SELECT id FROM teams WHERE name = 'LA Thieves'), 6),
('Sender', (SELECT id FROM teams WHERE name = 'LA Thieves'), 6),
('Crowder', (SELECT id FROM teams WHERE name = 'ATL FaZe'), 6),
('RJ', (SELECT id FROM teams WHERE name = 'ATL FaZe'), 6),
('Rambo', (SELECT id FROM teams WHERE name = 'VAN Surge'), 6),
('Methodz', (SELECT id FROM teams WHERE name = 'MIA Heretics'), 6),
('Sikotik', (SELECT id FROM teams WHERE name = 'MIA Heretics'), 6),
('DREAL', (SELECT id FROM teams WHERE name = 'LA Guerrillas M8'), 6),
('MarkyB', (SELECT id FROM teams WHERE name = 'LA Guerrillas M8'), 6),
('Loony', (SELECT id FROM teams WHERE name = 'MIN RØKKR'), 6),
('Alexdotzip', (SELECT id FROM teams WHERE name = 'MIN RØKKR'), 6),
('Accuracy', (SELECT id FROM teams WHERE name = 'Cloud9 NY'), 6),
('Arian', (SELECT id FROM teams WHERE name = 'Cloud9 NY'), 6),
('LewTee', (SELECT id FROM teams WHERE name = 'LV Falcons'), 6),
('Clayster', (SELECT id FROM teams WHERE name = 'LV Falcons'), 6);

-- Insert Major 1 Players (excluding coaches)
INSERT INTO players (gamertag, is_active) VALUES
-- OpTic TEX
('Shotzzy', true),
('Dashy', true),
('Kenny', true),
('Huke', true),
-- TOR Ultra
('CleanX', true),
('Insight', true),
('Beans', true),
('JoeDeceives', true),
-- BOS Breach
('Snoopy', true),
('Cammy', true),
('Owakening', true),
('Purj', true),
-- CAR Royal Ravens
('Gwinn', true),
('TJHaLy', true),
('SlasheR', true),
('Vivid', true),
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
-- VAN Surge
('Abuzah', true),
('04', true),
('Hicksy', true),
('Nastie', true),
-- MIA Heretics
('Lucky', true),
('MettalZ', true),
('ReeaL', true),
('RenKoR', true),
-- LA Guerrillas M8
('Lynz', true),
('KiSMET', true),
('Skyz', true),
('Priestahh', true),
-- MIN RØKKR
('Nero', true),
('Gio', true),
('Estreal', true),
('PaulEhx', true),
-- Cloud9 NY
('Sib', true),
('Attach', true),
('Kremp', true),
('Mack', true),
-- LV Falcons
('Roxas', true),
('Exnid', true),
('d7oom', true),
('KiinG', true)
ON CONFLICT DO NOTHING;

-- Assign players to teams for 2025 season (season_id = 6)
INSERT INTO team_rosters (team_id, player_id, season_id, is_starter, start_date) VALUES
-- OpTic TEX
((SELECT id FROM teams WHERE name = 'OpTic TEX'), (SELECT id FROM players WHERE gamertag = 'Shotzzy'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'OpTic TEX'), (SELECT id FROM players WHERE gamertag = 'Dashy'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'OpTic TEX'), (SELECT id FROM players WHERE gamertag = 'Kenny'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'OpTic TEX'), (SELECT id FROM players WHERE gamertag = 'Huke'), 6, true, '2024-11-01'),
-- TOR Ultra
((SELECT id FROM teams WHERE name = 'TOR Ultra'), (SELECT id FROM players WHERE gamertag = 'CleanX'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'TOR Ultra'), (SELECT id FROM players WHERE gamertag = 'Insight'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'TOR Ultra'), (SELECT id FROM players WHERE gamertag = 'Beans'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'TOR Ultra'), (SELECT id FROM players WHERE gamertag = 'JoeDeceives'), 6, true, '2024-11-01'),
-- BOS Breach
((SELECT id FROM teams WHERE name = 'BOS Breach'), (SELECT id FROM players WHERE gamertag = 'Snoopy'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'BOS Breach'), (SELECT id FROM players WHERE gamertag = 'Cammy'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'BOS Breach'), (SELECT id FROM players WHERE gamertag = 'Owakening'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'BOS Breach'), (SELECT id FROM players WHERE gamertag = 'Purj'), 6, true, '2024-11-01'),
-- CAR Royal Ravens
((SELECT id FROM teams WHERE name = 'CAR Royal Ravens'), (SELECT id FROM players WHERE gamertag = 'Gwinn'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'CAR Royal Ravens'), (SELECT id FROM players WHERE gamertag = 'TJHaLy'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'CAR Royal Ravens'), (SELECT id FROM players WHERE gamertag = 'SlasheR'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'CAR Royal Ravens'), (SELECT id FROM players WHERE gamertag = 'Vivid'), 6, true, '2024-11-01'),
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
-- VAN Surge
((SELECT id FROM teams WHERE name = 'VAN Surge'), (SELECT id FROM players WHERE gamertag = 'Abuzah'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'VAN Surge'), (SELECT id FROM players WHERE gamertag = '04'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'VAN Surge'), (SELECT id FROM players WHERE gamertag = 'Hicksy'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'VAN Surge'), (SELECT id FROM players WHERE gamertag = 'Nastie'), 6, true, '2024-11-01'),
-- MIA Heretics
((SELECT id FROM teams WHERE name = 'MIA Heretics'), (SELECT id FROM players WHERE gamertag = 'Lucky'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'MIA Heretics'), (SELECT id FROM players WHERE gamertag = 'MettalZ'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'MIA Heretics'), (SELECT id FROM players WHERE gamertag = 'ReeaL'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'MIA Heretics'), (SELECT id FROM players WHERE gamertag = 'RenKoR'), 6, true, '2024-11-01'),
-- LA Guerrillas M8
((SELECT id FROM teams WHERE name = 'LA Guerrillas M8'), (SELECT id FROM players WHERE gamertag = 'Lynz'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'LA Guerrillas M8'), (SELECT id FROM players WHERE gamertag = 'KiSMET'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'LA Guerrillas M8'), (SELECT id FROM players WHERE gamertag = 'Skyz'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'LA Guerrillas M8'), (SELECT id FROM players WHERE gamertag = 'Priestahh'), 6, true, '2024-11-01'),
-- MIN RØKKR
((SELECT id FROM teams WHERE name = 'MIN RØKKR'), (SELECT id FROM players WHERE gamertag = 'Nero'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'MIN RØKKR'), (SELECT id FROM players WHERE gamertag = 'Gio'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'MIN RØKKR'), (SELECT id FROM players WHERE gamertag = 'Estreal'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'MIN RØKKR'), (SELECT id FROM players WHERE gamertag = 'PaulEhx'), 6, true, '2024-11-01'),
-- Cloud9 NY
((SELECT id FROM teams WHERE name = 'Cloud9 NY'), (SELECT id FROM players WHERE gamertag = 'Sib'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'Cloud9 NY'), (SELECT id FROM players WHERE gamertag = 'Attach'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'Cloud9 NY'), (SELECT id FROM players WHERE gamertag = 'Kremp'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'Cloud9 NY'), (SELECT id FROM players WHERE gamertag = 'Mack'), 6, true, '2024-11-01'),
-- LV Falcons
((SELECT id FROM teams WHERE name = 'LV Falcons'), (SELECT id FROM players WHERE gamertag = 'Roxas'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'LV Falcons'), (SELECT id FROM players WHERE gamertag = 'Exnid'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'LV Falcons'), (SELECT id FROM players WHERE gamertag = 'd7oom'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'LV Falcons'), (SELECT id FROM players WHERE gamertag = 'KiinG'), 6, true, '2024-11-01');

-- Insert cumulative player stats for Major 1 (players only, using OVERALL K, D, KD from your table)
INSERT INTO player_tournament_stats (player_id, team_id, tournament_id, total_kills, total_deaths, total_assists, total_damage, kd_ratio, kda_ratio) VALUES
-- LA Thieves
((SELECT id FROM players WHERE gamertag = 'Scrap'), (SELECT id FROM teams WHERE name = 'LA Thieves'), 1, 438, 370, 0, 0, 1.184, 1.184),
((SELECT id FROM players WHERE gamertag = 'Envoy'), (SELECT id FROM teams WHERE name = 'LA Thieves'), 1, 420, 369, 0, 0, 1.138, 1.138),
((SELECT id FROM players WHERE gamertag = 'HyDra'), (SELECT id FROM teams WHERE name = 'LA Thieves'), 1, 401, 378, 0, 0, 1.061, 1.061),
((SELECT id FROM players WHERE gamertag = 'Ghosty'), (SELECT id FROM teams WHERE name = 'LA Thieves'), 1, 386, 368, 0, 0, 1.049, 1.049),
-- CAR Royal Ravens
((SELECT id FROM players WHERE gamertag = 'Gwinn'), (SELECT id FROM teams WHERE name = 'CAR Royal Ravens'), 1, 383, 372, 0, 0, 1.030, 1.030),
((SELECT id FROM players WHERE gamertag = 'TJHaLy'), (SELECT id FROM teams WHERE name = 'CAR Royal Ravens'), 1, 398, 394, 0, 0, 1.010, 1.010),
((SELECT id FROM players WHERE gamertag = 'SlasheR'), (SELECT id FROM teams WHERE name = 'CAR Royal Ravens'), 1, 340, 346, 0, 0, 0.983, 0.983),
((SELECT id FROM players WHERE gamertag = 'Vivid'), (SELECT id FROM teams WHERE name = 'CAR Royal Ravens'), 1, 381, 418, 0, 0, 0.912, 0.912),
-- MIA Heretics
((SELECT id FROM players WHERE gamertag = 'ReeaL'), (SELECT id FROM teams WHERE name = 'MIA Heretics'), 1, 396, 394, 0, 0, 1.005, 1.005),
((SELECT id FROM players WHERE gamertag = 'RenKoR'), (SELECT id FROM teams WHERE name = 'MIA Heretics'), 1, 362, 387, 0, 0, 0.935, 0.935),
((SELECT id FROM players WHERE gamertag = 'MettalZ'), (SELECT id FROM teams WHERE name = 'MIA Heretics'), 1, 354, 418, 0, 0, 0.847, 0.847),
((SELECT id FROM players WHERE gamertag = 'Lucky'), (SELECT id FROM teams WHERE name = 'MIA Heretics'), 1, 353, 368, 0, 0, 0.959, 0.959),
-- ATL FaZe
((SELECT id FROM players WHERE gamertag = 'Simp'), (SELECT id FROM teams WHERE name = 'ATL FaZe'), 1, 336, 314, 0, 0, 1.070, 1.070),
((SELECT id FROM players WHERE gamertag = 'Drazah'), (SELECT id FROM teams WHERE name = 'ATL FaZe'), 1, 331, 318, 0, 0, 1.041, 1.041),
((SELECT id FROM players WHERE gamertag = 'aBeZy'), (SELECT id FROM teams WHERE name = 'ATL FaZe'), 1, 326, 329, 0, 0, 0.991, 0.991),
((SELECT id FROM players WHERE gamertag = 'Cellium'), (SELECT id FROM teams WHERE name = 'ATL FaZe'), 1, 320, 289, 0, 0, 1.107, 1.107),
-- MIN RØKKR
((SELECT id FROM players WHERE gamertag = 'Gio'), (SELECT id FROM teams WHERE name = 'MIN RØKKR'), 1, 305, 247, 0, 0, 1.235, 1.235),
((SELECT id FROM players WHERE gamertag = 'Nero'), (SELECT id FROM teams WHERE name = 'MIN RØKKR'), 1, 294, 270, 0, 0, 1.089, 1.089),
((SELECT id FROM players WHERE gamertag = 'Estreal'), (SELECT id FROM teams WHERE name = 'MIN RØKKR'), 1, 261, 259, 0, 0, 1.008, 1.008),
((SELECT id FROM players WHERE gamertag = 'PaulEhx'), (SELECT id FROM teams WHERE name = 'MIN RØKKR'), 1, 255, 271, 0, 0, 0.941, 0.941),
-- OpTic TEX
((SELECT id FROM players WHERE gamertag = 'Huke'), (SELECT id FROM teams WHERE name = 'OpTic TEX'), 1, 253, 235, 0, 0, 1.077, 1.077),
((SELECT id FROM players WHERE gamertag = 'Shotzzy'), (SELECT id FROM teams WHERE name = 'OpTic TEX'), 1, 250, 260, 0, 0, 0.962, 0.962),
((SELECT id FROM players WHERE gamertag = 'Kenny'), (SELECT id FROM teams WHERE name = 'OpTic TEX'), 1, 242, 274, 0, 0, 0.883, 0.883),
((SELECT id FROM players WHERE gamertag = 'Dashy'), (SELECT id FROM teams WHERE name = 'OpTic TEX'), 1, 240, 258, 0, 0, 0.930, 0.930),
-- TOR Ultra
((SELECT id FROM players WHERE gamertag = 'CleanX'), (SELECT id FROM teams WHERE name = 'TOR Ultra'), 1, 231, 221, 0, 0, 1.045, 1.045),
((SELECT id FROM players WHERE gamertag = 'JoeDeceives'), (SELECT id FROM teams WHERE name = 'TOR Ultra'), 1, 225, 221, 0, 0, 1.018, 1.018),
((SELECT id FROM players WHERE gamertag = 'Beans'), (SELECT id FROM teams WHERE name = 'TOR Ultra'), 1, 210, 191, 0, 0, 1.100, 1.100),
((SELECT id FROM players WHERE gamertag = 'Insight'), (SELECT id FROM teams WHERE name = 'TOR Ultra'), 1, 203, 196, 0, 0, 1.036, 1.036),
-- Cloud9 NY
((SELECT id FROM players WHERE gamertag = 'Sib'), (SELECT id FROM teams WHERE name = 'Cloud9 NY'), 1, 193, 172, 0, 0, 1.122, 1.122),
((SELECT id FROM players WHERE gamertag = 'Kremp'), (SELECT id FROM teams WHERE name = 'Cloud9 NY'), 1, 188, 182, 0, 0, 1.033, 1.033),
((SELECT id FROM players WHERE gamertag = 'Mack'), (SELECT id FROM teams WHERE name = 'Cloud9 NY'), 1, 183, 197, 0, 0, 0.929, 0.929),
((SELECT id FROM players WHERE gamertag = 'Attach'), (SELECT id FROM teams WHERE name = 'Cloud9 NY'), 1, 162, 168, 0, 0, 0.964, 0.964),
-- VAN Surge
((SELECT id FROM players WHERE gamertag = '04'), (SELECT id FROM teams WHERE name = 'VAN Surge'), 1, 147, 166, 0, 0, 0.886, 0.886),
((SELECT id FROM players WHERE gamertag = 'Abuzah'), (SELECT id FROM teams WHERE name = 'VAN Surge'), 1, 145, 160, 0, 0, 0.906, 0.906),
((SELECT id FROM players WHERE gamertag = 'Hicksy'), (SELECT id FROM teams WHERE name = 'VAN Surge'), 1, 114, 142, 0, 0, 0.803, 0.803),
((SELECT id FROM players WHERE gamertag = 'Nastie'), (SELECT id FROM teams WHERE name = 'VAN Surge'), 1, 129, 146, 0, 0, 0.884, 0.884),
-- BOS Breach
((SELECT id FROM players WHERE gamertag = 'Snoopy'), (SELECT id FROM teams WHERE name = 'BOS Breach'), 1, 142, 138, 0, 0, 1.029, 1.029),
((SELECT id FROM players WHERE gamertag = 'Cammy'), (SELECT id FROM teams WHERE name = 'BOS Breach'), 1, 113, 152, 0, 0, 0.743, 0.743),
((SELECT id FROM players WHERE gamertag = 'Owakening'), (SELECT id FROM teams WHERE name = 'BOS Breach'), 1, 132, 137, 0, 0, 0.964, 0.964),
((SELECT id FROM players WHERE gamertag = 'Purj'), (SELECT id FROM teams WHERE name = 'BOS Breach'), 1, 124, 165, 0, 0, 0.752, 0.752),
-- LA Guerrillas M8
((SELECT id FROM players WHERE gamertag = 'Lynz'), (SELECT id FROM teams WHERE name = 'LA Guerrillas M8'), 1, 89, 79, 0, 0, 1.127, 1.127),
((SELECT id FROM players WHERE gamertag = 'KiSMET'), (SELECT id FROM teams WHERE name = 'LA Guerrillas M8'), 1, 94, 85, 0, 0, 1.106, 1.106),
((SELECT id FROM players WHERE gamertag = 'Skyz'), (SELECT id FROM teams WHERE name = 'LA Guerrillas M8'), 1, 74, 78, 0, 0, 0.949, 0.949),
((SELECT id FROM players WHERE gamertag = 'Priestahh'), (SELECT id FROM teams WHERE name = 'LA Guerrillas M8'), 1, 84, 83, 0, 0, 1.012, 1.012),
-- LV Falcons
((SELECT id FROM players WHERE gamertag = 'KiinG'), (SELECT id FROM teams WHERE name = 'LV Falcons'), 1, 72, 70, 0, 0, 1.029, 1.029),
((SELECT id FROM players WHERE gamertag = 'Exnid'), (SELECT id FROM teams WHERE name = 'LV Falcons'), 1, 63, 62, 0, 0, 1.016, 1.016),
((SELECT id FROM players WHERE gamertag = 'Roxas'), (SELECT id FROM teams WHERE name = 'LV Falcons'), 1, 51, 64, 0, 0, 0.797, 0.797),
((SELECT id FROM players WHERE gamertag = 'd7oom'), (SELECT id FROM teams WHERE name = 'LV Falcons'), 1, 46, 67, 0, 0, 0.687, 0.687)
ON CONFLICT DO NOTHING; 