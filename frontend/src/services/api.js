// src/services/api.js

// RED TEAM ATTACK SIMULATION

export const simulateSQLInjection = async () => {
  const response = await fetch(`${process.env.REACT_APP_RED}/sqlinjection`, {
    method: 'POST',
  });
  if (!response.ok) {
    throw new Error('Failed to simulate SQL Injection attack');
  }
  return response.text();
};

export const simulatePortScan = async () => {
  const response = await fetch(`${process.env.REACT_APP_RED}/portscan`, {
    method: 'POST',
  });
  if (!response.ok) {
    throw new Error('Failed to simulate Port Scan attack');
  }
  return response.text();
};

export const simulateMaliciousUpload = async (formData) => {
  try {
    const response = await fetch(`${process.env.REACT_APP_RED}/maliciousupload`, {
      method: 'POST',
      body: formData,
      
    });

    if (!response.ok) {
      const errorText = await response.text();
      throw new Error(errorText || 'Upload failed');
    }
    return await response.text();
  } catch (error) {
    console.error('Upload error:', error);
    throw error;
  }
};

export const simulateDDoS = async () => {
  const response = await fetch(`${process.env.REACT_APP_RED}/ddos`, {
    method: 'POST',
  });
  if (!response.ok) {
    throw new Error('Failed to simulate DDoS attack');
  }
  return response.text();
};

export const simulateSteganography = async (formData) => {
  const response = await fetch(`${process.env.REACT_APP_RED}/steganography`, {
    method: 'POST',
    body: formData,
  });
  if (!response.ok) {
    throw new Error('Failed to simulate Steganography attack');
  }
  return response.text();
};

// BLUE TEAM DEFENSE SIMULATION

export const defendSQLInjection = async () => {
  const response = await fetch(`${process.env.REACT_APP_BLUE}/defend-sqlinjection`, {
    method: 'POST',
  });
  if (!response.ok) {
    throw new Error('Failed to defend against SQL Injection');
  }
  return response.text();
};

export const defendPortScan = async () => {
  const response = await fetch(`${process.env.REACT_APP_BLUE}/defend-portscan`, {
    method: 'POST',
  });
  if (!response.ok) {
    throw new Error('Failed to defend against Port Scan');
  }
  return response.text();
};

export const defendMaliciousUpload = async (formData) => {
  const response = await fetch(`${process.env.REACT_APP_BLUE}/defend-maliciousupload`, {
    method: 'POST',
    body: formData,
  });
  if (!response.ok) {
    throw new Error('Failed to defend against Malicious Upload');
  }
  return response.text();
};

export const defendDDoS = async () => {
  const response = await fetch(`${process.env.REACT_APP_BLUE}/api/blueteam/ddos/defend`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
  });
  if (!response.ok) {
    throw new Error('Failed to defend against DDoS');
  }
  return response.text();
};


export const defendSteganography = async (formData) => {
  const response = await fetch(`${process.env.REACT_APP_BLUE}/defend-steganography`, {
    method: 'POST',
    body: formData,
  });
  if (!response.ok) {
    throw new Error('Failed to defend against Steganography');
  }
  return response.text();
};

// COMMON FETCHES

export const fetchAlerts = async () => {
  const response = await fetch(`${process.env.REACT_APP_API_BASE}/api/alerts`);
  if (!response.ok) {
    throw new Error('Failed to fetch alerts');
  }
  return response.json();
};

export const fetchBlueAlerts = async () => {
  const response = await fetch(`${process.env.REACT_APP_BLUE}/api/alerts/blue`);
  if (!response.ok) throw new Error('Failed to fetch blue team alerts');
  return response.json();
};

export const fetchRedAlerts = async () => {
  const response = await fetch(`${process.env.REACT_APP_RED}/api/alerts/red`);
  if (!response.ok) throw new Error('Failed to fetch red team alerts');
  return response.json();
};

export const fetchBlockedIPs = async () => {
  const response = await fetch(`${process.env.REACT_APP_BLUE}/api/blocked`);
  if (!response.ok) {
    throw new Error('Failed to fetch blocked IPs');
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
