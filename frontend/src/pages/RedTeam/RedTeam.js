import React, { useState } from 'react';
import { simulateAttack } from '../../services/api'; 

function RedTeam() {
  const [message, setMessage] = useState('');

  const handleAttack = async () => {
    try {
      const result = await simulateAttack();
      setMessage(result);
    } catch (error) {
      setMessage('Error simulating attack: ' + error.message);
    }
  };

  return (
    <div className="red-team">
      <h1>Red Team Operations</h1>
      <button onClick={handleAttack}>Simulate Attack</button>
      {message && <div className="result-message">{message}</div>}
    </div>
  );
}

export default RedTeam;
