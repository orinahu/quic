import React, { useState } from 'react';

const QUIC_URL = 'https://localhost:4433/webtransport';

const QuicClient = () => {
  const [message, setMessage] = useState('');
  const [response, setResponse] = useState('');

  const sendMessage = async () => {
    try {
      const transport = new WebTransport(QUIC_URL);

      await transport.ready;

      const writer = transport.datagrams.writable.getWriter();
      const encoder = new TextEncoder();
      const encodedMessage = encoder.encode(message);

      await writer.write(encodedMessage);
      writer.close();

      const reader = transport.datagrams.readable.getReader();
      const { value, done } = await reader.read();
      if (!done) {
        const decoder = new TextDecoder();
        const decodedResponse = decoder.decode(value);
        setResponse(decodedResponse);
      }

      transport.close();
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

export default QuicClient;
