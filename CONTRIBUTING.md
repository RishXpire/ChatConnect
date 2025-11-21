# Contributing to ChatConnect

Thank you for considering contributing to ChatConnect! This document provides guidelines for contributing to the project.

## Code of Conduct

By participating in this project, you agree to maintain a respectful and inclusive environment for all contributors.

## How to Contribute

### Reporting Bugs

Before creating bug reports, please check existing issues to avoid duplicates. When creating a bug report, include:

- Clear and descriptive title
- Steps to reproduce the issue
- Expected behavior
- Actual behavior
- System information (OS, Go version, Node version)
- Screenshots if applicable

### Suggesting Enhancements

Enhancement suggestions are welcome! Please provide:

- Clear description of the proposed feature
- Explanation of why this enhancement would be useful
- Examples of how the feature would work
- Any implementation ideas you may have

### Pull Requests

1. **Fork the Repository**
   ```bash
   git clone https://github.com/yourusername/ChatConnect.git
   cd ChatConnect
   ```

2. **Create a Branch**
   ```bash
   git checkout -b feature/your-feature-name
   ```

3. **Make Your Changes**
   - Follow the existing code style
   - Write clear commit messages
   - Add tests for new features
   - Update documentation as needed

4. **Test Your Changes**
   ```bash
   # Backend
   cd go-chat-backend
   go test ./...
   
   # Frontend
   cd chat-frontend-react
   npm test
   ```

5. **Commit Your Changes**
   ```bash
   git add .
   git commit -m "Add feature: description of your changes"
   ```

6. **Push to Your Fork**
   ```bash
   git push origin feature/your-feature-name
   ```

7. **Open a Pull Request**
   - Provide a clear description of the changes
   - Reference any related issues
   - Wait for review and address feedback

## Development Setup

### Prerequisites
- Go 1.19 or higher
- Node.js 16 or higher
- MongoDB (local or Atlas)

### Backend Development
```bash
cd go-chat-backend
go mod download
go run main.go
```

### Frontend Development
```bash
cd chat-frontend-react
npm install
npm start
```

## Coding Standards

### Go (Backend)
- Follow [Effective Go](https://golang.org/doc/effective_go.html) guidelines
- Use `gofmt` for code formatting
- Add comments for exported functions and types
- Write descriptive variable names
- Handle errors explicitly

### JavaScript/React (Frontend)
- Use functional components with hooks
- Follow ESLint recommendations
- Use meaningful component and variable names
- Keep components small and focused
- Add PropTypes for component props

### Git Commit Messages
- Use present tense ("Add feature" not "Added feature")
- Use imperative mood ("Move cursor to..." not "Moves cursor to...")
- Limit first line to 72 characters
- Reference issues and pull requests when applicable

## Project Structure

Maintain the existing project structure:
```
ChatApp/
├── go-chat-backend/
│   ├── main.go
│   ├── models/
│   └── websocket/
└── chat-frontend-react/
    └── src/
```

## Testing

- Write unit tests for new features
- Ensure all tests pass before submitting PR
- Maintain or improve code coverage

## Documentation

- Update README.md for new features
- Add code comments for complex logic
- Update API documentation for endpoint changes
- Include examples for new functionality

## Questions?

Feel free to open an issue for:
- Clarification on how to implement a feature
- Discussion about project direction
- Help with development setup

## License

By contributing, you agree that your contributions will be licensed under the MIT License.
