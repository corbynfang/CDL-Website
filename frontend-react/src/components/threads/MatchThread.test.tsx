import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import type { Session, User } from '@supabase/supabase-js'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import type { ThreadPost, ThreadResponse } from '../../types'
import MatchThread from './MatchThread'

vi.mock('../../services/threadApi', () => ({
  threadApi: {
    getThread: vi.fn(),
    createPost: vi.fn(),
    editPost: vi.fn(),
    deletePost: vi.fn(),
  },
}))

vi.mock('../../context/AuthContext', () => ({
  useAuth: vi.fn(),
}))

import { threadApi } from '../../services/threadApi'
import { useAuth } from '../../context/AuthContext'

const mockThreadApi = vi.mocked(threadApi)
const mockUseAuth = vi.mocked(useAuth)

const emptyThread: ThreadResponse = {
  thread_id: 1,
  data: [],
  pagination: { page: 1, limit: 25, total: 0, total_pages: 1 },
}

const threadWithPosts: ThreadResponse = {
  thread_id: 1,
  data: [
    {
      id: 1,
      thread_id: 1,
      user_id: 10,
      body: 'OpTic looking great!',
      edited: false,
      created_at: '2026-06-07T12:00:00Z',
      updated_at: '2026-06-07T12:00:00Z',
      user: { id: 10, supabase_uid: 'uid-owner', username: 'Corbyn', created_at: '', updated_at: '' },
    },
    {
      id: 2,
      thread_id: 1,
      user_id: 11,
      body: 'G2 strong too',
      edited: true,
      created_at: '2026-06-07T13:00:00Z',
      updated_at: '2026-06-07T14:00:00Z',
      user: { id: 11, supabase_uid: 'uid-other', username: 'OtherUser', created_at: '', updated_at: '' },
    },
  ],
  pagination: { page: 1, limit: 25, total: 2, total_pages: 1 },
}

function noAuth() {
  mockUseAuth.mockReturnValue({
    user: null, session: null, loading: false,
    showAuthModal: false, needsProfileSetup: false,
    openAuthModal: vi.fn(), closeAuthModal: vi.fn(),
    signUp: vi.fn(), signIn: vi.fn(), signInWithOAuth: vi.fn(),
    completeProfileSetup: vi.fn(), signOut: vi.fn(),
  })
}

function withAuth(uid = 'uid-owner') {
  mockUseAuth.mockReturnValue({
    user: { id: uid } as unknown as User,
    session: { access_token: 'tok' } as unknown as Session,
    loading: false,
    showAuthModal: false, needsProfileSetup: false,
    openAuthModal: vi.fn(), closeAuthModal: vi.fn(),
    signUp: vi.fn(), signIn: vi.fn(), signInWithOAuth: vi.fn(),
    completeProfileSetup: vi.fn(), signOut: vi.fn(),
  })
}

beforeEach(() => {
  vi.clearAllMocks()
  mockThreadApi.getThread.mockResolvedValue(emptyThread)
  noAuth()
})

describe('MatchThread — unauthenticated', () => {
  it('shows sign-in prompt when not logged in', async () => {
    render(<MatchThread matchId={1207} />)
    await waitFor(() => {
      expect(screen.getByText(/Sign in to join the discussion/i)).toBeInTheDocument()
    })
  })

  it('does not show the compose textarea', async () => {
    render(<MatchThread matchId={1207} />)
    await waitFor(() => screen.getByText(/Sign in to join the discussion/i))
    expect(screen.queryByPlaceholderText(/Write a comment/i)).not.toBeInTheDocument()
  })
})

describe('MatchThread — empty state', () => {
  it('shows empty state when there are no posts', async () => {
    render(<MatchThread matchId={1207} />)
    await waitFor(() => {
      expect(screen.getByText(/No comments yet/i)).toBeInTheDocument()
    })
  })

  it('shows 0 in the section header', async () => {
    render(<MatchThread matchId={1207} />)
    await waitFor(() => {
      expect(screen.getByText(/Discussion · 0 comments/i)).toBeInTheDocument()
    })
  })
})

describe('MatchThread — authenticated compose', () => {
  it('shows compose textarea when logged in', async () => {
    withAuth()
    render(<MatchThread matchId={1207} />)
    await waitFor(() => {
      expect(screen.getByPlaceholderText(/Write a comment/i)).toBeInTheDocument()
    })
  })

  it('submits a post and re-fetches the thread', async () => {
    withAuth()
    mockThreadApi.createPost.mockResolvedValue({} as unknown as ThreadPost)
    const user = userEvent.setup()
    render(<MatchThread matchId={1207} />)

    const textarea = await screen.findByPlaceholderText(/Write a comment/i)
    await user.type(textarea, 'New post!')
    await user.click(screen.getByRole('button', { name: /^Post$/i }))

    expect(mockThreadApi.createPost).toHaveBeenCalledWith(1207, 'New post!')
    expect(mockThreadApi.getThread).toHaveBeenCalledTimes(2)
  })

  it('shows an error message when posting fails', async () => {
    withAuth()
    mockThreadApi.createPost.mockRejectedValue({ isAxiosError: true, response: { data: { error: 'Rate limited' } } })
    const user = userEvent.setup()
    render(<MatchThread matchId={1207} />)

    const textarea = await screen.findByPlaceholderText(/Write a comment/i)
    await user.type(textarea, 'Some text')
    await user.click(screen.getByRole('button', { name: /^Post$/i }))

    await waitFor(() => {
      expect(screen.getByText('Rate limited')).toBeInTheDocument()
    })
  })
})

describe('MatchThread — posts rendering', () => {
  beforeEach(() => {
    mockThreadApi.getThread.mockResolvedValue(threadWithPosts)
  })

  it('renders post bodies', async () => {
    noAuth()
    render(<MatchThread matchId={1207} />)
    await waitFor(() => {
      expect(screen.getByText('OpTic looking great!')).toBeInTheDocument()
      expect(screen.getByText('G2 strong too')).toBeInTheDocument()
    })
  })

  it('renders usernames', async () => {
    noAuth()
    render(<MatchThread matchId={1207} />)
    await waitFor(() => {
      expect(screen.getByText('Corbyn')).toBeInTheDocument()
      expect(screen.getByText('OtherUser')).toBeInTheDocument()
    })
  })

  it('shows edited label for edited posts', async () => {
    noAuth()
    render(<MatchThread matchId={1207} />)
    await waitFor(() => {
      expect(screen.getByText('edited')).toBeInTheDocument()
    })
  })

  it('shows the correct comment count in the header', async () => {
    noAuth()
    render(<MatchThread matchId={1207} />)
    await waitFor(() => {
      expect(screen.getByText(/Discussion · 2 comments/i)).toBeInTheDocument()
    })
  })
})

describe('MatchThread — ownership controls', () => {
  beforeEach(() => {
    mockThreadApi.getThread.mockResolvedValue(threadWithPosts)
  })

  it('shows edit and delete only for the logged-in user own post', async () => {
    withAuth('uid-owner')
    render(<MatchThread matchId={1207} />)
    await waitFor(() => screen.getByText('OpTic looking great!'))

    expect(screen.getAllByText('Edit')).toHaveLength(1)
    expect(screen.getAllByText('Delete')).toHaveLength(1)
  })

  it('shows no edit or delete when logged-in user owns none of the posts', async () => {
    withAuth('uid-nobody')
    render(<MatchThread matchId={1207} />)
    await waitFor(() => screen.getByText('OpTic looking great!'))

    expect(screen.queryByText('Edit')).not.toBeInTheDocument()
    expect(screen.queryByText('Delete')).not.toBeInTheDocument()
  })
})

describe('MatchThread — pagination', () => {
  it('hides pagination when there is only one page', async () => {
    noAuth()
    render(<MatchThread matchId={1207} />)
    await waitFor(() => screen.getByText(/Discussion · 0 comments/i))
    expect(screen.queryByText('Prev')).not.toBeInTheDocument()
    expect(screen.queryByText('Next')).not.toBeInTheDocument()
  })

  it('shows pagination when there are multiple pages', async () => {
    noAuth()
    mockThreadApi.getThread.mockResolvedValue({
      ...emptyThread,
      pagination: { page: 1, limit: 25, total: 50, total_pages: 2 },
    })
    render(<MatchThread matchId={1207} />)
    await waitFor(() => {
      expect(screen.getByText('Prev')).toBeInTheDocument()
      expect(screen.getByText('Next')).toBeInTheDocument()
    })
  })
})
