import React, { useState, useEffect } from 'react';

import {
  simulateSQLInjection,
  simulatePortScan,
  simulateMaliciousUpload,
  simulateDDoS,
  simulateSteganography,
  fetchRedAlerts
} from '../../services/api';
import { Link } from 'react-router-dom';
import './RedTeam.css';

function RedTeam() {
  const [message, setMessage] = useState('');
  const [file, setFile] = useState(null);      
  const [malFile, setMalFile] = useState(null); 
  const [alerts, setAlerts] = useState([]);

  const handleSQLInjection = async () => {
    try {
      const result = await simulateSQLInjection();
      setMessage(result);
    } catch (error) {
      setMessage('Error simulating SQL Injection: ' + error.message);
    }
  };

  const handlePortScan = async () => {
    try {
      const result = await simulatePortScan();
      setMessage(result);
    } catch (error) {
      setMessage('Error simulating Port Scan: ' + error.message);
    }
  };

  const handleMaliciousUpload = async () => {
    if (!malFile) {
      setMessage("Please select a file first");
      return;
    }
  
    const formData = new FormData();
    formData.append('file', malFile); 
  
    try {
      const result = await simulateMaliciousUpload(formData);
      setMessage(`Upload successful: ${result}`);
      setMalFile(null);
      
      const fileInput = document.querySelector('input[type="file"][name="file"]');
      if (fileInput) fileInput.value = "";
    } catch (error) {
      setMessage(`Upload failed: ${error.message}`);
    }
  };

  const handleDDoS = async () => {
    try {
      const result = await simulateDDoS();  
      setMessage(result);
    } catch (error) {
      setMessage('Error simulating DDoS: ' + error.message);
    }
  };


  const fetchAlerts = async () => {
      try {
        const data = await fetchRedAlerts();
        setAlerts(data);
      } catch (error) {
        setMessage('Could not load alerts: ' + error.message);
      }
    };

  useEffect(() => {
    fetchAlerts();
  }, []);

  

  const handleSteganography = async () => {
    if (!file) {
        setMessage("Please select an image first");
        return;
    }

    const formData = new FormData();
    formData.append('image', file);

    try {
        const response = await fetch(`${process.env.REACT_APP_RED}/steganography`, {
            method: 'POST',
            body: formData,
        });
        
        const result = await response.json();
        setMessage(`${result.message} | File: ${result.filename}`);
    } catch (error) {
        setMessage(`Error: ${error.message}`);
    }
};

const systemsBreached = new Set(
  alerts
    .filter(a => a.status?.toLowerCase() !== "defended")
    .map(a => a.source_ip)
).size;

const attacksDeployed = new Set(
  alerts.map(a => a.source_ip || a.sourceIP)
).size;

const exploitKeywords = ["sql", "ddos", "steganography", "malicious", "port"];

const exploitsFound = new Set(
  alerts.flatMap(a =>
    exploitKeywords.filter(kw => a.message.toLowerCase().includes(kw))
  )
).size;


 return (
    <div className="dashboard-container">
      <aside className="sidebar">
        <h2>Red Team</h2>
        <nav>
          <ul>
            <li><Link to="/">Dashboard</Link></li>
            <li className="active"><Link to="/red">Red Team</Link></li>
            <li><Link to="/blue">Blue Team</Link></li>
          </ul>
        </nav>
      </aside>

      <main className="main-contentR">
        <header className="top-barR">
          <h1>Red Team Command Center</h1>
        </header>

        <section className="stats-gridR">
          <div className="Rcard">
            <h3>Attacks Deployed</h3>
            <p>{attacksDeployed}</p>
          </div>
          <div className="Rcard">
            <h3>Exploits Found</h3>
            <p>{exploitsFound}</p>
          </div>
          <div className="Rcard">
            <h3>Systems Breached</h3>
            <p>{systemsBreached}</p>
          </div>
        </section>

        <section className="actionsR">
          <h2>Offensive Actions</h2>
          <div className="button-groupR">
            <button onClick={handleSQLInjection}>SQL Injection</button>
            <button onClick={handleDDoS}>DDoS Attack</button>
            <button onClick={handleSteganography} disabled={!file}>Steganography</button>
            <button onClick={handleMaliciousUpload} disabled={!malFile}>Malicious Upload</button>
            <button onClick={handlePortScan}>Port Scan Attack</button>
          </div>

          {/* Upload inputs */}
          <div className="upload-section">
            <div style={{ marginTop: '1em' }}>
              <label htmlFor="malicious-upload"><strong>Malicious File:</strong></label>
              <input
                type="file"
                name="file"
                id="malicious-upload"
                onChange={(e) => setMalFile(e.target.files?.[0])}
              />
              {malFile && (
                <p>
                  Selected file: {malFile.name} ({Math.round(malFile.size / 1024)} KB)
                </p>
              )}
            </div>

            <div style={{ marginTop: '1em' }}>
              <label htmlFor="stego-upload"><strong>Steganography Image:</strong></label>
              <input
                type="file"
                id="stego-upload"
                accept="image/*"
                onChange={(e) => setFile(e.target.files?.[0])}
              />
              {file && <p>Selected image: {file.name}</p>}
            </div>
          </div>

          {/* Feedback Message */}
          {message && (
  <p className="status-message" style={{ color: '#000000', fontWeight: 'bold' }}>
    {message}
  </p>
)}
        </section>
      </main>
    </div>
  );
}


export default RedTeam;
