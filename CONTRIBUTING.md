# Contributing to Seentics Email

First off, thank you for considering contributing to Seentics Email! It's people like you that make this project such a great tool.

## Code of Conduct

This project and everyone participating in it is governed by our Code of Conduct. By participating, you are expected to uphold this code. Please report unacceptable behavior to the project maintainers.

## How Can I Contribute?

### Reporting Bugs

Before creating bug reports, please check the existing issues as you might find out that you don't need to create one. When you are creating a bug report, please include as many details as possible:

* **Use a clear and descriptive title**
* **Describe the exact steps to reproduce the problem**
* **Provide specific examples to demonstrate the steps**
* **Describe the behavior you observed and what you expected**
* **Include screenshots if relevant**
* **Include your environment details** (OS, Docker version, etc.)

### Suggesting Enhancements

Enhancement suggestions are tracked as GitHub issues. When creating an enhancement suggestion, please include:

* **Use a clear and descriptive title**
* **Provide a detailed description of the suggested enhancement**
* **Explain why this enhancement would be useful**
* **List any similar features in other projects**

### Pull Requests

* Fill in the required template
* Follow the coding style used throughout the project
* Include tests when adding new features
* Update documentation as needed
* End all files with a newline

## Development Setup

### Prerequisites

* Docker and Docker Compose
* Go 1.22+ (for backend development)
* Node.js 20+ (for frontend development)

### Setting Up Your Development Environment

1. **Fork the repository**

2. **Clone your fork**
   ```bash
   git clone https://github.com/YOUR_USERNAME/seentics-email.git
   cd seentics-email
   ```

3. **Start the development environment**
   ```bash
   docker-compose up -d
   ```

4. **Backend development**
   ```bash
   cd backend
   go mod download
   go run cmd/server/main.go
   ```

5. **Frontend development**
   ```bash
   cd frontend
   npm install
   npm run dev
   ```

## Coding Standards

### Backend (Go)

* Follow [Effective Go](https://golang.org/doc/effective_go.html) guidelines
* Use `gofmt` to format your code
* Run `go vet` before committing
* Write meaningful commit messages
* Add comments for exported functions and types

### Frontend (TypeScript/React)

* Use TypeScript for type safety
* Follow React best practices
* Use functional components with hooks
* Keep components small and focused
* Write meaningful component and variable names

### Git Commit Messages

* Use the present tense ("Add feature" not "Added feature")
* Use the imperative mood ("Move cursor to..." not "Moves cursor to...")
* Limit the first line to 72 characters or less
* Reference issues and pull requests liberally after the first line

Example:
```
Add email template support

- Implement template CRUD operations
- Add template variables parsing
- Update API documentation

Closes #123
```

## Project Structure

```
seentics-email/
├── backend/           # Go backend service
│   ├── cmd/          # Application entrypoints
│   ├── internal/     # Private application code
│   └── pkg/          # Public libraries
├── frontend/         # Next.js frontend
│   ├── app/          # Next.js app directory
│   ├── components/   # React components
│   └── lib/          # Utilities and helpers
├── docs/             # Documentation
└── docker-compose.yml
```

## Testing

### Backend Tests
```bash
cd backend
go test ./...
```

### Frontend Tests
```bash
cd frontend
npm test
```

## Documentation

* Update the README.md if you change functionality
* Add JSDoc/GoDoc comments for new functions
* Update API documentation in `docs/` if you modify endpoints

## Questions?

Feel free to open an issue with your question or reach out to the maintainers.

## License

By contributing, you agree that your contributions will be licensed under the MIT License.
