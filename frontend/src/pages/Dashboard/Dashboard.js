import React, { useEffect, useState } from 'react';
import { fetchAlerts, fetchLogs } from '../../services/api';

function Dashboard() {
  const [alerts, setAlerts] = useState([]);
  const [logs, setLogs] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  const fetchData = async () => {
    try {
      const [alertsData, logsData] = await Promise.all([
        fetchAlerts(),
        fetchLogs()
      ]);

      setAlerts(alertsData);
      setLogs(logsData);
      setError(null);
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchData();
    const interval = setInterval(fetchData, 5000);
    return () => clearInterval(interval);
  }, []);

  if (loading) return <div className="loading">Loading dashboard data...</div>;
  if (error) return <div className="error">Error: {error}</div>;

  return (
    <div className="dashboard">
      <h1>Security Dashboard</h1>
      <div className="grid-container">
        <div className="alerts-section">
          <h2>Recent Alerts ({alerts.length})</h2>
          <div className="alert-list">
            {alerts.map(alert => (
              <div key={alert.ID} className={`alert severity-${alert.Severity}`}>
                <div className="alert-header">
                  <span className="timestamp">{alert.Timestamp}</span>
                  <span className="severity">{alert.Severity}</span>
                  <span className="source-ip">{alert.SourceIP}</span>
                </div>
                <div className="alert-body">{alert.Message}</div>
              </div>
            ))}
          </div>
        </div>

        <div className="logs-section">
          <h2>Defense Logs ({logs.length})</h2>
          <div className="log-list">
            {logs.map(log => (
              <div key={log.ID} className="log-entry">
                <div className="log-header">
                  <span className="timestamp">{log.Timestamp}</span>
                  <span className="action">{log.Action}</span>
                  <span className="source-ip">{log.SourceIP}</span>
                </div>
                <div className="log-body">{log.Description}</div>
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
  );
}

export default Dashboard;
