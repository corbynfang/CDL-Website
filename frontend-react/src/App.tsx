import { lazy, Suspense } from "react";
import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import { AuthProvider } from "./context/AuthContext";
import ErrorBoundary from "./components/ErrorBoundary";
import Layout from "./components/Layout";
import NotFound from "./components/NotFound";

const Home          = lazy(() => import("./components/Home"));
const Players       = lazy(() => import("./components/Players"));
const PlayerDetail  = lazy(() => import("./components/players/PlayerDetail"));
const Teams         = lazy(() => import("./components/Teams"));
const TeamDetail    = lazy(() => import("./components/teams/TeamDetail"));
const Stats         = lazy(() => import("./components/Stats"));
const Transfers     = lazy(() => import("./components/Transfers"));
const MatchDetail   = lazy(() => import("./components/matches/MatchDetail"));
const Events      = lazy(() => import("./components/Events"));
const EventDetail = lazy(() => import("./components/EventDetail"));
const PrivacyPage     = lazy(() => import("./components/legal/PrivacyPage"));
const TermsPage       = lazy(() => import("./components/legal/TermsPage"));
const DisclaimerPage  = lazy(() => import("./components/legal/DisclaimerPage"));

function App() {
  return (
    <ErrorBoundary>
      <AuthProvider>
      <Router>
        <Routes>
          <Route path="/" element={<Layout />}>
            <Route index element={
              <Suspense><Home /></Suspense>
            } />
            <Route path="players" element={
              <Suspense><Players /></Suspense>
            } />
            <Route path="players/:id" element={
              <Suspense><PlayerDetail /></Suspense>
            } />
            <Route path="teams" element={
              <Suspense><Teams /></Suspense>
            } />
            <Route path="teams/:id" element={
              <Suspense><TeamDetail /></Suspense>
            } />
            <Route path="matches/:id" element={
              <Suspense><MatchDetail /></Suspense>
            } />
            <Route path="stats" element={
              <Suspense><Stats /></Suspense>
            } />
            <Route path="transfers" element={
              <Suspense><Transfers /></Suspense>
            } />
            <Route path="events" element={
              <Suspense><Events /></Suspense>
            } />
            <Route path="events/:slug" element={
              <Suspense><EventDetail /></Suspense>
            } />
            <Route path="privacy"    element={<Suspense><PrivacyPage /></Suspense>} />
            <Route path="terms"      element={<Suspense><TermsPage /></Suspense>} />
            <Route path="disclaimer" element={<Suspense><DisclaimerPage /></Suspense>} />
            <Route path="*" element={<NotFound />} />
          </Route>
        </Routes>
      </Router>
      </AuthProvider>
    </ErrorBoundary>
  );
}

export default App;
