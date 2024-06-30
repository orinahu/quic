import React from "react";
import HttpsClient from "./Http/Http";
import QuicClient from "./QuicClient/QuicClient";

function App() {
  return (
    <div className="App">
      <h1>WebTransport Client</h1>
      <QuicClient />
      <h1>Https Client</h1>
      <HttpsClient />
    </div>
  );
}

export default App;
