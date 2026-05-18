// Dynamically import every avatar and logo at build time so Vite bundles them.
const avatarModules = import.meta.glob(
  '../assets/avatars/*.{webp,png,jpg,jpeg}',
  { eager: true }
) as Record<string, { default: string }>;

const logoModules = import.meta.glob(
  '../assets/logos/*.{webp,png,jpg,jpeg}',
  { eager: true }
) as Record<string, { default: string }>;

// Build maps: lowercase filename (no extension) → resolved URL
const avatarMap: Record<string, string> = {};
for (const [path, mod] of Object.entries(avatarModules)) {
  const key = path.split('/').pop()!.replace(/\.[^.]+$/, '').toLowerCase();
  avatarMap[key] = mod.default;
}

const logoMap: Record<string, string> = {};
for (const [path, mod] of Object.entries(logoModules)) {
  const key = path.split('/').pop()!.replace(/\.[^.]+$/, '').toLowerCase();
  logoMap[key] = mod.default;
}

// Map DB team names → logo file keys (lowercase name → lowercase filename without extension)
export const teamLogoKeys: Record<string, string> = {
  // ── CDL Franchises ──────────────────────────────────────────────────────────

  // Atlanta / FaZe
  'atlanta faze':              'atlantalogo',
  'faze clan':                 'fazeclanold',
  'faze black':                'fazeblack',

  // Boston
  'boston breach':             'bostonbreach',

  // Carolina
  'carolina royal ravens':     'carolinaroyalravens',

  // Dallas / OpTic lineage
  'dallas empire':             'dallasempire',
  'optic texas':               'optictexaslogo',
  'optic gaming':              'opticgaming',
  'optic chicago':             'optic_chicagologo_std',

  // Florida
  'florida mutineers':         'florida_mutineerslogo_square',

  // LA Guerrillas
  'los angeles guerrillas':    'losanglesgurreillas',
  'la guerrillas m8':          'lagurreillasm8',
  'los angeles guerrillas m8': 'lagurreillasm8',

  // LA Thieves / 100T
  'los angeles thieves':       'latheiveslogo',
  '100 thieves':               'latheiveslogo',

  // Las Vegas / Falcons
  'las vegas legion':          'lasvegaslegion',
  'las vegas falcons':         'vegasfalconslogo',
  'vegas falcons':             'vegasfalconslogo',

  // London
  'london royal ravens':       'londonroyalravens',

  // Miami / Team Heretics
  'miami heretics':            'miamiheretics',
  'team heretics':             'teamheretics',

  // Minnesota
  'minnesota røkkr':           'minnesotarokker',
  'g2 minnesota':              'g2minnesota',

  // New York
  'new york subliners':        'newyorksubliners',
  'cloud9 new york':           'cloud9logo',
  'cloud9':                    'cloud9logo',

  // Paris
  'paris legion':              'parislegion',

  // Riyadh / Team Falcons (EWC name)
  'riyadh falcons':            'riyadhfalconslogo',
  'new riyadh falcons':        'newriyadhfalconslogo',
  'team falcons':              'newriyadhfalconslogo',

  // Toronto
  'toronto ultra':             'torontonewlogo',
  'toronto koi':               'torontokoi',

  // Vancouver / Seattle
  'vancouver surge':           'vancouversurge',
  'seattle surge':             'seattlesurgelogo',

  // ── Academy / Sub-teams ─────────────────────────────────────────────────────
  'lag academy':               'laglogo',
  'boston academy':            'bostonacademylogo',
  'falcons academy':           'falconsacademylogo',
  'røkkr academy':             'minnesotarokker',
  'toronto ultra academy na':  'toronto_ultraacademylogo',

  // ── CDL Challengers placeholder ─────────────────────────────────────────────
  '18andcracked':              'callofdutychallengersplaceholderforchallengerteams',
  'slammed':                   'callofdutychallengersplaceholderforchallengerteams',

  // ── Challenger Orgs (real logos) ─────────────────────────────────────────────
  'fantastic four':            'fantasticfourlogo',
  'fc black':                  'fc_blacklogo_square',
  'infamous esports':          'infamous_esportslogo_square',
  'millennium 7':              'millennium_7logo_square',
  'ut crew':                   'utcrewlogo',
  'we are trying now':         'wearetryingnowlogo',

  // ── Challenger Orgs ─────────────────────────────────────────────────────────
  '705 esports':               '705esportslogo',
  'aphelion esports':          'aphelionesportslogo',
  'clutch rayn esport':        'clutchraynesportlogo',
  'convoy gaming':             'convoygaminglogo',
  'decimate gaming':           'decimategaminglogo',
  'deviance x notorious':      'deviancexnotoriouslogo',
  'eastr':                     'eastrlogo',
  'electrify steel gaming black': 'electrifysteelgamingblacklogo',
  'hive':                      'hivelogo',
  'houston spartans':          'houstonspartanslogo',
  'iron blood gaming':         'ironbloodgaminglogo',
  'katana gaming':             'katanagaminglogo',
  'kc pioneers':               'kcpioneerslogo',
  'omit':                      'omit',
  'omit brooklyn':             'omit',
  'omit eu':                   'omiteulogo',
  'singapore syndicate':       'singaporesyndicatelogo',
  'stallions':                 'stallionslogo',
  'texas nation':              'texasnationlogo',
  'the vicious':               'theviciouslogo',
  'twisted theory':            'twistedtheorylogo',
  'westr':                     'westrlogo',
  'whateverittakes':           'whateverittakeslogo',

  // ── Shared / Variant Logos ──────────────────────────────────────────────────
  'fivefears':                 'fivefearslogo',
  'five fears':                'fivefearslogo',
  'fivefears white':           'fivefearslogo',
  'gentle mates':              'gentlemates',
  'gentlemates':               'gentlemates',
  'koi':                       'koi',
  'lore':                      'lorelogo',
  'lore gaming':               'lorelogo',
  'lore black':                'lorelogo',
  'lore gold':                 'lorelogo',
  'project 7':                 'project7',
  'team orchid':               'orchidlogo',
  'team valcons':              'teamvalcons',
  'team war':                  'teamwarlogo',
};

// Gamertags whose DB spelling differs from the avatar filename.
export const avatarNicknames: Record<string, string> = {
  'hicksy':   'hicksey',
  'mercules': 'merc',
  'ojohnny':  'ojohnyy',
  'purj':     'purj_',
  'reeal':    'real',
  'lyynnz':   'lynz',
  // Note: '04' avatar — 04.webp (old) and New04.webp (current) both exist.
  // Delete src/assets/avatars/04.webp to make this nickname active.
  '04':       'new04',

  // No public photo available — explicitly use the placeholder.
  '5aldx':    'unknown',
  'felony':   'unknown',
  'hamza':    'unknown',
  'markyb':   'unknown',
  'qk4b':     'unknown',
};

/**
 * Returns the local avatar image URL for a player's gamertag.
 * Falls back to the Unknown avatar if no match is found.
 */
export function getPlayerAvatar(gamertag: string): string {
  const key = gamertag.toLowerCase();
  return (
    avatarMap[key] ??
    avatarMap[avatarNicknames[key] ?? ''] ??
    avatarMap['unknown'] ??
    ''
  );
}

/**
 * Returns the local logo image URL for a team name.
 * Returns empty string if no logo is available.
 */
export function getTeamLogo(teamName: string): string {
  const key = teamLogoKeys[teamName.toLowerCase()];
  return key ? (logoMap[key] ?? '') : '';
}
