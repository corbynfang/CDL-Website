const avatarModules = import.meta.glob(
  '../assets/avatars/*.{webp,png,jpg,jpeg}',
  { eager: true }
) as Record<string, { default: string }>;

const avatarMap: Record<string, string> = {};
for (const [path, mod] of Object.entries(avatarModules)) {
  const key = path.split('/').pop()!.replace(/\.[^.]+$/, '').toLowerCase();
  avatarMap[key] = mod.default;
}

export const avatarNicknames: Record<string, string> = {
  'hicksy':   'hicksey',
  'mercules': 'merc',
  'ojohnny':  'ojohnyy',
  'purj':     'purj_',
  'reeal':    'real',
  'lyynnz':   'lynz',
  '04':       'new04',
  '5aldx':    'unknown',
  'felony':   'unknown',
  'hamza':    'unknown',
  'markyb':   'unknown',
  'qk4b':     'unknown',
};

export function getPlayerAvatar(gamertag: string): string {
  const key = gamertag.toLowerCase();
  return (
    avatarMap[key] ??
    avatarMap[avatarNicknames[key] ?? ''] ??
    avatarMap['unknown'] ??
    ''
  );
}
