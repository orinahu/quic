import React, { useState } from 'react';

const HTTPS_URL = 'https://localhost:4433/your-endpoint';

export const HttpsClient = () => {
  const [message, setMessage] = useState('');
  const [response, setResponse] = useState('');

  const sendMessage = async () => {
    try {
      const response = await fetch(HTTPS_URL, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ message }),
      });

      const data = await response.json();
      setResponse(data.response);
    } catch (error) {
      console.error('Error sending message:', error);
    }
  };

  return (
    <div>
      <input
        type="text"
        value={message}
        onChange={(e) => setMessage(e.target.value)}
        placeholder="Enter message"
      />
      <button onClick={sendMessage}>Send Message</button>
      {response && <p>Response from server: {response}</p>}
    </div>
  );
};

export default HttpsClient;
