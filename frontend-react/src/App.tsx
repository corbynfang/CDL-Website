import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import ErrorBoundary from "./components/ErrorBoundary";
import Layout from "./components/Layout";
import Home from "./components/Home";
import Players from "./components/Players";
import PlayerDetail from "./components/PlayerDetail";
import Teams from "./components/Teams";
import TeamDetail from "./components/TeamDetail";
import Stats from "./components/Stats";
import Transfers from "./components/Transfers";
import MatchDetail from "./components/MatchDetail";
import NotFound from "./components/NotFound";

function App() {
  return (
    <ErrorBoundary>
      <Router>
        <Routes>
          <Route path="/" element={<Layout />}>
            <Route index element={<Home />} />
            <Route path="players" element={<Players />} />
            <Route path="players/:id" element={<PlayerDetail />} />
            <Route path="teams" element={<Teams />} />
            <Route path="teams/:id" element={<TeamDetail />} />
            <Route path="matches/:id" element={<MatchDetail />} />
            <Route path="stats" element={<Stats />} />
            <Route path="transfers" element={<Transfers />} />
            <Route path="*" element={<NotFound />} />
          </Route>
        </Routes>
      </Router>
    </ErrorBoundary>
  );
}

export default App;
