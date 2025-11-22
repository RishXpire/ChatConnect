const socket = new WebSocket("ws://localhost:8000/ws");

const form = document.getElementById('form');
const input = document.getElementById('input');
const messages = document.getElementById('messages');

// 1. We'll store the username here. Start it as null.
let username = null;

socket.onopen = (event) => {
    console.log("Successfully connected to WebSocket server");
    // DO NOT prompt here. Let the connection stay idle.
};

socket.onmessage = (event) => {
    // This will now run 29 times as soon as the page loads,
    // because the browser is not frozen.
    const message = JSON.parse(event.data);
    console.log("Message from server:", message);
    
    const displayTime = new Date(message.timestamp).toLocaleTimeString('en-US', {
        hour: '2-digit',
        minute: '2-digit'
    });
    const li = document.createElement('li');
    li.textContent = `[${displayTime}] ${message.username}: ${message.content}`;
    messages.appendChild(li);
};

socket.onclose = (event) => {
    console.log("Disconnected from WebSocket server", event);
    const li = document.createElement('li');
    li.textContent = "Disconnected from server.";
    messages.appendChild(li);
};

socket.onerror = (error) => {
    console.error("WebSocket Error:", error);
    const li = document.createElement('li');
    li.textContent = "WebSocket Error!";
    messages.appendChild(li);
};

// 2. This event listener is now outside 'onopen'
form.addEventListener('submit', (event) => {
    event.preventDefault();

    // 3. Check if we have a username. If not, ask for it.
    if (!username) {
        username = prompt("Please enter your username:");
        if (!username) { // If they click "Cancel"
            username = "Anonymous";
        }
    }

    // 4. Now, send the message.
    if (input.value) {
        const message = {
            username: username,
            content: input.value
        };
        socket.send(JSON.stringify(message));
        input.value = '';
    } else {
        console.log("Input is empty. Nothing was sent.");
    }
});
