-- Major 2 Data for CDL 2025
-- Assumes schema.sql has already been run

-- Insert Major 2 Teams
INSERT INTO teams (name, abbreviation, city, is_active) VALUES
('LA Thieves', 'LAT', 'Los Angeles', true),
('ATL FaZe', 'ATL', 'Atlanta', true),
('TOR Ultra', 'TOR', 'Toronto', true),
('CAR Royal Ravens', 'CAR', 'Charlotte', true),
('VAN Surge', 'VAN', 'Vancouver', true),
('Cloud9 NY', 'NYC', 'New York', true),
('MIN RØKKR', 'MIN', 'Minneapolis', true),
('BOS Breach', 'BOS', 'Boston', true),
('MIA Heretics', 'MIA', 'Miami', true),
('OpTic TEX', 'TEX', 'Dallas', true),
('LA Guerrillas M8', 'LAG', 'Los Angeles', true),
('LV Falcons', 'LVF', 'Las Vegas', true)
ON CONFLICT DO NOTHING;

-- Insert Major 2 Coaches
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
('Saintt', (SELECT id FROM teams WHERE name = 'CAR Royal Ravens'), 6),
('Rambo', (SELECT id FROM teams WHERE name = 'VAN Surge'), 6),
('Accuracy', (SELECT id FROM teams WHERE name = 'Cloud9 NY'), 6),
('Arian', (SELECT id FROM teams WHERE name = 'Cloud9 NY'), 6),
('Loony', (SELECT id FROM teams WHERE name = 'MIN RØKKR'), 6),
('Alexdotzip', (SELECT id FROM teams WHERE name = 'MIN RØKKR'), 6),
('Seany', (SELECT id FROM teams WHERE name = 'BOS Breach'), 6),
('Magxck', (SELECT id FROM teams WHERE name = 'BOS Breach'), 6),
('MethodZ', (SELECT id FROM teams WHERE name = 'MIA Heretics'), 6),
('Sikotik', (SELECT id FROM teams WHERE name = 'MIA Heretics'), 6),
('Karma', (SELECT id FROM teams WHERE name = 'OpTic TEX'), 6),
('JP Krez', (SELECT id FROM teams WHERE name = 'OpTic TEX'), 6),
('DREAL', (SELECT id FROM teams WHERE name = 'LA Guerrillas M8'), 6),
('MarkyB', (SELECT id FROM teams WHERE name = 'LA Guerrillas M8'), 6),
('Clayster', (SELECT id FROM teams WHERE name = 'LV Falcons'), 6);

-- Insert Major 2 Players (excluding coaches)
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
-- CAR Royal Ravens
('Gwinn', true),
('TJHaLy', true),
('SlasheR', true),
('Vivid', true),
-- VAN Surge
('Abuzah', true),
('04', true),
('Nastie', true),
('Neptune', true),
-- Cloud9 NY
('Sib', true),
('Attach', true),
('Kremp', true),
('Mack', true),
-- MIN RØKKR
('Nero', true),
('Gio', true),
('Estreal', true),
('PaulEhx', true),
-- BOS Breach
('Snoopy', true),
('Cammy', true),
('Owakening', true),
('Purj', true),
-- MIA Heretics
('Lucky', true),
('MettalZ', true),
('ReeaL', true),
('RenKoR', true),
-- OpTic TEX
('Shotzzy', true),
('Dashy', true),
('Pred', true),
('Skyz', true),
-- LA Guerrillas M8
('Lynz', true),
('KiSMET', true),
('Priestahh', true),
('Lunarz', true),
-- LV Falcons
('Exnid', true),
('d7oom', true),
('KiinG', true),
('WXSL', true)
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
-- CAR Royal Ravens
((SELECT id FROM teams WHERE name = 'CAR Royal Ravens'), (SELECT id FROM players WHERE gamertag = 'Gwinn'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'CAR Royal Ravens'), (SELECT id FROM players WHERE gamertag = 'TJHaLy'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'CAR Royal Ravens'), (SELECT id FROM players WHERE gamertag = 'SlasheR'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'CAR Royal Ravens'), (SELECT id FROM players WHERE gamertag = 'Vivid'), 6, true, '2024-11-01'),
-- VAN Surge
((SELECT id FROM teams WHERE name = 'VAN Surge'), (SELECT id FROM players WHERE gamertag = 'Abuzah'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'VAN Surge'), (SELECT id FROM players WHERE gamertag = '04'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'VAN Surge'), (SELECT id FROM players WHERE gamertag = 'Nastie'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'VAN Surge'), (SELECT id FROM players WHERE gamertag = 'Neptune'), 6, true, '2024-11-01'),
-- Cloud9 NY
((SELECT id FROM teams WHERE name = 'Cloud9 NY'), (SELECT id FROM players WHERE gamertag = 'Sib'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'Cloud9 NY'), (SELECT id FROM players WHERE gamertag = 'Attach'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'Cloud9 NY'), (SELECT id FROM players WHERE gamertag = 'Kremp'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'Cloud9 NY'), (SELECT id FROM players WHERE gamertag = 'Mack'), 6, true, '2024-11-01'),
-- MIN RØKKR
((SELECT id FROM teams WHERE name = 'MIN RØKKR'), (SELECT id FROM players WHERE gamertag = 'Nero'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'MIN RØKKR'), (SELECT id FROM players WHERE gamertag = 'Gio'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'MIN RØKKR'), (SELECT id FROM players WHERE gamertag = 'Estreal'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'MIN RØKKR'), (SELECT id FROM players WHERE gamertag = 'PaulEhx'), 6, true, '2024-11-01'),
-- BOS Breach
((SELECT id FROM teams WHERE name = 'BOS Breach'), (SELECT id FROM players WHERE gamertag = 'Snoopy'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'BOS Breach'), (SELECT id FROM players WHERE gamertag = 'Cammy'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'BOS Breach'), (SELECT id FROM players WHERE gamertag = 'Owakening'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'BOS Breach'), (SELECT id FROM players WHERE gamertag = 'Purj'), 6, true, '2024-11-01'),
-- MIA Heretics
((SELECT id FROM teams WHERE name = 'MIA Heretics'), (SELECT id FROM players WHERE gamertag = 'Lucky'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'MIA Heretics'), (SELECT id FROM players WHERE gamertag = 'MettalZ'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'MIA Heretics'), (SELECT id FROM players WHERE gamertag = 'ReeaL'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'MIA Heretics'), (SELECT id FROM players WHERE gamertag = 'RenKoR'), 6, true, '2024-11-01'),
-- OpTic TEX
((SELECT id FROM teams WHERE name = 'OpTic TEX'), (SELECT id FROM players WHERE gamertag = 'Shotzzy'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'OpTic TEX'), (SELECT id FROM players WHERE gamertag = 'Dashy'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'OpTic TEX'), (SELECT id FROM players WHERE gamertag = 'Pred'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'OpTic TEX'), (SELECT id FROM players WHERE gamertag = 'Skyz'), 6, true, '2024-11-01'),
-- LA Guerrillas M8
((SELECT id FROM teams WHERE name = 'LA Guerrillas M8'), (SELECT id FROM players WHERE gamertag = 'Lynz'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'LA Guerrillas M8'), (SELECT id FROM players WHERE gamertag = 'KiSMET'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'LA Guerrillas M8'), (SELECT id FROM players WHERE gamertag = 'Priestahh'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'LA Guerrillas M8'), (SELECT id FROM players WHERE gamertag = 'Lunarz'), 6, true, '2024-11-01'),
-- LV Falcons
((SELECT id FROM teams WHERE name = 'LV Falcons'), (SELECT id FROM players WHERE gamertag = 'Exnid'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'LV Falcons'), (SELECT id FROM players WHERE gamertag = 'd7oom'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'LV Falcons'), (SELECT id FROM players WHERE gamertag = 'KiinG'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'LV Falcons'), (SELECT id FROM players WHERE gamertag = 'WXSL'), 6, true, '2024-11-01');

-- Insert cumulative player stats for Major 2 (players only, using OVERALL K, D, KD from your table)
INSERT INTO player_tournament_stats (player_id, team_id, tournament_id, total_kills, total_deaths, total_assists, total_damage, kd_ratio, kda_ratio) VALUES
-- LA Thieves
((SELECT id FROM players WHERE gamertag = 'Ghosty'), (SELECT id FROM teams WHERE name = 'LA Thieves'), 2, 294, 303, 0, 0, 0.970, 0.970),
((SELECT id FROM players WHERE gamertag = 'Scrap'), (SELECT id FROM teams WHERE name = 'LA Thieves'), 2, 356, 326, 0, 0, 1.092, 1.092),
((SELECT id FROM players WHERE gamertag = 'HyDra'), (SELECT id FROM teams WHERE name = 'LA Thieves'), 2, 355, 336, 0, 0, 1.057, 1.057),
((SELECT id FROM players WHERE gamertag = 'Envoy'), (SELECT id FROM teams WHERE name = 'LA Thieves'), 2, 323, 352, 0, 0, 0.918, 0.918),
-- ATL FaZe
((SELECT id FROM players WHERE gamertag = 'aBeZy'), (SELECT id FROM teams WHERE name = 'ATL FaZe'), 2, 289, 307, 0, 0, 0.941, 0.941),
((SELECT id FROM players WHERE gamertag = 'Cellium'), (SELECT id FROM teams WHERE name = 'ATL FaZe'), 2, 324, 250, 0, 0, 1.296, 1.296),
((SELECT id FROM players WHERE gamertag = 'Simp'), (SELECT id FROM teams WHERE name = 'ATL FaZe'), 2, 319, 300, 0, 0, 1.063, 1.063),
((SELECT id FROM players WHERE gamertag = 'Drazah'), (SELECT id FROM teams WHERE name = 'ATL FaZe'), 2, 310, 280, 0, 0, 1.107, 1.107),
-- TOR Ultra
((SELECT id FROM players WHERE gamertag = 'CleanX'), (SELECT id FROM teams WHERE name = 'TOR Ultra'), 2, 381, 368, 0, 0, 1.035, 1.035),
((SELECT id FROM players WHERE gamertag = 'Insight'), (SELECT id FROM teams WHERE name = 'TOR Ultra'), 2, 315, 315, 0, 0, 1.000, 1.000),
((SELECT id FROM players WHERE gamertag = 'Beans'), (SELECT id FROM teams WHERE name = 'TOR Ultra'), 2, 355, 325, 0, 0, 1.092, 1.092),
((SELECT id FROM players WHERE gamertag = 'JoeDeceives'), (SELECT id FROM teams WHERE name = 'TOR Ultra'), 2, 399, 346, 0, 0, 1.153, 1.153),
-- CAR Royal Ravens
((SELECT id FROM players WHERE gamertag = 'Gwinn'), (SELECT id FROM teams WHERE name = 'CAR Royal Ravens'), 2, 154, 156, 0, 0, 0.987, 0.987),
((SELECT id FROM players WHERE gamertag = 'TJHaLy'), (SELECT id FROM teams WHERE name = 'CAR Royal Ravens'), 2, 172, 183, 0, 0, 0.940, 0.940),
((SELECT id FROM players WHERE gamertag = 'SlasheR'), (SELECT id FROM teams WHERE name = 'CAR Royal Ravens'), 2, 145, 148, 0, 0, 0.980, 0.980),
((SELECT id FROM players WHERE gamertag = 'Vivid'), (SELECT id FROM teams WHERE name = 'CAR Royal Ravens'), 2, 146, 170, 0, 0, 0.859, 0.859),
-- VAN Surge
((SELECT id FROM players WHERE gamertag = 'Abuzah'), (SELECT id FROM teams WHERE name = 'VAN Surge'), 2, 390, 409, 0, 0, 0.954, 0.954),
((SELECT id FROM players WHERE gamertag = '04'), (SELECT id FROM teams WHERE name = 'VAN Surge'), 2, 424, 420, 0, 0, 1.010, 1.010),
((SELECT id FROM players WHERE gamertag = 'Nastie'), (SELECT id FROM teams WHERE name = 'VAN Surge'), 2, 389, 357, 0, 0, 1.090, 1.090),
((SELECT id FROM players WHERE gamertag = 'Neptune'), (SELECT id FROM teams WHERE name = 'VAN Surge'), 2, 459, 418, 0, 0, 1.098, 1.098),
-- Cloud9 NY
((SELECT id FROM players WHERE gamertag = 'Sib'), (SELECT id FROM teams WHERE name = 'Cloud9 NY'), 2, 202, 204, 0, 0, 0.990, 0.990),
((SELECT id FROM players WHERE gamertag = 'Attach'), (SELECT id FROM teams WHERE name = 'Cloud9 NY'), 2, 161, 188, 0, 0, 0.856, 0.856),
((SELECT id FROM players WHERE gamertag = 'Kremp'), (SELECT id FROM teams WHERE name = 'Cloud9 NY'), 2, 215, 197, 0, 0, 1.091, 1.091),
((SELECT id FROM players WHERE gamertag = 'Mack'), (SELECT id FROM teams WHERE name = 'Cloud9 NY'), 2, 199, 211, 0, 0, 0.943, 0.943),
-- MIN RØKKR
((SELECT id FROM players WHERE gamertag = 'Nero'), (SELECT id FROM teams WHERE name = 'MIN RØKKR'), 2, 80, 114, 0, 0, 0.702, 0.702),
((SELECT id FROM players WHERE gamertag = 'Gio'), (SELECT id FROM teams WHERE name = 'MIN RØKKR'), 2, 120, 98, 0, 0, 1.225, 1.225),
((SELECT id FROM players WHERE gamertag = 'Estreal'), (SELECT id FROM teams WHERE name = 'MIN RØKKR'), 2, 106, 116, 0, 0, 0.914, 0.914),
((SELECT id FROM players WHERE gamertag = 'PaulEhx'), (SELECT id FROM teams WHERE name = 'MIN RØKKR'), 2, 66, 114, 0, 0, 0.579, 0.579),
-- BOS Breach
((SELECT id FROM players WHERE gamertag = 'Snoopy'), (SELECT id FROM teams WHERE name = 'BOS Breach'), 2, 244, 244, 0, 0, 1.000, 1.000),
((SELECT id FROM players WHERE gamertag = 'Cammy'), (SELECT id FROM teams WHERE name = 'BOS Breach'), 2, 240, 249, 0, 0, 0.964, 0.964),
((SELECT id FROM players WHERE gamertag = 'Owakening'), (SELECT id FROM teams WHERE name = 'BOS Breach'), 2, 228, 241, 0, 0, 0.946, 0.946),
((SELECT id FROM players WHERE gamertag = 'Purj'), (SELECT id FROM teams WHERE name = 'BOS Breach'), 2, 218, 258, 0, 0, 0.845, 0.845),
-- MIA Heretics
((SELECT id FROM players WHERE gamertag = 'Lucky'), (SELECT id FROM teams WHERE name = 'MIA Heretics'), 2, 68, 92, 0, 0, 0.739, 0.739),
((SELECT id FROM players WHERE gamertag = 'MettalZ'), (SELECT id FROM teams WHERE name = 'MIA Heretics'), 2, 97, 94, 0, 0, 1.032, 1.032),
((SELECT id FROM players WHERE gamertag = 'ReeaL'), (SELECT id FROM teams WHERE name = 'MIA Heretics'), 2, 117, 89, 0, 0, 1.315, 1.315),
((SELECT id FROM players WHERE gamertag = 'RenKoR'), (SELECT id FROM teams WHERE name = 'MIA Heretics'), 2, 96, 84, 0, 0, 1.143, 1.143),
-- OpTic TEX
((SELECT id FROM players WHERE gamertag = 'Shotzzy'), (SELECT id FROM teams WHERE name = 'OpTic TEX'), 2, 107, 100, 0, 0, 1.070, 1.070),
((SELECT id FROM players WHERE gamertag = 'Dashy'), (SELECT id FROM teams WHERE name = 'OpTic TEX'), 2, 108, 96, 0, 0, 1.125, 1.125),
((SELECT id FROM players WHERE gamertag = 'Pred'), (SELECT id FROM teams WHERE name = 'OpTic TEX'), 2, 116, 107, 0, 0, 1.084, 1.084),
((SELECT id FROM players WHERE gamertag = 'Skyz'), (SELECT id FROM teams WHERE name = 'OpTic TEX'), 2, 86, 113, 0, 0, 0.761, 0.761),
-- LA Guerrillas M8
((SELECT id FROM players WHERE gamertag = 'Lynz'), (SELECT id FROM teams WHERE name = 'LA Guerrillas M8'), 2, 39, 51, 0, 0, 0.765, 0.765),
((SELECT id FROM players WHERE gamertag = 'KiSMET'), (SELECT id FROM teams WHERE name = 'LA Guerrillas M8'), 2, 47, 56, 0, 0, 0.839, 0.839),
((SELECT id FROM players WHERE gamertag = 'Priestahh'), (SELECT id FROM teams WHERE name = 'LA Guerrillas M8'), 2, 35, 56, 0, 0, 0.625, 0.625),
((SELECT id FROM players WHERE gamertag = 'Lunarz'), (SELECT id FROM teams WHERE name = 'LA Guerrillas M8'), 2, 47, 57, 0, 0, 0.825, 0.825),
-- LV Falcons
((SELECT id FROM players WHERE gamertag = 'Exnid'), (SELECT id FROM teams WHERE name = 'LV Falcons'), 2, 44, 53, 0, 0, 0.830, 0.830),
((SELECT id FROM players WHERE gamertag = 'd7oom'), (SELECT id FROM teams WHERE name = 'LV Falcons'), 2, 37, 49, 0, 0, 0.755, 0.755),
((SELECT id FROM players WHERE gamertag = 'KiinG'), (SELECT id FROM teams WHERE name = 'LV Falcons'), 2, 43, 49, 0, 0, 0.878, 0.878),
((SELECT id FROM players WHERE gamertag = 'WXSL'), (SELECT id FROM teams WHERE name = 'LV Falcons'), 2, 33, 52, 0, 0, 0.635, 0.635)
ON CONFLICT DO NOTHING; 