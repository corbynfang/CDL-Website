-- Major 3 Data for CDL 2025
-- Assumes schema.sql has already been run

-- Insert Major 3 Teams
INSERT INTO teams (name, abbreviation, city, is_active) VALUES
('LA Thieves', 'LAT', 'Los Angeles', true),
('BOS Breach', 'BOS', 'Boston', true),
('VAN Surge', 'VAN', 'Vancouver', true),
('MIA Heretics', 'MIA', 'Miami', true),
('TOR Ultra', 'TOR', 'Toronto', true),
('ATL FaZe', 'ATL', 'Atlanta', true),
('OpTic TEX', 'TEX', 'Dallas', true),
('CAR Royal Ravens', 'CAR', 'Charlotte', true),
('MIN RØKKR', 'MIN', 'Minneapolis', true),
('LV Falcons', 'LVF', 'Las Vegas', true),
('Cloud9 NY', 'NYC', 'New York', true),
('LA Guerrillas M8', 'LAG', 'Los Angeles', true)
ON CONFLICT DO NOTHING;

-- Insert Major 3 Coaches
CREATE TABLE IF NOT EXISTS coaches (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    team_id INTEGER REFERENCES teams(id),
    season_id INTEGER REFERENCES seasons(id)
);

INSERT INTO coaches (name, team_id, season_id) VALUES
('ShAnE', (SELECT id FROM teams WHERE name = 'LA Thieves'), 6),
('Sender', (SELECT id FROM teams WHERE name = 'LA Thieves'), 6),
('Seany', (SELECT id FROM teams WHERE name = 'BOS Breach'), 6),
('Magxck', (SELECT id FROM teams WHERE name = 'BOS Breach'), 6),
('Rambo', (SELECT id FROM teams WHERE name = 'VAN Surge'), 6),
('MethodZ', (SELECT id FROM teams WHERE name = 'MIA Heretics'), 6),
('Sikotik', (SELECT id FROM teams WHERE name = 'MIA Heretics'), 6),
('Flux', (SELECT id FROM teams WHERE name = 'TOR Ultra'), 6),
('Joee', (SELECT id FROM teams WHERE name = 'TOR Ultra'), 6),
('Crowder', (SELECT id FROM teams WHERE name = 'ATL FaZe'), 6),
('RJ', (SELECT id FROM teams WHERE name = 'ATL FaZe'), 6),
('Karma', (SELECT id FROM teams WHERE name = 'OpTic TEX'), 6),
('JPKrez', (SELECT id FROM teams WHERE name = 'OpTic TEX'), 6),
('Saintt', (SELECT id FROM teams WHERE name = 'CAR Royal Ravens'), 6),
('Loony', (SELECT id FROM teams WHERE name = 'MIN RØKKR'), 6),
('Alexdotzip', (SELECT id FROM teams WHERE name = 'MIN RØKKR'), 6),
('Clayster', (SELECT id FROM teams WHERE name = 'LV Falcons'), 6),
('Accuracy', (SELECT id FROM teams WHERE name = 'Cloud9 NY'), 6),
('Arian', (SELECT id FROM teams WHERE name = 'Cloud9 NY'), 6),
('DREAL', (SELECT id FROM teams WHERE name = 'LA Guerrillas M8'), 6),
('MarkyB', (SELECT id FROM teams WHERE name = 'LA Guerrillas M8'), 6);

-- Insert Major 3 Players (excluding coaches)
INSERT INTO players (gamertag, is_active) VALUES
-- LA Thieves
('Ghosty', true),
('Scrap', true),
('HyDra', true),
('Envoy', true),
-- BOS Breach
('Snoopy', true),
('Cammy', true),
('Owakening', true),
('Purj', true),
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
-- TOR Ultra
('CleanX', true),
('Insight', true),
('JoeDeceives', true),
('Mercules', true),
-- ATL FaZe
('aBeZy', true),
('Cellium', true),
('Simp', true),
('Drazah', true),
-- OpTic TEX
('Shotzzy', true),
('Dashy', true),
('Skyz', true),
('Huke', true),
-- CAR Royal Ravens
('Gwinn', true),
('TJHaLy', true),
('SlasheR', true),
('Wrecks', true),
-- MIN RØKKR
('Nero', true),
('Gio', true),
('Estreal', true),
('Kenny', true),
-- LV Falcons
('Exnid', true),
('Pred', true),
('Arcitys', true),
('Priestahh', true),
-- Cloud9 NY
('Sib', true),
('Attach', true),
('Kremp', true),
('Mack', true),
-- LA Guerrillas M8
('KiSMET', true),
('Lunarz', true),
('FeLo', true),
('oJohnny', true)
ON CONFLICT DO NOTHING;

-- Assign players to teams for 2025 season (season_id = 6)
INSERT INTO team_rosters (team_id, player_id, season_id, is_starter, start_date) VALUES
-- LA Thieves
((SELECT id FROM teams WHERE name = 'LA Thieves'), (SELECT id FROM players WHERE gamertag = 'Ghosty'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'LA Thieves'), (SELECT id FROM players WHERE gamertag = 'Scrap'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'LA Thieves'), (SELECT id FROM players WHERE gamertag = 'HyDra'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'LA Thieves'), (SELECT id FROM players WHERE gamertag = 'Envoy'), 6, true, '2024-11-01'),
-- BOS Breach
((SELECT id FROM teams WHERE name = 'BOS Breach'), (SELECT id FROM players WHERE gamertag = 'Snoopy'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'BOS Breach'), (SELECT id FROM players WHERE gamertag = 'Cammy'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'BOS Breach'), (SELECT id FROM players WHERE gamertag = 'Owakening'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'BOS Breach'), (SELECT id FROM players WHERE gamertag = 'Purj'), 6, true, '2024-11-01'),
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
-- TOR Ultra
((SELECT id FROM teams WHERE name = 'TOR Ultra'), (SELECT id FROM players WHERE gamertag = 'CleanX'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'TOR Ultra'), (SELECT id FROM players WHERE gamertag = 'Insight'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'TOR Ultra'), (SELECT id FROM players WHERE gamertag = 'JoeDeceives'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'TOR Ultra'), (SELECT id FROM players WHERE gamertag = 'Mercules'), 6, true, '2024-11-01'),
-- ATL FaZe
((SELECT id FROM teams WHERE name = 'ATL FaZe'), (SELECT id FROM players WHERE gamertag = 'aBeZy'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'ATL FaZe'), (SELECT id FROM players WHERE gamertag = 'Cellium'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'ATL FaZe'), (SELECT id FROM players WHERE gamertag = 'Simp'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'ATL FaZe'), (SELECT id FROM players WHERE gamertag = 'Drazah'), 6, true, '2024-11-01'),
-- OpTic TEX
((SELECT id FROM teams WHERE name = 'OpTic TEX'), (SELECT id FROM players WHERE gamertag = 'Shotzzy'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'OpTic TEX'), (SELECT id FROM players WHERE gamertag = 'Dashy'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'OpTic TEX'), (SELECT id FROM players WHERE gamertag = 'Skyz'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'OpTic TEX'), (SELECT id FROM players WHERE gamertag = 'Huke'), 6, true, '2024-11-01'),
-- CAR Royal Ravens
((SELECT id FROM teams WHERE name = 'CAR Royal Ravens'), (SELECT id FROM players WHERE gamertag = 'Gwinn'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'CAR Royal Ravens'), (SELECT id FROM players WHERE gamertag = 'TJHaLy'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'CAR Royal Ravens'), (SELECT id FROM players WHERE gamertag = 'SlasheR'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'CAR Royal Ravens'), (SELECT id FROM players WHERE gamertag = 'Wrecks'), 6, true, '2024-11-01'),
-- MIN RØKKR
((SELECT id FROM teams WHERE name = 'MIN RØKKR'), (SELECT id FROM players WHERE gamertag = 'Nero'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'MIN RØKKR'), (SELECT id FROM players WHERE gamertag = 'Gio'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'MIN RØKKR'), (SELECT id FROM players WHERE gamertag = 'Estreal'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'MIN RØKKR'), (SELECT id FROM players WHERE gamertag = 'Kenny'), 6, true, '2024-11-01'),
-- LV Falcons
((SELECT id FROM teams WHERE name = 'LV Falcons'), (SELECT id FROM players WHERE gamertag = 'Exnid'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'LV Falcons'), (SELECT id FROM players WHERE gamertag = 'Pred'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'LV Falcons'), (SELECT id FROM players WHERE gamertag = 'Arcitys'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'LV Falcons'), (SELECT id FROM players WHERE gamertag = 'Priestahh'), 6, true, '2024-11-01'),
-- Cloud9 NY
((SELECT id FROM teams WHERE name = 'Cloud9 NY'), (SELECT id FROM players WHERE gamertag = 'Sib'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'Cloud9 NY'), (SELECT id FROM players WHERE gamertag = 'Attach'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'Cloud9 NY'), (SELECT id FROM players WHERE gamertag = 'Kremp'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'Cloud9 NY'), (SELECT id FROM players WHERE gamertag = 'Mack'), 6, true, '2024-11-01'),
-- LA Guerrillas M8
((SELECT id FROM teams WHERE name = 'LA Guerrillas M8'), (SELECT id FROM players WHERE gamertag = 'KiSMET'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'LA Guerrillas M8'), (SELECT id FROM players WHERE gamertag = 'Lunarz'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'LA Guerrillas M8'), (SELECT id FROM players WHERE gamertag = 'FeLo'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'LA Guerrillas M8'), (SELECT id FROM players WHERE gamertag = 'oJohnny'), 6, true, '2024-11-01');

-- Insert cumulative player stats for Major 3 (players only, using OVERALL K, D, KD from your table)
INSERT INTO player_tournament_stats (player_id, team_id, tournament_id, total_kills, total_deaths, total_assists, total_damage, kd_ratio, kda_ratio) VALUES
-- LA Thieves
((SELECT id FROM players WHERE gamertag = 'Ghosty'), (SELECT id FROM teams WHERE name = 'LA Thieves'), 3, 239, 244, 0, 0, 0.980, 0.980),
((SELECT id FROM players WHERE gamertag = 'Scrap'), (SELECT id FROM teams WHERE name = 'LA Thieves'), 3, 299, 245, 0, 0, 1.220, 1.220),
((SELECT id FROM players WHERE gamertag = 'HyDra'), (SELECT id FROM teams WHERE name = 'LA Thieves'), 3, 311, 245, 0, 0, 1.269, 1.269),
((SELECT id FROM players WHERE gamertag = 'Envoy'), (SELECT id FROM teams WHERE name = 'LA Thieves'), 3, 272, 261, 0, 0, 1.042, 1.042),
-- BOS Breach
((SELECT id FROM players WHERE gamertag = 'Snoopy'), (SELECT id FROM teams WHERE name = 'BOS Breach'), 3, 195, 218, 0, 0, 0.895, 0.895),
((SELECT id FROM players WHERE gamertag = 'Cammy'), (SELECT id FROM teams WHERE name = 'BOS Breach'), 3, 181, 221, 0, 0, 0.819, 0.819),
((SELECT id FROM players WHERE gamertag = 'Owakening'), (SELECT id FROM teams WHERE name = 'BOS Breach'), 3, 197, 212, 0, 0, 0.929, 0.929),
((SELECT id FROM players WHERE gamertag = 'Purj'), (SELECT id FROM teams WHERE name = 'BOS Breach'), 3, 218, 220, 0, 0, 0.991, 0.991),
-- VAN Surge
((SELECT id FROM players WHERE gamertag = 'Abuzah'), (SELECT id FROM teams WHERE name = 'VAN Surge'), 3, 381, 370, 0, 0, 1.030, 1.030),
((SELECT id FROM players WHERE gamertag = '04'), (SELECT id FROM teams WHERE name = 'VAN Surge'), 3, 366, 412, 0, 0, 0.888, 0.888),
((SELECT id FROM players WHERE gamertag = 'Nastie'), (SELECT id FROM teams WHERE name = 'VAN Surge'), 3, 382, 353, 0, 0, 1.082, 1.082),
((SELECT id FROM players WHERE gamertag = 'Neptune'), (SELECT id FROM teams WHERE name = 'VAN Surge'), 3, 411, 408, 0, 0, 1.007, 1.007),
-- MIA Heretics
((SELECT id FROM players WHERE gamertag = 'MettalZ'), (SELECT id FROM teams WHERE name = 'MIA Heretics'), 3, 318, 321, 0, 0, 0.991, 0.991),
((SELECT id FROM players WHERE gamertag = 'ReeaL'), (SELECT id FROM teams WHERE name = 'MIA Heretics'), 3, 292, 299, 0, 0, 0.977, 0.977),
((SELECT id FROM players WHERE gamertag = 'RenKoR'), (SELECT id FROM teams WHERE name = 'MIA Heretics'), 3, 296, 293, 0, 0, 1.010, 1.010),
((SELECT id FROM players WHERE gamertag = 'SupeR'), (SELECT id FROM teams WHERE name = 'MIA Heretics'), 3, 280, 272, 0, 0, 1.029, 1.029),
-- TOR Ultra
((SELECT id FROM players WHERE gamertag = 'CleanX'), (SELECT id FROM teams WHERE name = 'TOR Ultra'), 3, 483, 490, 0, 0, 0.986, 0.986),
((SELECT id FROM players WHERE gamertag = 'Insight'), (SELECT id FROM teams WHERE name = 'TOR Ultra'), 3, 426, 418, 0, 0, 1.019, 1.019),
((SELECT id FROM players WHERE gamertag = 'JoeDeceives'), (SELECT id FROM teams WHERE name = 'TOR Ultra'), 3, 507, 454, 0, 0, 1.117, 1.117),
((SELECT id FROM players WHERE gamertag = 'Mercules'), (SELECT id FROM teams WHERE name = 'TOR Ultra'), 3, 474, 440, 0, 0, 1.077, 1.077),
-- ATL FaZe
((SELECT id FROM players WHERE gamertag = 'aBeZy'), (SELECT id FROM teams WHERE name = 'ATL FaZe'), 3, 258, 256, 0, 0, 1.008, 1.008),
((SELECT id FROM players WHERE gamertag = 'Cellium'), (SELECT id FROM teams WHERE name = 'ATL FaZe'), 3, 256, 228, 0, 0, 1.123, 1.123),
((SELECT id FROM players WHERE gamertag = 'Simp'), (SELECT id FROM teams WHERE name = 'ATL FaZe'), 3, 276, 255, 0, 0, 1.082, 1.082),
((SELECT id FROM players WHERE gamertag = 'Drazah'), (SELECT id FROM teams WHERE name = 'ATL FaZe'), 3, 242, 258, 0, 0, 0.938, 0.938),
-- OpTic TEX
((SELECT id FROM players WHERE gamertag = 'Shotzzy'), (SELECT id FROM teams WHERE name = 'OpTic TEX'), 3, 201, 183, 0, 0, 1.098, 1.098),
((SELECT id FROM players WHERE gamertag = 'Dashy'), (SELECT id FROM teams WHERE name = 'OpTic TEX'), 3, 190, 177, 0, 0, 1.073, 1.073),
((SELECT id FROM players WHERE gamertag = 'Skyz'), (SELECT id FROM teams WHERE name = 'OpTic TEX'), 3, 146, 177, 0, 0, 0.825, 0.825),
((SELECT id FROM players WHERE gamertag = 'Huke'), (SELECT id FROM teams WHERE name = 'OpTic TEX'), 3, 182, 177, 0, 0, 1.028, 1.028),
-- CAR Royal Ravens
((SELECT id FROM players WHERE gamertag = 'Gwinn'), (SELECT id FROM teams WHERE name = 'CAR Royal Ravens'), 3, 136, 136, 0, 0, 1.000, 1.000),
((SELECT id FROM players WHERE gamertag = 'TJHaLy'), (SELECT id FROM teams WHERE name = 'CAR Royal Ravens'), 3, 121, 146, 0, 0, 0.829, 0.829),
((SELECT id FROM players WHERE gamertag = 'SlasheR'), (SELECT id FROM teams WHERE name = 'CAR Royal Ravens'), 3, 111, 121, 0, 0, 0.917, 0.917),
((SELECT id FROM players WHERE gamertag = 'Wrecks'), (SELECT id FROM teams WHERE name = 'CAR Royal Ravens'), 3, 110, 143, 0, 0, 0.769, 0.769),
-- MIN RØKKR
((SELECT id FROM players WHERE gamertag = 'Nero'), (SELECT id FROM teams WHERE name = 'MIN RØKKR'), 3, 166, 155, 0, 0, 1.071, 1.071),
((SELECT id FROM players WHERE gamertag = 'Gio'), (SELECT id FROM teams WHERE name = 'MIN RØKKR'), 3, 163, 155, 0, 0, 1.052, 1.052),
((SELECT id FROM players WHERE gamertag = 'Estreal'), (SELECT id FROM teams WHERE name = 'MIN RØKKR'), 3, 156, 159, 0, 0, 0.981, 0.981),
((SELECT id FROM players WHERE gamertag = 'Kenny'), (SELECT id FROM teams WHERE name = 'MIN RØKKR'), 3, 158, 165, 0, 0, 0.958, 0.958),
-- LV Falcons
((SELECT id FROM players WHERE gamertag = 'Exnid'), (SELECT id FROM teams WHERE name = 'LV Falcons'), 3, 62, 61, 0, 0, 1.016, 1.016),
((SELECT id FROM players WHERE gamertag = 'Pred'), (SELECT id FROM teams WHERE name = 'LV Falcons'), 3, 48, 59, 0, 0, 0.814, 0.814),
((SELECT id FROM players WHERE gamertag = 'Arcitys'), (SELECT id FROM teams WHERE name = 'LV Falcons'), 3, 49, 57, 0, 0, 0.860, 0.860),
((SELECT id FROM players WHERE gamertag = 'Priestahh'), (SELECT id FROM teams WHERE name = 'LV Falcons'), 3, 38, 60, 0, 0, 0.633, 0.633),
-- Cloud9 NY
((SELECT id FROM players WHERE gamertag = 'Sib'), (SELECT id FROM teams WHERE name = 'Cloud9 NY'), 3, 46, 47, 0, 0, 0.979, 0.979),
((SELECT id FROM players WHERE gamertag = 'Attach'), (SELECT id FROM teams WHERE name = 'Cloud9 NY'), 3, 30, 42, 0, 0, 0.714, 0.714),
((SELECT id FROM players WHERE gamertag = 'Kremp'), (SELECT id FROM teams WHERE name = 'Cloud9 NY'), 3, 33, 45, 0, 0, 0.733, 0.733),
((SELECT id FROM players WHERE gamertag = 'Mack'), (SELECT id FROM teams WHERE name = 'Cloud9 NY'), 3, 42, 46, 0, 0, 0.913, 0.913),
-- LA Guerrillas M8
((SELECT id FROM players WHERE gamertag = 'KiSMET'), (SELECT id FROM teams WHERE name = 'LA Guerrillas M8'), 3, 40, 63, 0, 0, 0.635, 0.635),
((SELECT id FROM players WHERE gamertag = 'Lunarz'), (SELECT id FROM teams WHERE name = 'LA Guerrillas M8'), 3, 45, 50, 0, 0, 0.900, 0.900),
((SELECT id FROM players WHERE gamertag = 'FeLo'), (SELECT id FROM teams WHERE name = 'LA Guerrillas M8'), 3, 42, 50, 0, 0, 0.840, 0.840),
((SELECT id FROM players WHERE gamertag = 'oJohnny'), (SELECT id FROM teams WHERE name = 'LA Guerrillas M8'), 3, 42, 60, 0, 0, 0.700, 0.700)
ON CONFLICT DO NOTHING; 