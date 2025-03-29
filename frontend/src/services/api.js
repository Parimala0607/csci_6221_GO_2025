// src/services/api.js

export const simulateDefense = async () => {
    const response = await fetch(`${process.env.REACT_APP_BLUE}/defend`, {
      method: 'POST',
    });
  
    if (!response.ok) {
      throw new Error('Failed to simulate defense');
    }
  
    return response.text();
  };
  
  // src/services/api.js

export const simulateAttack = async () => {
    const response = await fetch(`${process.env.REACT_APP_RED}/attack`, {
      method: 'POST',
    });
  
    if (!response.ok) {
      throw new Error('Failed to simulate attack');
    }
  
    return response.text();
  };


  // src/services/api.js

export const fetchAlerts = async () => {
    const response = await fetch(`${process.env.REACT_APP_API_BASE}/api/alerts`);
    if (!response.ok) {
      throw new Error('Failed to fetch alerts');
    }
    return response.json();
  };
  
  export const fetchLogs = async () => {
    const response = await fetch(`${process.env.REACT_APP_API_BASE}/api/logs`);
    if (!response.ok) {
      throw new Error('Failed to fetch logs');
    }
    return response.json();
  };
  