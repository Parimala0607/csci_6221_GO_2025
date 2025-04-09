import React, { useEffect, useState } from 'react';
import { fetchBlueAlerts } from '../../services/api';
import { fetchBlockedIPs } from '../../services/api';
import { Link } from 'react-router-dom';
import './BlueTeam.css';



function BlueTeam() {
  const [message, setMessage] = useState('');
  const [alerts, setAlerts] = useState([]);
  const [blockedIPs, setBlockedIPs] = useState([]);

  // Defense functions
  const defendSQLInjection = async () => {
    try {
      const response = await fetch(`${process.env.REACT_APP_BLUE}/defend-sqlinjection`, {
        method: 'POST',
      });
      if (!response.ok) throw new Error('SQL Injection defense failed');
      const result = await response.text();
      setMessage(result);
      await fetchAlerts();
    } catch (error) {
      setMessage('Error: ' + error.message);
    }
  };

  const defendPortScan = async () => {
    try {
      const response = await fetch(`${process.env.REACT_APP_BLUE}/defend-portscan`, {  
        method: 'POST',
      });
      if (!response.ok) throw new Error('Port Scan defense failed');
      const result = await response.text();
      setMessage(result);
      await fetchAlerts();
    } catch (error) {
      setMessage('Error: ' + error.message);
    }
  };

  const defendDDoS = async () => {
    try {
      const response = await fetch(`${process.env.REACT_APP_BLUE}/api/blueteam/ddos/defend`, {
        method: 'POST',
      });
  
      if (!response.ok) throw new Error('DDoS defense failed');
  
      const contentType = response.headers.get('content-type');
  
      let message = '';
      if (contentType && contentType.includes('application/json')) {
        const result = await response.json();
        
        message = `Defended ${result.total_alerts_defended || 0} alerts.\n` +
                  `Blocked IPs: ${(result.ips_blocked || []).join(', ') || 'None'}\n` +
                  `Defended IPs: ${(result.ips_defended || []).join(', ') || 'None'}`;
      } else {
        message = await response.text(); 
      }
  
      setMessage(message);
      await fetchAlerts(); 
    } catch (error) {
      setMessage('Error: ' + error.message);
    }
  };

const defendSteganography = async () => {
    try {
        const response = await fetch(`${process.env.REACT_APP_BLUE}/defend-steganography`, {
            method: 'POST',
        });
        
        const result = await response.json();
        
        setMessage([
            `${result.status.toUpperCase()}`,
            `File: ${result.filename}`,
            result.hidden_data ? `Hidden Data: ${result.hidden_data}` : "No hidden data found",
            `Source IP: ${result.source_ip}`
        ].join('\n'));
        
        await fetchAlerts();
    } catch (error) {
        setMessage(`Defense failed: ${error.message}`);
    }
};

const defendMaliciousUpload = async () => {
  try {
    setMessage("Analyzing files...");
    
    const response = await fetch(`${process.env.REACT_APP_BLUE}/defend-maliciousupload`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Accept': 'application/json'
      }
    });

    if (!response.ok) {
      const errorText = await response.text();
      throw new Error(errorText || `HTTP ${response.status}`);
    }

    const result = await response.json();
    
    // Build detailed defense report
    let report = `Defense Report:\n`;
    report += `File: ${result.filename}\n`;
    report += `Size: ${result.file_size} bytes\n`;
    report += `Source IP: ${result.source_ip}\n`;
    report += `Found: ${result.keywords_found?.join(', ') || 'None'}\n`;
    report += `Content:\n${result.file_content}\n\n`;
    
    if (result.blocked) {
      report += `ACTION TAKEN: ${result.block_message}`;
    } else {
      report += `File cleared - no malicious content detected`;
    }

    setMessage(report);
    await fetchAlerts(); 

  } catch (error) {
    setMessage(`Defense failed: ${error.message}`);
    console.error("Defense error:", error);
  }
};

 
  const fetchAlerts = async () => {
    try {
      const data = await fetchBlueAlerts();
      setAlerts(data);
    } catch (error) {
      setMessage('Could not load alerts: ' + error.message);
    }
  };

 
  const fetchblockIPS = async () => {
    try {
      const data = await fetchBlockedIPs();
      setBlockedIPs(data);
    } catch (error) {
      setMessage('Could not load blocked IPs: ' + error.message);
    }
  };

 
  useEffect(() => {
    let ws;
    let reconnectAttempts = 0;
    const maxReconnectAttempts = 5;

    const connectWebSocket = () => {
      const wsUrl = process.env.REACT_APP_WS_URL || 'ws://localhost:8081';
      ws = new WebSocket(`${wsUrl}/ws`);

      ws.onmessage = (event) => {
        console.log("New alert received:", event.data);
        fetchAlerts(); 
        
      };
      


      ws.onopen = () => {
        console.log("WebSocket connected to Blue Team backend");
        reconnectAttempts = 0; 
      };
          
      ws.onerror = (error) => {
        console.error("WebSocket error:", error);
      };

      ws.onclose = () => {
        console.log("WebSocket connection closed");
        if (reconnectAttempts < maxReconnectAttempts) {
          const delay = Math.min(1000 * Math.pow(2, reconnectAttempts), 30000);
          console.log(`Attempting reconnect in ${delay}ms...`);
          setTimeout(connectWebSocket, delay);
          reconnectAttempts++;
        }
      };
    };

    
    // Initial fetch and WebSocket connection
    fetchAlerts().catch(error => {
      console.error('Initial alerts fetch failed:', error);
    });
    connectWebSocket();

    return () => {
      if (ws) {
        ws.close();
      }
    };
  }, []);

  useEffect(() => {
    fetchblockIPS();
  }, []);


  return (
    <div className="blue-container">
      <aside className="sidebar">
        <h2>Blue Team</h2>
        <nav>
          <ul>
            <li><Link to="/">Dashboard</Link></li>
            <li><Link to="/red">Red Team</Link></li>
            <li className="active"><Link to="/blue">Blue Team</Link></li>
          </ul>
        </nav>
      </aside>

      <main className="main-contentB">
        <header className="top-barB">
          <h1>Blue Team Control Center</h1>
        </header>

        <section className="stats-gridB">
          <div className="Bcard">
            <h3>Active Threats</h3>
            <p>{
      new Set(
        alerts
          .filter(a => a.status !== 'defended')
          .map(a => a.source_ip)
      ).size
    }</p>
          </div>
          <div className="Bcard">
            <h3>Blocked Attacks</h3>
            <p>{
      new Set(
        alerts
          .filter(a => a.status === 'defended')
          .map(a => a.source_ip)
      ).size
    }</p>
          </div>
          <div className="Bcard">
  <h3>Blocked IPs</h3>
  <p>{blockedIPs.length}</p>
</div>
        </section>

        <section className="actionsB">
          <h2>Defensive Actions</h2>
          <div className="button-groupB">
            <button onClick={defendSQLInjection}>Defend SQL Injection</button>
            <button onClick={defendDDoS}>Defend DDoS</button>
            <button onClick={defendSteganography}>Defend Steganography</button>
            <button onClick={defendMaliciousUpload}>Defend Malicious Upload</button>
            <button onClick={defendPortScan}>Defend Port Scan</button>
          </div>
          

          {message && (
  <p className="status-message" style={{ color: '#000000', fontWeight: 'bold' }}>
    {message}
  </p>
)}
        </section>

  

        <section className="recent-alerts team-boxB">
  <h2>Current Alerts</h2>
  {alerts.length ? (
    alerts
      .sort((a, b) => new Date(b.timestamp) - new Date(a.timestamp)) 
      .slice(0, 15)
      .map((alert) => (
        <div className="alert activityB" key={alert.id}>
          <strong>{alert.message}</strong>
          <p>Source IP: {alert.source_ip}</p>
          <p>Severity: {alert.severity}</p>
          <span style={{ color: alert.status === 'defended' ? 'green' : 'red' }}>
            {alert.status === 'defended' ? 'Defended' : 'Not Defended'}
          </span>
        </div>
      ))
  ) : (
    <p>No current alerts.</p>
  )}
</section>
      </main>
    </div>
  );
}

export default BlueTeam;