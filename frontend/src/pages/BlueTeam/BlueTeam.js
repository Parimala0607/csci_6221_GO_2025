import React, { useState } from 'react';
import { simulateDefense } from '../../services/api'; //  Correct path

function BlueTeam() {
  const [message, setMessage] = useState('');

  const handleDefense = async () => {
    try {
      const result = await simulateDefense();
      setMessage(result);
    } catch (error) {
      setMessage('Error logging defense: ' + error.message);
    }
  };

  return (
    <div className="blue-team">
      <h1>Blue Team Operations</h1>
      <button onClick={handleDefense}>Simulate Defense Action</button>
      {message && <div className="result-message">{message}</div>}
    </div>
  );
}

export default BlueTeam;
