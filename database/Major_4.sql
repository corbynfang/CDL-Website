-- Major 4 Data for CDL 2025
-- Assumes schema.sql has already been run

-- Insert Major 4 Teams
INSERT INTO teams (name, abbreviation, city, is_active) VALUES
('MIA Heretics', 'MIA', 'Miami', true),
('ATL FaZe', 'ATL', 'Atlanta', true),
('OpTic TEX', 'TEX', 'Dallas', true),
('LA Thieves', 'LAT', 'Los Angeles', true),
('VAN Surge', 'VAN', 'Vancouver', true),
('TOR Ultra', 'TOR', 'Toronto', true),
('Cloud9 NY', 'NYC', 'New York', true),
('MIN RØKKR', 'MIN', 'Minneapolis', true),
('BOS Breach', 'BOS', 'Boston', true),
('LA Guerrillas M8', 'LAG', 'Los Angeles', true),
('CAR Royal Ravens', 'CAR', 'Charlotte', true),
('LV Falcons', 'LVF', 'Las Vegas', true)
ON CONFLICT DO NOTHING;

-- Insert Major 4 Coaches
CREATE TABLE IF NOT EXISTS coaches (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    team_id INTEGER REFERENCES teams(id),
    season_id INTEGER REFERENCES seasons(id)
);

INSERT INTO coaches (name, team_id, season_id) VALUES
('MethodZ', (SELECT id FROM teams WHERE name = 'MIA Heretics'), 6),
('Sikotik', (SELECT id FROM teams WHERE name = 'MIA Heretics'), 6),
('Crowder', (SELECT id FROM teams WHERE name = 'ATL FaZe'), 6),
('RJ', (SELECT id FROM teams WHERE name = 'ATL FaZe'), 6),
('Karma', (SELECT id FROM teams WHERE name = 'OpTic TEX'), 6),
('JP Krez', (SELECT id FROM teams WHERE name = 'OpTic TEX'), 6),
('ShAnE', (SELECT id FROM teams WHERE name = 'LA Thieves'), 6),
('Sender', (SELECT id FROM teams WHERE name = 'LA Thieves'), 6),
('Rambo', (SELECT id FROM teams WHERE name = 'VAN Surge'), 6),
('Flux', (SELECT id FROM teams WHERE name = 'TOR Ultra'), 6),
('Joee', (SELECT id FROM teams WHERE name = 'TOR Ultra'), 6),
('Accuracy', (SELECT id FROM teams WHERE name = 'Cloud9 NY'), 6),
('Arian', (SELECT id FROM teams WHERE name = 'Cloud9 NY'), 6),
('Loony', (SELECT id FROM teams WHERE name = 'MIN RØKKR'), 6),
('Alexdotzip', (SELECT id FROM teams WHERE name = 'MIN RØKKR'), 6),
('Seany', (SELECT id FROM teams WHERE name = 'BOS Breach'), 6),
('Magxck', (SELECT id FROM teams WHERE name = 'BOS Breach'), 6),
('DREAL', (SELECT id FROM teams WHERE name = 'LA Guerrillas M8'), 6),
('MarkyB', (SELECT id FROM teams WHERE name = 'LA Guerrillas M8'), 6),
('Saintt', (SELECT id FROM teams WHERE name = 'CAR Royal Ravens'), 6),
('Clayster', (SELECT id FROM teams WHERE name = 'LV Falcons'), 6);

-- Insert Major 4 Players (excluding coaches)
INSERT INTO players (gamertag, is_active) VALUES
-- MIA Heretics
('MettalZ', true),
('ReeaL', true),
('RenKoR', true),
('SupeR', true),
-- ATL FaZe
('aBeZy', true),
('Cellium', true),
('Simp', true),
('Drazah', true),
-- OpTic TEX
('Shotzzy', true),
('Dashy', true),
('Huke', true),
('Mercules', true),
-- LA Thieves
('Ghosty', true),
('Scrap', true),
('HyDra', true),
('Envoy', true),
-- VAN Surge
('Abuzah', true),
('04', true),
('Nastie', true),
('Neptune', true),
-- TOR Ultra
('CleanX', true),
('Insight', true),
('Beans', true),
('JoeDeceives', true),
-- Cloud9 NY
('Sib', true),
('Kremp', true),
('Capsidal', true),
('Spart', true),
-- MIN RØKKR
('Nero', true),
('Gio', true),
('Estreal', true),
('Kenny', true),
-- BOS Breach
('Snoopy', true),
('Cammy', true),
('Owakening', true),
('Purj', true),
-- LA Guerrillas M8
('KiSMET', true),
('Lunarz', true),
('oJohnny', true),
('FeLo', true),
-- CAR Royal Ravens
('Gwinn', true),
('TJHaLy', true),
('SlasheR', true),
('Wrecks', true),
-- LV Falcons
('Exnid', true),
('Pred', true),
('Arcitys', true),
('Priestahh', true)
ON CONFLICT DO NOTHING;

-- Assign players to teams for 2025 season (season_id = 6)
INSERT INTO team_rosters (team_id, player_id, season_id, is_starter, start_date) VALUES
-- MIA Heretics
((SELECT id FROM teams WHERE name = 'MIA Heretics'), (SELECT id FROM players WHERE gamertag = 'MettalZ'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'MIA Heretics'), (SELECT id FROM players WHERE gamertag = 'ReeaL'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'MIA Heretics'), (SELECT id FROM players WHERE gamertag = 'RenKoR'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'MIA Heretics'), (SELECT id FROM players WHERE gamertag = 'SupeR'), 6, true, '2024-11-01'),
-- ATL FaZe
((SELECT id FROM teams WHERE name = 'ATL FaZe'), (SELECT id FROM players WHERE gamertag = 'aBeZy'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'ATL FaZe'), (SELECT id FROM players WHERE gamertag = 'Cellium'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'ATL FaZe'), (SELECT id FROM players WHERE gamertag = 'Simp'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'ATL FaZe'), (SELECT id FROM players WHERE gamertag = 'Drazah'), 6, true, '2024-11-01'),
-- OpTic TEX
((SELECT id FROM teams WHERE name = 'OpTic TEX'), (SELECT id FROM players WHERE gamertag = 'Shotzzy'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'OpTic TEX'), (SELECT id FROM players WHERE gamertag = 'Dashy'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'OpTic TEX'), (SELECT id FROM players WHERE gamertag = 'Huke'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'OpTic TEX'), (SELECT id FROM players WHERE gamertag = 'Mercules'), 6, true, '2024-11-01'),
-- LA Thieves
((SELECT id FROM teams WHERE name = 'LA Thieves'), (SELECT id FROM players WHERE gamertag = 'Ghosty'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'LA Thieves'), (SELECT id FROM players WHERE gamertag = 'Scrap'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'LA Thieves'), (SELECT id FROM players WHERE gamertag = 'HyDra'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'LA Thieves'), (SELECT id FROM players WHERE gamertag = 'Envoy'), 6, true, '2024-11-01'),
-- VAN Surge
((SELECT id FROM teams WHERE name = 'VAN Surge'), (SELECT id FROM players WHERE gamertag = 'Abuzah'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'VAN Surge'), (SELECT id FROM players WHERE gamertag = '04'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'VAN Surge'), (SELECT id FROM players WHERE gamertag = 'Nastie'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'VAN Surge'), (SELECT id FROM players WHERE gamertag = 'Neptune'), 6, true, '2024-11-01'),
-- TOR Ultra
((SELECT id FROM teams WHERE name = 'TOR Ultra'), (SELECT id FROM players WHERE gamertag = 'CleanX'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'TOR Ultra'), (SELECT id FROM players WHERE gamertag = 'Insight'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'TOR Ultra'), (SELECT id FROM players WHERE gamertag = 'Beans'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'TOR Ultra'), (SELECT id FROM players WHERE gamertag = 'JoeDeceives'), 6, true, '2024-11-01'),
-- Cloud9 NY
((SELECT id FROM teams WHERE name = 'Cloud9 NY'), (SELECT id FROM players WHERE gamertag = 'Sib'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'Cloud9 NY'), (SELECT id FROM players WHERE gamertag = 'Kremp'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'Cloud9 NY'), (SELECT id FROM players WHERE gamertag = 'Capsidal'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'Cloud9 NY'), (SELECT id FROM players WHERE gamertag = 'Spart'), 6, true, '2024-11-01'),
-- MIN RØKKR
((SELECT id FROM teams WHERE name = 'MIN RØKKR'), (SELECT id FROM players WHERE gamertag = 'Nero'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'MIN RØKKR'), (SELECT id FROM players WHERE gamertag = 'Gio'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'MIN RØKKR'), (SELECT id FROM players WHERE gamertag = 'Estreal'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'MIN RØKKR'), (SELECT id FROM players WHERE gamertag = 'Kenny'), 6, true, '2024-11-01'),
-- BOS Breach
((SELECT id FROM teams WHERE name = 'BOS Breach'), (SELECT id FROM players WHERE gamertag = 'Snoopy'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'BOS Breach'), (SELECT id FROM players WHERE gamertag = 'Cammy'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'BOS Breach'), (SELECT id FROM players WHERE gamertag = 'Owakening'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'BOS Breach'), (SELECT id FROM players WHERE gamertag = 'Purj'), 6, true, '2024-11-01'),
-- LA Guerrillas M8
((SELECT id FROM teams WHERE name = 'LA Guerrillas M8'), (SELECT id FROM players WHERE gamertag = 'KiSMET'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'LA Guerrillas M8'), (SELECT id FROM players WHERE gamertag = 'Lunarz'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'LA Guerrillas M8'), (SELECT id FROM players WHERE gamertag = 'oJohnny'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'LA Guerrillas M8'), (SELECT id FROM players WHERE gamertag = 'FeLo'), 6, true, '2024-11-01'),
-- CAR Royal Ravens
((SELECT id FROM teams WHERE name = 'CAR Royal Ravens'), (SELECT id FROM players WHERE gamertag = 'Gwinn'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'CAR Royal Ravens'), (SELECT id FROM players WHERE gamertag = 'TJHaLy'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'CAR Royal Ravens'), (SELECT id FROM players WHERE gamertag = 'SlasheR'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'CAR Royal Ravens'), (SELECT id FROM players WHERE gamertag = 'Wrecks'), 6, true, '2024-11-01'),
-- LV Falcons
((SELECT id FROM teams WHERE name = 'LV Falcons'), (SELECT id FROM players WHERE gamertag = 'Exnid'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'LV Falcons'), (SELECT id FROM players WHERE gamertag = 'Pred'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'LV Falcons'), (SELECT id FROM players WHERE gamertag = 'Arcitys'), 6, true, '2024-11-01'),
((SELECT id FROM teams WHERE name = 'LV Falcons'), (SELECT id FROM players WHERE gamertag = 'Priestahh'), 6, true, '2024-11-01');

-- Insert cumulative player stats for Major 4 (players only, using OVERALL K, D, KD from your table)
INSERT INTO player_tournament_stats (player_id, team_id, tournament_id, total_kills, total_deaths, total_assists, total_damage, kd_ratio, kda_ratio) VALUES
-- MIA Heretics
((SELECT id FROM players WHERE gamertag = 'MettalZ'), (SELECT id FROM teams WHERE name = 'MIA Heretics'), 4, 257, 295, 0, 0, 0.871, 0.871),
((SELECT id FROM players WHERE gamertag = 'ReeaL'), (SELECT id FROM teams WHERE name = 'MIA Heretics'), 4, 296, 279, 0, 0, 1.061, 1.061),
((SELECT id FROM players WHERE gamertag = 'RenKoR'), (SELECT id FROM teams WHERE name = 'MIA Heretics'), 4, 310, 279, 0, 0, 1.111, 1.111),
((SELECT id FROM players WHERE gamertag = 'SupeR'), (SELECT id FROM teams WHERE name = 'MIA Heretics'), 4, 243, 254, 0, 0, 0.957, 0.957),
-- ATL FaZe
((SELECT id FROM players WHERE gamertag = 'aBeZy'), (SELECT id FROM teams WHERE name = 'ATL FaZe'), 4, 349, 350, 0, 0, 0.997, 0.997),
((SELECT id FROM players WHERE gamertag = 'Cellium'), (SELECT id FROM teams WHERE name = 'ATL FaZe'), 4, 345, 310, 0, 0, 1.113, 1.113),
((SELECT id FROM players WHERE gamertag = 'Simp'), (SELECT id FROM teams WHERE name = 'ATL FaZe'), 4, 389, 322, 0, 0, 1.208, 1.208),
((SELECT id FROM players WHERE gamertag = 'Drazah'), (SELECT id FROM teams WHERE name = 'ATL FaZe'), 4, 343, 347, 0, 0, 0.989, 0.989),
-- OpTic TEX
((SELECT id FROM players WHERE gamertag = 'Shotzzy'), (SELECT id FROM teams WHERE name = 'OpTic TEX'), 4, 320, 325, 0, 0, 0.985, 0.985),
((SELECT id FROM players WHERE gamertag = 'Dashy'), (SELECT id FROM teams WHERE name = 'OpTic TEX'), 4, 288, 304, 0, 0, 0.947, 0.947),
((SELECT id FROM players WHERE gamertag = 'Huke'), (SELECT id FROM teams WHERE name = 'OpTic TEX'), 4, 311, 308, 0, 0, 1.010, 1.010),
((SELECT id FROM players WHERE gamertag = 'Mercules'), (SELECT id FROM teams WHERE name = 'OpTic TEX'), 4, 322, 317, 0, 0, 1.016, 1.016),
-- LA Thieves
((SELECT id FROM players WHERE gamertag = 'Ghosty'), (SELECT id FROM teams WHERE name = 'LA Thieves'), 4, 434, 435, 0, 0, 0.998, 0.998),
((SELECT id FROM players WHERE gamertag = 'Scrap'), (SELECT id FROM teams WHERE name = 'LA Thieves'), 4, 486, 445, 0, 0, 1.092, 1.092),
((SELECT id FROM players WHERE gamertag = 'HyDra'), (SELECT id FROM teams WHERE name = 'LA Thieves'), 4, 506, 457, 0, 0, 1.107, 1.107),
((SELECT id FROM players WHERE gamertag = 'Envoy'), (SELECT id FROM teams WHERE name = 'LA Thieves'), 4, 472, 471, 0, 0, 1.002, 1.002),
-- VAN Surge
((SELECT id FROM players WHERE gamertag = 'Abuzah'), (SELECT id FROM teams WHERE name = 'VAN Surge'), 4, 115, 135, 0, 0, 0.852, 0.852),
((SELECT id FROM players WHERE gamertag = '04'), (SELECT id FROM teams WHERE name = 'VAN Surge'), 4, 103, 144, 0, 0, 0.715, 0.715),
((SELECT id FROM players WHERE gamertag = 'Nastie'), (SELECT id FROM teams WHERE name = 'VAN Surge'), 4, 112, 130, 0, 0, 0.862, 0.862),
((SELECT id FROM players WHERE gamertag = 'Neptune'), (SELECT id FROM teams WHERE name = 'VAN Surge'), 4, 117, 139, 0, 0, 0.842, 0.842),
-- TOR Ultra
((SELECT id FROM players WHERE gamertag = 'CleanX'), (SELECT id FROM teams WHERE name = 'TOR Ultra'), 4, 363, 367, 0, 0, 0.989, 0.989),
((SELECT id FROM players WHERE gamertag = 'Insight'), (SELECT id FROM teams WHERE name = 'TOR Ultra'), 4, 292, 331, 0, 0, 0.882, 0.882),
((SELECT id FROM players WHERE gamertag = 'Beans'), (SELECT id FROM teams WHERE name = 'TOR Ultra'), 4, 325, 338, 0, 0, 0.962, 0.962),
((SELECT id FROM players WHERE gamertag = 'JoeDeceives'), (SELECT id FROM teams WHERE name = 'TOR Ultra'), 4, 371, 325, 0, 0, 1.142, 1.142),
-- Cloud9 NY
((SELECT id FROM players WHERE gamertag = 'Sib'), (SELECT id FROM teams WHERE name = 'Cloud9 NY'), 4, 124, 158, 0, 0, 0.785, 0.785),
((SELECT id FROM players WHERE gamertag = 'Kremp'), (SELECT id FROM teams WHERE name = 'Cloud9 NY'), 4, 132, 149, 0, 0, 0.886, 0.886),
((SELECT id FROM players WHERE gamertag = 'Capsidal'), (SELECT id FROM teams WHERE name = 'Cloud9 NY'), 4, 126, 159, 0, 0, 0.793, 0.793),
((SELECT id FROM players WHERE gamertag = 'Spart'), (SELECT id FROM teams WHERE name = 'Cloud9 NY'), 4, 128, 141, 0, 0, 0.908, 0.908),
-- MIN RØKKR
((SELECT id FROM players WHERE gamertag = 'Nero'), (SELECT id FROM teams WHERE name = 'MIN RØKKR'), 4, 370, 326, 0, 0, 1.135, 1.135),
((SELECT id FROM players WHERE gamertag = 'Gio'), (SELECT id FROM teams WHERE name = 'MIN RØKKR'), 4, 296, 288, 0, 0, 1.028, 1.028),
((SELECT id FROM players WHERE gamertag = 'Estreal'), (SELECT id FROM teams WHERE name = 'MIN RØKKR'), 4, 283, 303, 0, 0, 0.934, 0.934),
((SELECT id FROM players WHERE gamertag = 'Kenny'), (SELECT id FROM teams WHERE name = 'MIN RØKKR'), 4, 315, 309, 0, 0, 1.019, 1.019),
-- BOS Breach
((SELECT id FROM players WHERE gamertag = 'Snoopy'), (SELECT id FROM teams WHERE name = 'BOS Breach'), 4, 96, 94, 0, 0, 1.021, 1.021),
((SELECT id FROM players WHERE gamertag = 'Cammy'), (SELECT id FROM teams WHERE name = 'BOS Breach'), 4, 97, 104, 0, 0, 0.933, 0.933),
((SELECT id FROM players WHERE gamertag = 'Owakening'), (SELECT id FROM teams WHERE name = 'BOS Breach'), 4, 92, 100, 0, 0, 0.920, 0.920),
((SELECT id FROM players WHERE gamertag = 'Purj'), (SELECT id FROM teams WHERE name = 'BOS Breach'), 4, 95, 109, 0, 0, 0.872, 0.872),
-- LA Guerrillas M8
((SELECT id FROM players WHERE gamertag = 'KiSMET'), (SELECT id FROM teams WHERE name = 'LA Guerrillas M8'), 4, 147, 158, 0, 0, 0.930, 0.930),
((SELECT id FROM players WHERE gamertag = 'Lunarz'), (SELECT id FROM teams WHERE name = 'LA Guerrillas M8'), 4, 140, 129, 0, 0, 1.085, 1.085),
((SELECT id FROM players WHERE gamertag = 'oJohnny'), (SELECT id FROM teams WHERE name = 'LA Guerrillas M8'), 4, 164, 145, 0, 0, 1.131, 1.131),
((SELECT id FROM players WHERE gamertag = 'FeLo'), (SELECT id FROM teams WHERE name = 'LA Guerrillas M8'), 4, 137, 138, 0, 0, 0.993, 0.993),
-- CAR Royal Ravens
((SELECT id FROM players WHERE gamertag = 'Gwinn'), (SELECT id FROM teams WHERE name = 'CAR Royal Ravens'), 4, 89, 87, 0, 0, 1.023, 1.023),
((SELECT id FROM players WHERE gamertag = 'TJHaLy'), (SELECT id FROM teams WHERE name = 'CAR Royal Ravens'), 4, 90, 103, 0, 0, 0.874, 0.874),
((SELECT id FROM players WHERE gamertag = 'SlasheR'), (SELECT id FROM teams WHERE name = 'CAR Royal Ravens'), 4, 94, 82, 0, 0, 1.146, 1.146),
((SELECT id FROM players WHERE gamertag = 'Wrecks'), (SELECT id FROM teams WHERE name = 'CAR Royal Ravens'), 4, 75, 83, 0, 0, 0.904, 0.904),
-- LV Falcons
((SELECT id FROM players WHERE gamertag = 'Exnid'), (SELECT id FROM teams WHERE name = 'LV Falcons'), 4, 186, 176, 0, 0, 1.057, 1.057),
((SELECT id FROM players WHERE gamertag = 'Pred'), (SELECT id FROM teams WHERE name = 'LV Falcons'), 4, 193, 157, 0, 0, 1.229, 1.229),
((SELECT id FROM players WHERE gamertag = 'Arcitys'), (SELECT id FROM teams WHERE name = 'LV Falcons'), 4, 163, 167, 0, 0, 0.976, 0.976),
((SELECT id FROM players WHERE gamertag = 'Priestahh'), (SELECT id FROM teams WHERE name = 'LV Falcons'), 4, 135, 175, 0, 0, 0.771, 0.771)
ON CONFLICT DO NOTHING; 