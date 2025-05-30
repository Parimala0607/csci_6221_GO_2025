import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import React from 'react';
import Dashboard from './pages/Dashboard/Dashboard';
import RedTeam from './pages/RedTeam/RedTeam';
import BlueTeam from './pages/BlueTeam/BlueTeam';

import './styles/App.css';

function App() {
  return (
    <Router>
      <div className="App">
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
