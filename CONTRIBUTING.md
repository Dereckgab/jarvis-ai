# Contributing to JARVIS Full IA

Thank you for your interest in contributing to JARVIS Full IA! We welcome contributions from the community and are excited to work with you.

## Code of Conduct

Please be respectful and professional in all interactions. We are committed to providing a welcoming and inclusive environment for all contributors.

## Getting Started

1. **Fork the repository** on GitHub
2. **Clone your fork** locally
   ```bash
   git clone https://github.com/yourusername/jarvis-fullia.git
   cd jarvis-fullia
   ```
3. **Create a feature branch**
   ```bash
   git checkout -b feature/your-feature-name
   ```

## Development Workflow

### Backend Development

1. **Set up the development environment**
   ```bash
   cd backend
   go mod download
   ```

2. **Follow Go best practices**
   - Use `gofmt` for code formatting
   - Run `go vet` to check for common errors
   - Write tests for new functionality

3. **Commit your changes**
   ```bash
   git commit -m "feat: add new feature"
   ```

### Frontend Development

1. **Set up the development environment**
   ```bash
   cd frontend
   npm install
   npm run dev
   ```

2. **Follow React best practices**
   - Use functional components with hooks
   - Keep components small and focused
   - Write TypeScript for type safety

3. **Commit your changes**
   ```bash
   git commit -m "feat: add new feature"
   ```

## Commit Message Guidelines

We follow conventional commits:

- `feat:` A new feature
- `fix:` A bug fix
- `docs:` Documentation changes
- `style:` Code style changes (formatting, etc.)
- `refactor:` Code refactoring
- `test:` Adding or updating tests
- `chore:` Maintenance tasks

Example:
```
feat: add semantic memory search functionality
```

## Pull Request Process

1. **Update your branch** with the latest main
   ```bash
   git fetch origin
   git rebase origin/main
   ```

2. **Push your changes**
   ```bash
   git push origin feature/your-feature-name
   ```

3. **Open a Pull Request** on GitHub with:
   - Clear title describing the changes
   - Description of what was changed and why
   - Reference to any related issues (#123)
   - Screenshots for UI changes

4. **Address review feedback** promptly

5. **Ensure all checks pass** (tests, linting, etc.)

## Code Style

### Go
- Follow [Effective Go](https://golang.org/doc/effective_go)
- Use `gofmt` for formatting
- Use meaningful variable names
- Add comments for exported functions

### TypeScript/React
- Use ESLint configuration provided
- Use Prettier for code formatting
- Follow React hooks best practices
- Use meaningful component names

## Testing

### Backend
```bash
cd backend
go test ./... -v -cover
```

### Frontend
```bash
cd frontend
npm test
```

## Documentation

- Update README.md if adding new features
- Add comments to complex logic
- Update API documentation for new endpoints
- Include examples for new functionality

## Reporting Issues

When reporting bugs, please include:
- Clear description of the issue
- Steps to reproduce
- Expected behavior
- Actual behavior
- Environment details (OS, Go version, Node version, etc.)
- Screenshots or logs if applicable

## Feature Requests

For feature requests, please:
- Describe the feature clearly
- Explain the use case
- Provide examples if possible
- Consider backwards compatibility

## Questions?

Feel free to open an issue or discussion if you have questions about contributing.

---

Thank you for contributing to JARVIS Full IA! 🚀
