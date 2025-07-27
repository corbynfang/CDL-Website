import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import Layout from './components/Layout';
import Home from './components/Home';
import Teams from './components/Teams';
import TeamDetails from './components/TeamDetails';
import TeamPlayers from './components/TeamPlayers';
import Players from './components/Players';
import KDStats from './components/KDStats';
import PlayerKDStats from './components/PlayerKDStats';
import PlayerDetails from './components/PlayerDetails';
import Transfers from './components/Transfers';


function App() {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<Layout />}>
          <Route index element={<Home />} />
          <Route path="teams" element={<Teams />} />
          <Route path="teams/:id" element={<TeamDetails />} />
          <Route path="teams/:id/players" element={<TeamPlayers />} />
          <Route path="players" element={<Players />} />
          <Route path="players/:id" element={<PlayerDetails />} />
          <Route path="kd-stats" element={<KDStats />} />
          <Route path="players/:id/kd-stats" element={<PlayerKDStats />} />
          <Route path="transfers" element={<Transfers />} />

          {/* Add more routes as needed */}
        </Route>
      </Routes>
    </Router>
  );
}

export default App;
