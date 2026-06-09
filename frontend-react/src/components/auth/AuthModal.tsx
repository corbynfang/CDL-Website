import { useState } from 'react';
import { useAuth } from '../../context/AuthContext';

const OAUTH_PROVIDERS = [
  {
    id: 'github' as const,
    label: 'Continue with GitHub',
    icon: (
      <svg viewBox="0 0 24 24" fill="currentColor" className="w-4 h-4">
        <path d="M12 0C5.37 0 0 5.37 0 12c0 5.31 3.435 9.795 8.205 11.385.6.105.825-.255.825-.57 0-.285-.015-1.23-.015-2.235-3.015.555-3.795-.735-4.035-1.41-.135-.345-.72-1.41-1.23-1.695-.42-.225-1.02-.78-.015-.795.945-.015 1.62.87 1.845 1.23 1.08 1.815 2.805 1.305 3.495.99.105-.78.42-1.305.765-1.605-2.67-.3-5.46-1.335-5.46-5.925 0-1.305.465-2.385 1.23-3.225-.12-.3-.54-1.53.12-3.18 0 0 1.005-.315 3.3 1.23.96-.27 1.98-.405 3-.405s2.04.135 3 .405c2.295-1.56 3.3-1.23 3.3-1.23.66 1.65.24 2.88.12 3.18.765.84 1.23 1.905 1.23 3.225 0 4.605-2.805 5.625-5.475 5.925.435.375.81 1.095.81 2.22 0 1.605-.015 2.895-.015 3.3 0 .315.225.69.825.57A12.02 12.02 0 0 0 24 12c0-6.63-5.37-12-12-12z" />
      </svg>
    ),
  },
  {
    id: 'twitch' as const,
    label: 'Continue with Twitch',
    icon: (
      <svg viewBox="0 0 24 24" fill="currentColor" className="w-4 h-4">
        <path d="M11.571 4.714h1.715v5.143H11.57zm4.715 0H18v5.143h-1.714zM6 0L1.714 4.286v15.428h5.143V24l4.286-4.286h3.428L22.286 12V0zm14.571 11.143l-3.428 3.428h-3.429l-3 3v-3H6.857V1.714h13.714z" />
      </svg>
    ),
  },
  {
    id: 'google' as const,
    label: 'Continue with Google',
    icon: (
      <svg viewBox="0 0 24 24" className="w-4 h-4">
        <path fill="#4285F4" d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z" />
        <path fill="#34A853" d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z" />
        <path fill="#FBBC05" d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z" />
        <path fill="#EA4335" d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z" />
      </svg>
    ),
  },
] as const;

type EmailMode = 'login' | 'signup';

interface Props {
  onClose: () => void;
}

export default function AuthModal({ onClose }: Props) {
  const { signIn, signUp, signInWithOAuth, completeProfileSetup, needsProfileSetup } = useAuth();

  const [showEmailForm, setShowEmailForm] = useState(false);
  const [emailMode, setEmailMode] = useState<EmailMode>('login');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [username, setUsername] = useState('');
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);
  const [emailSent, setEmailSent] = useState(false);

  async function handleOAuth(provider: 'github' | 'twitch' | 'google') {
    setError(null);
    const { error } = await signInWithOAuth(provider);
    if (error) setError(error);
  }

  async function handleEmailSubmit(e: React.FormEvent) {
    e.preventDefault();
    setError(null);
    setLoading(true);

    if (emailMode === 'login') {
      const { error } = await signIn(email, password);
      if (error) { setError(error); setLoading(false); return; }
      onClose();
    } else {
      const { error } = await signUp(email, password);
      if (error) { setError(error); setLoading(false); return; }
      // Profile setup happens after email confirmation + first sign-in, not here —
      // there is no session until the user clicks the confirmation link.
      setEmailSent(true);
    }
    setLoading(false);
  }

  async function handleUsernameSetup(e: React.FormEvent) {
    e.preventDefault();
    setError(null);
    setLoading(true);
    const { error } = await completeProfileSetup(username);
    if (error) setError(error);
    setLoading(false);
  }

  const canClose = !needsProfileSetup;

  function handleOverlayClick() {
    if (canClose) onClose();
  }

  if (emailSent) {
    return (
      <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/70 backdrop-blur-sm" onClick={onClose}>
        <div className="bg-[#111111] border border-[#1a1a1a] p-8 w-full max-w-sm mx-4" onClick={e => e.stopPropagation()}>
          <p className="text-xs uppercase tracking-widest text-[#737373] mb-2">Check your email</p>
          <p className="text-white text-sm">We sent a confirmation link to <span className="text-[#a3a3a3]">{email}</span>. Click it to activate your account.</p>
          <button onClick={onClose} className="mt-6 w-full py-2 text-xs uppercase tracking-widest text-[#737373] border border-[#1a1a1a] hover:border-[#404040] hover:text-white transition-colors">
            Close
          </button>
        </div>
      </div>
    );
  }

  if (needsProfileSetup) {
    return (
      <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/70 backdrop-blur-sm">
        <div className="bg-[#111111] border border-[#1a1a1a] p-8 w-full max-w-sm mx-4">
          <p className="text-xs uppercase tracking-widest text-[#737373] mb-1">One last step</p>
          <h2 className="font-grotesk font-bold text-white text-lg mb-6">Choose a username</h2>
          <form onSubmit={handleUsernameSetup}>
            <input
              type="text"
              placeholder="Username (3–30 characters)"
              value={username}
              onChange={e => setUsername(e.target.value)}
              required
              minLength={3}
              maxLength={30}
              autoFocus
              className="w-full bg-[#0a0a0a] border border-[#1a1a1a] text-white text-sm placeholder-[#404040] px-3 py-2.5 focus:outline-none focus:border-[#404040] transition-colors"
            />
            {error && <p className="text-red-400 text-xs mt-2">{error}</p>}
            <button
              type="submit"
              disabled={loading || username.trim().length < 3}
              className="mt-4 w-full py-2.5 text-xs font-grotesk font-semibold uppercase tracking-widest bg-white text-black disabled:opacity-30 disabled:cursor-not-allowed hover:bg-[#e0e0e0] transition-colors"
            >
              {loading ? 'Setting up…' : 'Continue'}
            </button>
          </form>
        </div>
      </div>
    );
  }

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/70 backdrop-blur-sm" onClick={handleOverlayClick}>
      <div className="bg-[#111111] border border-[#1a1a1a] p-8 w-full max-w-sm mx-4" onClick={e => e.stopPropagation()}>
        <div className="flex items-start justify-between mb-6">
          <div>
            <p className="text-xs uppercase tracking-widest text-[#737373] mb-1">CDLytics</p>
            <h2 className="font-grotesk font-bold text-white text-lg">
              {showEmailForm ? (emailMode === 'login' ? 'Sign In' : 'Create Account') : 'Sign In'}
            </h2>
          </div>
          <button onClick={onClose} className="text-[#404040] hover:text-white transition-colors text-lg leading-none mt-1">×</button>
        </div>

        {!showEmailForm ? (
          <div className="space-y-3">
            {OAUTH_PROVIDERS.map(p => (
              <button
                key={p.id}
                onClick={() => handleOAuth(p.id)}
                className="w-full flex items-center gap-3 px-4 py-2.5 border border-[#1a1a1a] text-[#a3a3a3] hover:border-[#404040] hover:text-white text-sm transition-colors"
              >
                {p.icon}
                <span>{p.label}</span>
              </button>
            ))}

            <div className="flex items-center gap-3 py-1">
              <div className="flex-1 h-px bg-[#1a1a1a]" />
              <span className="text-[#404040] text-[10px] uppercase tracking-widest">or</span>
              <div className="flex-1 h-px bg-[#1a1a1a]" />
            </div>

            <button
              onClick={() => setShowEmailForm(true)}
              className="w-full py-2.5 text-xs uppercase tracking-widest text-[#737373] border border-[#1a1a1a] hover:border-[#404040] hover:text-white transition-colors"
            >
              Continue with Email
            </button>

            {error && <p className="text-red-400 text-xs mt-1">{error}</p>}
          </div>
        ) : (
          <form onSubmit={handleEmailSubmit} className="space-y-3">
            <input
              type="email"
              placeholder="Email"
              value={email}
              onChange={e => setEmail(e.target.value)}
              required
              autoComplete="email"
              className="w-full bg-[#0a0a0a] border border-[#1a1a1a] text-white text-sm placeholder-[#404040] px-3 py-2.5 focus:outline-none focus:border-[#404040] transition-colors"
            />
            <input
              type="password"
              placeholder="Password"
              value={password}
              onChange={e => setPassword(e.target.value)}
              required
              autoComplete={emailMode === 'login' ? 'current-password' : 'new-password'}
              minLength={8}
              className="w-full bg-[#0a0a0a] border border-[#1a1a1a] text-white text-sm placeholder-[#404040] px-3 py-2.5 focus:outline-none focus:border-[#404040] transition-colors"
            />
            {emailMode === 'signup' && (
              <p className="text-[#737373] text-xs">
                You'll choose a username after confirming your email.
              </p>
            )}
            {error && <p className="text-red-400 text-xs">{error}</p>}
            <button
              type="submit"
              disabled={loading}
              className="w-full py-2.5 text-xs font-grotesk font-semibold uppercase tracking-widest bg-white text-black disabled:opacity-30 disabled:cursor-not-allowed hover:bg-[#e0e0e0] transition-colors"
            >
              {loading ? 'Please wait…' : emailMode === 'login' ? 'Sign In' : 'Create Account'}
            </button>

            <div className="flex items-center justify-between pt-1">
              <button
                type="button"
                onClick={() => { setShowEmailForm(false); setError(null); }}
                className="text-[10px] uppercase tracking-widest text-[#404040] hover:text-[#737373] transition-colors"
              >
                ← Back
              </button>
              <button
                type="button"
                onClick={() => { setEmailMode(m => m === 'login' ? 'signup' : 'login'); setError(null); }}
                className="text-[10px] uppercase tracking-widest text-[#404040] hover:text-[#737373] transition-colors"
              >
                {emailMode === 'login' ? 'Create account' : 'Sign in instead'}
              </button>
            </div>
          </form>
        )}
      </div>
    </div>
  );
}
