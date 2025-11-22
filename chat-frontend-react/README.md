# ChatConnect Frontend

React-based frontend for the ChatConnect real-time chat application.

## Features

- Modern glassmorphism UI design
- Real-time messaging with WebSocket
- User authentication interface
- Room creation and management
- Message history display
- Responsive layout

## Technology Stack

- React 18
- WebSocket API
- CSS3 (Custom styling)
- Fetch API for HTTP requests

## Getting Started

### Prerequisites

- Node.js 16 or higher
- npm 8 or higher

### Installation

```bash
# Install dependencies
npm install

# Start development server
npm start
```

Application opens at `http://localhost:3000`

### Building for Production

```bash
# Create production build
npm run build

# Serve build folder
npx serve -s build
```

## Project Structure

```
src/
├── App.js              # Main application component
├── App.css            # Application styles
├── components/
│   └── InputForm.js   # Message input component
├── index.js           # Entry point
└── index.css          # Global styles
```

## Components

### App
Main application component managing:
- View state (login, lobby, chat)
- WebSocket connections
- User authentication
- Room management

### InputForm
Reusable input component for:
- Message sending
- Form submissions

## Configuration

Update backend URL in `App.js`:

```javascript
const host = window.location.hostname;
// Connects to http://${host}:8000
```

For production, create `.env.production`:
```
REACT_APP_API_URL=https://your-api-domain.com
```

## Available Scripts

### `npm start`
Runs development server on `http://localhost:3000`

### `npm test`
Launches test runner

### `npm run build`
Builds app for production to `build/` folder

### `npm run eject`
Ejects from Create React App (one-way operation)

## Styling

Custom CSS with:
- Glassmorphism effects
- Gradient backgrounds
- Smooth animations
- Responsive design

## WebSocket Integration

Connects to backend WebSocket server:
```javascript
const socket = new WebSocket(`ws://${host}:8000/ws/${room}?username=${username}`);
```

## Deployment

### Vercel
```bash
npm install -g vercel
vercel
```

### Netlify
```bash
npm run build
# Deploy build/ folder
```

### GitHub Pages
```bash
npm install --save gh-pages

# Add to package.json:
"homepage": "https://username.github.io/repo-name",
"scripts": {
  "predeploy": "npm run build",
  "deploy": "gh-pages -d build"
}

npm run deploy
```

## Browser Support

- Chrome (latest)
- Firefox (latest)
- Safari (latest)
- Edge (latest)

## Environment Variables

Create `.env.local`:
```
REACT_APP_API_URL=http://localhost:8000
REACT_APP_WS_URL=ws://localhost:8000
```

## Contributing

See main project CONTRIBUTING.md for guidelines.

## License

MIT License - see LICENSE file in root directory.
