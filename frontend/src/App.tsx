import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import Layout from './components/Layout';
import Home from './components/Home';
import Teams from './components/Teams';
import Players from './components/Players';
import KDStats from './components/KDStats';
import PlayerKDStats from './components/PlayerKDStats';
import PlayerDetails from './components/PlayerDetails';


function App() {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<Layout />}>
          <Route index element={<Home />} />
          <Route path="teams" element={<Teams />} />
          <Route path="players" element={<Players />} />
          <Route path="players/:id" element={<PlayerDetails />} />
          <Route path="kd-stats" element={<KDStats />} />
          <Route path="players/:id/kd-stats" element={<PlayerKDStats />} />

          {/* Add more routes as needed */}
        </Route>
      </Routes>
    </Router>
  );
}

export default App;
