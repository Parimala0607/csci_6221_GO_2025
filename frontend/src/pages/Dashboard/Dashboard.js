import React, { useEffect, useState } from 'react';
import { fetchAlerts, fetchLogs } from '../../services/api';
import { NavLink } from 'react-router-dom';
import './Dashboard.css';
import {
  BarChart, Bar, XAxis, YAxis, Tooltip, Legend, ResponsiveContainer,
  PieChart, Pie, Cell
} from 'recharts';

const CHART_COLORS = ["#3b82f6", "#10b981", "#f59e0b", "#ef4444", "#8b5cf6", "#0ea5e9"];

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
      setAlerts(alertsData || []);
      setLogs(logsData || []);
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

  const classifyAlertType = (message) => {
    if (message.toLowerCase().includes("sql")) return "SQL Injection";
    if (message.toLowerCase().includes("ddos")) return "DDoS";
    if (message.toLowerCase().includes("steganography")) return "Steganography";
    if (message.toLowerCase().includes("malicious")) return "Malicious Upload";
    if (message.toLowerCase().includes("port")) return "Port Scan";
    return "Other";
  };

  //const totalAlerts = alerts.length;

  const totalAlerts = new Set(
    alerts
      .filter(alert => alert.Status && alert.Status.toLowerCase() === "not_defended")
      .map(alert => alert.SourceIP)
  ).size;
//const uniqueVulnSources = new Set(alerts.map(a => a.SourceIP)).size;

const uniqueVulnSources = new Set(
  alerts
    .map(a => a.SourceIP)
    .filter(ip => ip) 
).size;

const incidentsResolved = new Set(
  alerts
    .filter(alert => alert.Status && alert.Status.toLowerCase() === "defended")
    .map(alert => alert.SourceIP)
).size;

const total = totalAlerts + incidentsResolved;
const rawScore = total > 0
  ? (incidentsResolved / total) * 100
  : 100;
const securityScore = Math.round(Math.min(rawScore, 100)) + "%";

const severityData = ["low", "medium", "high"].map(level => {
  const ipSet = new Set(
    alerts
      .filter(a => a.Severity && a.Severity.toLowerCase() === level)
      .map(a => a.SourceIP)
  );
  return { severity: level, count: ipSet.size };
});

const typeIPMap = new Map();

alerts.forEach(alert => {
  const type = classifyAlertType(alert.Message);
  const ip = alert.SourceIP;

  if (!typeIPMap.has(type)) {
    typeIPMap.set(type, new Set());
  }

  typeIPMap.get(type).add(ip);
});

const typeData = Array.from(typeIPMap.entries()).map(([name, ipSet]) => ({
  name,
  value: ipSet.size
}));

  
const barColors = ["#3b82f6", "#f59e0b", "#ef4444"];

  if (loading) return <div className="loading">Loading dashboard data...</div>;
  if (error) return <div className="error">Error: {error}</div>;

  return (
    <div className="dashboard-container">
      {/* Sidebar Navigation */}
      <aside className="sidebar">
        <h2>Security Dashboard</h2>
        <nav>
          <ul>
            <li className="active"><NavLink to="/">Dashboard</NavLink></li>
            <li><NavLink to="/red">Red Team</NavLink></li>
            <li><NavLink to="/blue">Blue Team</NavLink></li>
          </ul>
        </nav>
      </aside>

      {/* Main Content */}
      <main className="main-contentD">
        {/* Stats Cards */}
        <section className="stats-gridD">
        <div className="Dcard"><h3>Active Threats</h3><p>{totalAlerts}</p></div>
  <div className="Dcard"><h3>Attack Attempts Made</h3><p>{uniqueVulnSources}</p></div>
  <div className="Dcard"><h3>Security Score</h3><p>{securityScore}</p></div>
  <div className="Dcard"><h3>Incidents Resolved</h3><p>{incidentsResolved}</p></div>
        </section>

        {/* Charts Section */}
        <section className="charts-grid">
          {/* Bar Chart: Attacks by Severity */}
          <div className="chart-box">
            <h3>Attacks by Severity</h3>
            <ResponsiveContainer width="100%" height={300}>
              <BarChart data={severityData}>
                <XAxis dataKey="severity" />
                <YAxis />
                <Tooltip />
                <Legend />
                <Bar dataKey="count" name="Severity">
  {severityData.map((entry, index) => (
    <Cell key={`cell-${index}`} fill={barColors[index % barColors.length]} />
  ))}
</Bar>
              </BarChart>
            </ResponsiveContainer>
          </div>

          {/* Pie Chart: Alerts by Type */}
          <div className="chart-box">
            <h3>Alerts by Type</h3>
            <ResponsiveContainer width="100%" height={300}>
              <PieChart>
                <Pie
                  data={typeData}
                  dataKey="value"
                  nameKey="name"
                  cx="50%"
                  cy="50%"
                  outerRadius={100}
                  label
                >
                  {typeData.map((_, index) => (
                    <Cell key={`cell-${index}`} fill={CHART_COLORS[index % CHART_COLORS.length]} />
                  ))}
                </Pie>
                <Tooltip />
                <Legend />
              </PieChart>
            </ResponsiveContainer>
          </div>
        </section>

        {/* Activity Summary */}
        <section className="team-activities">
  {/* Red Team */}
  <div className="team-box">
    <h3>Red Team Activity Summary</h3>
    {alerts.length === 0 && <p>No alerts found.</p>}
    {[...alerts]
      .sort((a, b) => new Date(b.Timestamp) - new Date(a.Timestamp))
      .slice(0, 10)
      .map(alert => (
        <div key={alert.ID} className={`alert-entry activity severity-${alert.Severity}`}>
          <div className="alert-header">
            <span className="timestamp">{alert.Timestamp}</span>
            <span className="severity">{alert.Severity}</span>
            <span className="source-ip">{alert.SourceIP}</span>
          </div>
          <div className="alert-body">{alert.Message}</div>
        </div>
      ))}
  </div>

  {/* Blue Team */}
  <div className="team-box">
    <h3>Blue Team Activity Summary</h3>
    {logs.length === 0 && <p>No defense logs found.</p>}
    {[...logs]
      .sort((a, b) => new Date(b.Timestamp) - new Date(a.Timestamp))
      .slice(0, 10)
      .map(log => (
        <div key={log.ID} className="log-entry activity">
          <div className="log-header">
            <span className="timestamp">{log.Timestamp}</span>
            <span className="action">{log.Action}</span>
            <span className="source-ip">{log.SourceIP}</span>
          </div>
          <div className="log-body">{log.Description}</div>
        </div>
      ))}
  </div>
</section>
      </main>
    </div>
  );
}

export default Dashboard;