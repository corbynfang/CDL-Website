import { describe, it, expect } from 'vitest'
import { getTeamLogo, getPlayerAvatar } from './assets'

describe('getTeamLogo', () => {
  it('returns a non-empty string for a known team name', () => {
    expect(getTeamLogo('Atlanta Faze')).not.toBe('')
  })

  it('is case-insensitive', () => {
    expect(getTeamLogo('ATLANTA FAZE')).toBe(getTeamLogo('atlanta faze'))
  })

  it('returns empty string for an unknown team', () => {
    expect(getTeamLogo('Some Random Team')).toBe('')
  })

  it('resolves Boston Breach', () => {
    expect(getTeamLogo('Boston Breach')).not.toBe('')
  })

  it('resolves OpTic Texas', () => {
    expect(getTeamLogo('OpTic Texas')).not.toBe('')
  })

  it('resolves LA Guerrillas M8 alias', () => {
    expect(getTeamLogo('LA Guerrillas M8')).not.toBe('')
  })
})

describe('getPlayerAvatar', () => {
  it('returns a non-empty string for a known gamertag', () => {
    expect(getPlayerAvatar('AbeZy')).not.toBe('')
  })

  it('is case-insensitive', () => {
    expect(getPlayerAvatar('abezy')).toBe(getPlayerAvatar('ABEZY'))
  })

  it('resolves nickname aliases (Hicksy → hicksey file)', () => {
    expect(getPlayerAvatar('Hicksy')).not.toBe('')
  })

  it('returns empty string or unknown avatar for completely unknown gamertag', () => {
    const result = getPlayerAvatar('zzz_nobody_zzz')
    expect(typeof result).toBe('string')
  })
})
