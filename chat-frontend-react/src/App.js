import React, { useState, useEffect, useRef } from 'react';
import './App.css';
import InputForm from './components/InputForm';

function App() {
  const [view, setView] = useState('login'); // 'login', 'lobby', 'chat'
  const [username, setUsername] = useState('');
  const [passwordHash, setPasswordHash] = useState('');
  const [isRegistering, setIsRegistering] = useState(false);
  const [room, setRoom] = useState('');
  const [messages, setMessages] = useState([]);
  const [isConnected, setIsConnected] = useState(false);
  const [availableRooms, setAvailableRooms] = useState([]);

  const socket = useRef(null);
  const messagesEndRef = useRef(null);

  // --- EFFECTS ---

  // Auto-scroll to bottom
  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
  }, [messages]);

  // Fetch rooms when entering lobby
  useEffect(() => {
    if (view === 'lobby') {
      fetchRooms();
      // Optional: Poll for rooms every 5 seconds
      const interval = setInterval(fetchRooms, 5000);
      return () => clearInterval(interval);
    }
  }, [view]);

  // Connect to WebSocket when entering chat
  useEffect(() => {
    if (view !== 'chat' || !room) return;

    const wsHost = window.location.hostname;
    socket.current = new WebSocket(`ws://${wsHost}:8000/ws/${room}?username=${encodeURIComponent(username)}`);


    socket.current.onopen = () => {
      console.log("Connected");
      setIsConnected(true);
    };

    socket.current.onmessage = (event) => {
      const message = JSON.parse(event.data);
      setMessages((prev) => [...prev, message]);
    };

    socket.current.onclose = () => {
      console.log("Disconnected");
      setIsConnected(false);
    };

    return () => {
      if (socket.current) socket.current.close();
    };
  }, [view, room]);

  // --- HELPERS ---

  const fetchRooms = async () => {
    try {
      const host = window.location.hostname;
      const res = await fetch(`http://${host}:8000/rooms`);
      const data = await res.json();
      setAvailableRooms(data || []); // Fallback to empty array if null
    } catch (err) {
      console.error("Failed to fetch rooms:", err);
      setAvailableRooms([]); // Set to empty array on error
    }
  };

  const deleteRoom = async (roomName) => {
    if (!window.confirm(`Are you sure you want to delete "${roomName}"?`)) return;

    try {
      const host = window.location.hostname;
      const res = await fetch(`http://${host}:8000/rooms?name=${encodeURIComponent(roomName)}&username=${encodeURIComponent(username)}`, {
        method: 'DELETE',
      });

      if (res.ok) {
        alert('Room deleted successfully!');
        fetchRooms(); // Refresh room list
      } else {
        const error = await res.text();
        alert(`Failed to delete room: ${error}`);
      }
    } catch (err) {
      console.error('Delete room error:', err);
      alert('Failed to delete room');
    }
  };

  const getInitials = (name) => (name ? name.substring(0, 2).toUpperCase() : "?");

  const stringToColor = (str) => {
    let hash = 0;
    for (let i = 0; i < str.length; i++) {
      hash = str.charCodeAt(i) + ((hash << 5) - hash);
    }
    const c = (hash & 0x00FFFFFF).toString(16).toUpperCase();
    return '#' + "00000".substring(0, 6 - c.length) + c;
  };

  // --- HANDLERS ---

  // --- HANDLERS ---

  const handleAuth = async (e) => {
    e.preventDefault();
    const endpoint = isRegistering ? 'register' : 'login';
    const host = window.location.hostname;

    try {
      const res = await fetch(`http://${host}:8000/${endpoint}`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ username, password_hash: passwordHash })
      });

      const data = await res.json();

      if (!res.ok) {
        alert(data.message || "Authentication failed");
        return;
      }

      if (isRegistering) {
        alert("Registration successful! Please login.");
        setIsRegistering(false);
      } else {
        setView('lobby');
      }
    } catch (err) {
      console.error("Auth error:", err);
      alert("Failed to connect to server");
    }
  };

  const handleJoinRoom = (roomName) => {
    setRoom(roomName);
    setMessages([]); // Clear old messages
    setView('chat');
  };

  const handleCreateRoom = (e) => {
    e.preventDefault();
    const newRoom = e.target.elements.roomName.value.trim();
    if (newRoom) {
      handleJoinRoom(newRoom);
    }
  };

  const handleLeaveRoom = () => {
    if (socket.current) socket.current.close();
    setRoom('');
    setView('lobby');
  };

  const handleSendMessage = (text) => {
    if (socket.current && socket.current.readyState === WebSocket.OPEN) {
      const message = { roomID: room, username, content: text };
      socket.current.send(JSON.stringify(message));
    }
  };

  // --- VIEWS ---

  if (view === 'login') {
    return (
      <div className="App">
        <div className="join-room-container">
          <div className="login-logo">💬</div>
          <h1 className="login-title">Chat Connect</h1>
          <p className="login-subtitle">
            {isRegistering ? "Create an account to get started" : "Enter your credentials to join"}
          </p>

          <form className="join-room-form" onSubmit={handleAuth}>
            <input
              type="text"
              placeholder="Username"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              required
            />
            <input
              type="password"
              placeholder="Password"
              value={passwordHash}
              onChange={(e) => setPasswordHash(e.target.value)}
              required
            />
            <button type="submit">{isRegistering ? "Sign Up" : "Login"}</button>
          </form>

          <button className="toggle-auth" onClick={() => setIsRegistering(!isRegistering)}>
            {isRegistering ? "Already have an account? Login" : "New here? Create Account"}
          </button>
        </div>
      </div>
    );
  }

  if (view === 'lobby') {
    return (
      <div className="App">
        <header className="App-header">
          <h1>Lobby</h1>
          <div className="user-info">Logged in as: <strong>{username}</strong></div>
        </header>
        <div className="lobby-container">
          <div className="create-room-section">
            <h3>Create New Room</h3>
            <form className="join-room-form" onSubmit={handleCreateRoom}>
              <input name="roomName" type="text" placeholder="Room Name" />
              <button type="submit">Create & Join</button>
            </form>
          </div>

          <div className="active-rooms-section">
            <h3>Active Rooms</h3>
            <div className="room-list">
              {availableRooms.length === 0 ? (
                <p className="empty-rooms">No active rooms</p>
              ) : (
                <ul>
                  {availableRooms.map((r) => (
                    <li key={r.name || r} className="room-item">
                      <div className="room-info">
                        <span className="room-name">{r.name || r}</span>
                        {r.creator && <span className="room-creator">by {r.creator}</span>}
                      </div>
                      <div className="room-actions">
                        <button onClick={() => handleJoinRoom(r.name || r)} className="join-btn">
                          Join
                        </button>
                        {r.creator === username && (
                          <button
                            onClick={() => deleteRoom(r.name || r)}
                            className="delete-btn"
                            title="Delete this room"
                          >
                            🗑️
                          </button>
                        )}
                      </div>
                    </li>
                  ))}
                </ul>
              )}
            </div>
          </div>
        </div>
      </div>
    );
  }

  // Chat View
  return (
    <div className="App">
      <header className="App-header chat-header">
        <button className="back-button" onClick={handleLeaveRoom}>← Lobby</button>
        <h1>
          {room}
          <span
            className={`status-indicator ${isConnected ? 'online' : 'offline'}`}
            title={isConnected ? "Connected" : "Disconnected"}
          />
        </h1>
      </header>
      <div className="message-list">
        {messages.length === 0 && (
          <div className="empty-state">
            <p>No messages yet. Start the conversation!</p>
          </div>
        )}

        {messages.map((msg, index) => {
          const isMe = msg.username === username;
          const avatarColor = stringToColor(msg.username);

          return (
            <div key={index} className={`message-row ${isMe ? 'is-me' : ''}`}>
              {!isMe && (
                <div className="avatar" style={{ backgroundColor: avatarColor }}>
                  {getInitials(msg.username)}
                </div>
              )}

              <div className={`message ${isMe ? 'is-me' : 'is-other'}`}>
                {!isMe && (
                  <span className="message-username">{msg.username}</span>
                )}
                <span className="message-content">{msg.content}</span>
                <span className="message-time">
                  {new Date(msg.timestamp).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}
                </span>
              </div>
            </div>
          );
        })}
        <div ref={messagesEndRef} />
      </div>

      <InputForm onSendMessage={handleSendMessage} />
    </div>
  );
}

export default App;
