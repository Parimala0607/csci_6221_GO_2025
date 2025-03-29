import { BrowserRouter as Router, Routes, Route, Link } from 'react-router-dom';

import Dashboard from './pages/Dashboard/Dashboard';
import RedTeam from './pages/RedTeam/RedTeam';
import BlueTeam from './pages/BlueTeam/BlueTeam';

import './styles/App.css';

function App() {
  return (
    <Router>
      <div className="App">
        <nav className="navigation">
          <Link to="/">Dashboard</Link>
          <Link to="/red">Red Team</Link>
          <Link to="/blue">Blue Team</Link>
        </nav>

        <Routes>
          <Route path="/" element={<Dashboard />} />
          <Route path="/red" element={<RedTeam />} />
          <Route path="/blue" element={<BlueTeam />} />
        </Routes>
      </div>
    </Router>
  );
}

export default App;
