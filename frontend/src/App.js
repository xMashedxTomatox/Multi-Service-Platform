import React, { useState } from "react";
import AuthForm from "./AuthForm";
import Feedback from "./Feedback";
import {jwtDecode} from "jwt-decode";

function App() {
  const [token, setToken] = useState(null);
  const [userId, setUserId] = useState(null);

  const handleLogin = (jwtToken) => {
    setToken(jwtToken);
    try {
      const decoded = jwtDecode(jwtToken);
      setUserId(decoded.sub); // backend sets sub = user_id
    } catch (err) {
      console.error("Failed to decode token:", err);
    }
  };

  const handleLogout = () => {
    setToken(null);
    setUserId(null);
  };

  return (
    <div>
      {!token ? (
        <AuthForm setToken={handleLogin} />
      ) : (
        <div>
          <p>Logged in as user ID: {userId}</p>
          <button onClick={handleLogout} style={{ marginBottom: "10px" }}>
            Sign Out
          </button>
          <Feedback token={token} />
        </div>
      )}
    </div>
  );
}

export default App;
