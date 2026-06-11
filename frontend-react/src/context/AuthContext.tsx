import { createContext, useContext, useEffect, useRef, useState } from "react";
import type { ReactNode } from "react";
import axios from "axios";
import type { Session, User } from "@supabase/supabase-js";
import { supabase } from "../lib/supabaseClient";
import api from "../services/api";

type OAuthProvider = "github" | "twitch" | "google";

interface AuthContextValue {
  session: Session | null;
  user: User | null;
  loading: boolean;
  showAuthModal: boolean;
  needsProfileSetup: boolean;
  openAuthModal: () => void;
  closeAuthModal: () => void;
  signUp: (
    email: string,
    password: string,
    captchaToken?: string,
  ) => Promise<{ error: string | null }>;
  signIn: (
    email: string,
    password: string,
  ) => Promise<{ error: string | null }>;
  signInWithOAuth: (
    provider: OAuthProvider,
  ) => Promise<{ error: string | null }>;
  completeProfileSetup: (username: string) => Promise<{ error: string | null }>;
  signOut: () => Promise<void>;
}

const AuthContext = createContext<AuthContextValue | null>(null);

async function checkProfile(token: string): Promise<boolean> {
  try {
    await api.get("/auth/me", {
      headers: { Authorization: `Bearer ${token}` },
    });
    return true;
  } catch (err: unknown) {
    if (axios.isAxiosError(err) && err.response?.status === 422) {
      return false;
    }
    return true;
  }
}

export function AuthProvider({ children }: { children: ReactNode }) {
  const [session, setSession] = useState<Session | null>(null);
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);
  const [showAuthModal, setShowAuthModal] = useState(false);
  const [needsProfileSetup, setNeedsProfileSetup] = useState(false);
  const sessionRef = useRef<Session | null>(null);

  useEffect(() => {
    supabase.auth.getSession().then(async ({ data: { session } }) => {
      sessionRef.current = session;
      setSession(session);
      setUser(session?.user ?? null);
      if (session) {
        const hasProfile = await checkProfile(session.access_token);
        setNeedsProfileSetup(!hasProfile);
      }
      setLoading(false);
    });

    const {
      data: { subscription },
    } = supabase.auth.onAuthStateChange(async (event, session) => {
      sessionRef.current = session;
      setSession(session);
      setUser(session?.user ?? null);

      if (event === "SIGNED_IN" && session) {
        const hasProfile = await checkProfile(session.access_token);
        setNeedsProfileSetup(!hasProfile);
        setShowAuthModal(false);
      }
      if (event === "SIGNED_OUT") {
        setNeedsProfileSetup(false);
      }
    });

    return () => subscription.unsubscribe();
  }, []);

  useEffect(() => {
    const interceptor = api.interceptors.request.use((config) => {
      const token = sessionRef.current?.access_token;
      if (token) {
        config.headers.set("Authorization", `Bearer ${token}`);
      }
      return config;
    });
    return () => api.interceptors.request.eject(interceptor);
  }, []);

  function openAuthModal() {
    setShowAuthModal(true);
  }
  function closeAuthModal() {
    setShowAuthModal(false);
  }

  async function signUp(
    email: string,
    password: string,
    captchaToken?: string,
  ) {
    const { error } = await supabase.auth.signUp({
      email,
      password,
      options: captchaToken ? { captchaToken } : undefined,
    });
    return { error: error?.message ?? null };
  }

  async function signIn(email: string, password: string) {
    const { error } = await supabase.auth.signInWithPassword({
      email,
      password,
    });
    return { error: error?.message ?? null };
  }

  async function signInWithOAuth(provider: OAuthProvider) {
    const { error } = await supabase.auth.signInWithOAuth({
      provider,
      options: { redirectTo: window.location.origin },
    });
    return { error: error?.message ?? null };
  }

  async function completeProfileSetup(username: string) {
    try {
      const {
        data: { session: fresh },
      } = await supabase.auth.getSession();
      await api.post(
        "/auth/profile",
        { username },
        {
          headers: fresh?.access_token
            ? { Authorization: `Bearer ${fresh.access_token}` }
            : {},
        },
      );
      setNeedsProfileSetup(false);
      return { error: null };
    } catch (err: unknown) {
      const msg = axios.isAxiosError<{ error: string }>(err)
        ? (err.response?.data?.error ?? "Failed to set up profile")
        : "Failed to set up profile";
      return { error: msg };
    }
  }

  async function signOut() {
    await supabase.auth.signOut();
  }

  return (
    <AuthContext.Provider
      value={{
        session,
        user,
        loading,
        showAuthModal,
        needsProfileSetup,
        openAuthModal,
        closeAuthModal,
        signUp,
        signIn,
        signInWithOAuth,
        completeProfileSetup,
        signOut,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
}

// eslint-disable-next-line react-refresh/only-export-components
export function useAuth() {
  const ctx = useContext(AuthContext);
  if (!ctx) throw new Error("useAuth must be used inside AuthProvider");
  return ctx;
}
