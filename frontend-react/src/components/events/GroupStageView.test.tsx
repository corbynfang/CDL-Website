import { describe, it, expect, vi } from 'vitest'
import { render, screen } from '@testing-library/react'
import { MemoryRouter } from 'react-router-dom'
import GroupStageView from './GroupStageView'
import { cdlGroupBracketData, ewcBracketData } from '../../test/fixtures/matches'

vi.mock('../../utils/assets', () => ({
  getTeamLogo: vi.fn().mockReturnValue(null),
  getPlayerAvatar: vi.fn().mockReturnValue('/placeholder.png'),
}))

function wrap(ui: React.ReactElement) {
  return render(<MemoryRouter>{ui}</MemoryRouter>)
}

describe('GroupStageView', () => {
  it('CDL format: renders flat round columns', () => {
    wrap(
      <GroupStageView
        groupStage={cdlGroupBracketData.group_stage!}
        format="cdl_major_group_stage_bracket"
      />
    )
    expect(screen.getByText('Round 1')).toBeInTheDocument()
    expect(screen.getByText('Qualification Match')).toBeInTheDocument()
    expect(screen.getByText('Losers Bracket')).toBeInTheDocument()
  })

  it('CDL format: does not render group cards', () => {
    wrap(
      <GroupStageView
        groupStage={cdlGroupBracketData.group_stage!}
        format="cdl_major_group_stage_bracket"
      />
    )
    expect(screen.queryByText(/group a/i)).not.toBeInTheDocument()
  })

  it('EWC format: renders cross-group rounds section', () => {
    wrap(
      <GroupStageView
        groupStage={ewcBracketData.group_stage!}
        format="ewc_group_stage_single_elim"
      />
    )
    expect(screen.getByText('Opening Match')).toBeInTheDocument()
    expect(screen.getByText('Winners Match')).toBeInTheDocument()
  })

  it('EWC format: renders per-group sections for groups present in data', () => {
    wrap(
      <GroupStageView
        groupStage={ewcBracketData.group_stage!}
        format="ewc_group_stage_single_elim"
      />
    )
    expect(screen.getByText('Group A')).toBeInTheDocument()
    expect(screen.getByText('Group B')).toBeInTheDocument()
    // Groups C and D have no keys in fixture
    expect(screen.queryByText('Group C')).not.toBeInTheDocument()
    expect(screen.queryByText('Group D')).not.toBeInTheDocument()
  })

  it('EWC format: strips group prefix from round labels', () => {
    wrap(
      <GroupStageView
        groupStage={ewcBracketData.group_stage!}
        format="ewc_group_stage_single_elim"
      />
    )
    // group_play_a_winners_round_1 → "Winners Round 1" under Group A
    expect(screen.getByText('Winners Round 1')).toBeInTheDocument()
  })

  it('renders empty state when group stage has no rounds', () => {
    wrap(<GroupStageView groupStage={{}} format="cdl_major_group_stage_bracket" />)
    expect(screen.getByText(/no group stage data/i)).toBeInTheDocument()
  })
})
