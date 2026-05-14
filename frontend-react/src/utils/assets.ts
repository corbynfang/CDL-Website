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

// Map DB team names → logo file keys
const teamLogoKeys: Record<string, string> = {
  'vegas falcons':          'vegasfalconslogo',
  'carolina royal ravens':  'carolinalogo',
  'atlanta faze':           'fazenewlogo',
  'faze clan':              'fazenewlogo',
  'toronto ultra':          'torontonewlogo',
  'la guerrillas m8':       'laglogo',
  'los angeles guerrillas': 'laglogo',
  'lag academy':            'laglogo',
  'minnesota røkkr':        'minnesotalogo',
  'g2 minnesota':           'minnesotalogo',
  'optic texas':            'optictexaslogo',
  'optic gaming':           'optictexaslogo',
  'optic chicago':          'optictexaslogo',
  'dallas empire':          'optictexaslogo',
  'boston breach':          'bostonlogo',
  'miami heretics':         'miamilogo',
  'team heretics':          'miamilogo',
  'vancouver surge':        'vancouverlogo',
  'seattle surge':          'vancouverlogo',
  'los angeles thieves':    'latheiveslogo',
  '100 thieves':            'latheiveslogo',
  'team orchid':            'orchidlogo',
  'team falcons':           'newriyadhfalconslogo',
  'new riyadh falcons':     'newriyadhfalconslogo',
  'team war':               'teamwarlogo',
  'fivefears':              'fivefearslogo',
  'five fears':             'fivefearslogo',
  'lore':                   'lorelogo',
};

// Gamertags whose DB spelling differs from the avatar filename.
const avatarNicknames: Record<string, string> = {
  'hicksy':   'hicksey',
  'mercules': 'merc',
  'ojohnny':  'ojohnyy',
  'purj':     'purj_',
  'reeal':    'real',
  'lyynnz':   'lynz',   // same player, different spelling in some CSVs
  '04':       'new04',  // avatar file saved as New04
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
