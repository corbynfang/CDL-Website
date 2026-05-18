import { lazy, Suspense } from "react";
import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import ErrorBoundary from "./components/ErrorBoundary";
import Layout from "./components/Layout";

// Layout, ErrorBoundary, and NotFound stay eagerly loaded — they're tiny
// and needed on every route before any page renders.
import NotFound from "./components/NotFound";

const Home          = lazy(() => import("./components/Home"));
const Players       = lazy(() => import("./components/Players"));
const PlayerDetail  = lazy(() => import("./components/PlayerDetail"));
const Teams         = lazy(() => import("./components/Teams"));
const TeamDetail    = lazy(() => import("./components/TeamDetail"));
const Stats         = lazy(() => import("./components/Stats"));
const Transfers     = lazy(() => import("./components/Transfers"));
const MatchDetail   = lazy(() => import("./components/MatchDetail"));
const EventsPage    = lazy(() => import("./components/EventsPage"));
const EventDetailPage = lazy(() => import("./components/EventDetailPage"));

function App() {
  return (
    <ErrorBoundary>
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
              <Suspense><EventsPage /></Suspense>
            } />
            <Route path="events/:slug" element={
              <Suspense><EventDetailPage /></Suspense>
            } />
            <Route path="*" element={<NotFound />} />
          </Route>
        </Routes>
      </Router>
    </ErrorBoundary>
  );
}

export default App;
