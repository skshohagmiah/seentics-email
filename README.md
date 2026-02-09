# Seentics Email 📮

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Docker](https://img.shields.io/badge/docker-%230db7ed.svg?style=flat&logo=docker&logoColor=white)](https://www.docker.com/)
[![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=flat&logo=go&logoColor=white)](https://golang.org/)
[![Next.js](https://img.shields.io/badge/Next.js-black?style=flat&logo=next.js&logoColor=white)](https://nextjs.org/)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](CONTRIBUTING.md)

> **Open-source email management platform** - A self-hosted alternative to Resend, Mailgun, and SendGrid

A comprehensive email management platform built with Next.js, Go, and Postal server. Send, track, and manage transactional emails with full control over your infrastructure.

## ✨ Features

- **Frontend**: Next.js 16 with TypeScript and Tailwind CSS
- **Backend**: Go with Gin framework
- **Email Server**: Postal (open-source mail delivery platform)
- **Database**: PostgreSQL for application data
- **Cache**: Redis for rate limiting
- **Message Queue**: RabbitMQ (for Postal)

## Features

- ✅ User authentication with JWT
- ✅ API key management with rate limiting
- ✅ Email sending via Postal
- ✅ Email tracking and analytics
- ✅ Domain management and verification
- ✅ Webhook configuration for email events
- ✅ Modern, responsive dashboard

## Quick Start

### Prerequisites

- Docker and Docker Compose
- Node.js 20+ (for local frontend development)
- Go 1.22+ (for local backend development)

### Running with Docker Compose

The entire stack (including Postal server) runs with a single command:

```bash
docker-compose up -d
```

This will start:
- PostgreSQL (port 5432)
- Redis (port 6379)
- Postal MySQL (internal)
- Postal RabbitMQ (internal)
- Postal Server (ports 5000 for API, 25 for SMTP)
- Backend API (port 8080)
- Frontend (port 3000)

Access the application at: **http://localhost:3000**

### Initial Postal Setup

After starting the services, you need to initialize Postal:

```bash
# Enter the Postal container
docker exec -it seentics-email-postal bash

# Initialize Postal
postal initialize

# Create an organization and mail server
postal make-user
postal make-organization
postal make-server
```

Follow the prompts to create your first organization and mail server. Save the API key generated - you'll need to add it to your environment.

### Environment Variables

Create a `.env` file in the root directory:

```env
# Postal Configuration
POSTAL_API_KEY=your-postal-api-key-here

# JWT Secret (change in production)
JWT_SECRET=your-super-secret-jwt-key-change-in-production
```

## Development

### Backend Development

```bash
cd backend
cp .env.example .env
# Edit .env with your configuration
go mod download
go run cmd/server/main.go
```

### Frontend Development

```bash
cd frontend
npm install
npm run dev
```

## API Documentation

### Authentication

- `POST /api/auth/signup` - Create new account
- `POST /api/auth/login` - Login
- `GET /api/profile` - Get user profile (requires JWT)

### API Keys

- `GET /api/keys` - List API keys
- `POST /api/keys` - Create API key
- `PUT /api/keys/:id` - Update API key
- `DELETE /api/keys/:id` - Delete API key

### Email Sending

- `POST /api/send` - Send email (requires API key in `X-API-Key` header)
- `GET /api/emails` - List sent emails
- `GET /api/emails/:id` - Get email details

### Domains

- `GET /api/domains` - List domains
- `POST /api/domains` - Add domain
- `GET /api/domains/:id/verify` - Get DNS verification records
- `DELETE /api/domains/:id` - Delete domain

### Webhooks

- `GET /api/webhooks` - List webhooks
- `POST /api/webhooks` - Create webhook
- `DELETE /api/webhooks/:id` - Delete webhook

## Sending Emails

### Using API Key

```bash
curl -X POST http://localhost:8080/api/send \
  -H "X-API-Key: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "to": ["recipient@example.com"],
    "from": "sender@yourdomain.com",
    "subject": "Hello from Seentics Email",
    "html_body": "<h1>Hello!</h1><p>This is a test email.</p>"
  }'
```

## Project Structure

```
.
├── backend/
│   ├── cmd/server/          # Main application
│   ├── internal/
│   │   ├── config/          # Configuration
│   │   ├── database/        # Database connection
│   │   ├── handlers/        # HTTP handlers
│   │   ├── middleware/      # Middleware (auth, rate limiting)
│   │   ├── models/          # Data models
│   │   └── postal/          # Postal API client
│   └── Dockerfile
├── frontend/
│   ├── src/
│   │   ├── app/             # Next.js pages
│   │   ├── components/      # React components
│   │   └── lib/             # Utilities (API client, auth)
│   └── Dockerfile
└── docker-compose.yml       # Full stack orchestration
```

## 🤝 Contributing

We love contributions! Seentics Email is an open-source project and we welcome contributions of all kinds:

- 🐛 Bug reports and fixes
- ✨ Feature requests and implementations
- 📖 Documentation improvements
- 🎨 UI/UX enhancements
- 🧪 Tests and quality improvements

Please read our [Contributing Guide](CONTRIBUTING.md) to get started.

### Development Setup

```bash
# Clone the repository
git clone https://github.com/yourusername/seentics-email.git
cd seentics-email

# Start development environment
docker-compose up -d

# Backend development
cd backend && go run cmd/server/main.go

# Frontend development
cd frontend && npm run dev
```

## 🌟 Community

- **GitHub Discussions**: Ask questions and share ideas
- **Issues**: Report bugs and request features
- **Discord**: Join our community server (coming soon)
- **Twitter**: Follow [@SeenticsEmail](https://twitter.com/seenticsemail) for updates

## 📝 Roadmap

- [ ] Email templates with variables
- [ ] Batch email sending
- [ ] Email scheduling
- [ ] Advanced analytics and reporting
- [ ] Multi-user organizations
- [ ] SMTP relay support
- [ ] Email verification service
- [ ] Terraform/Kubernetes deployment

See the [open issues](https://github.com/yourusername/seentics-email/issues) for a full list of proposed features.

## 🙏 Acknowledgments

- [Postal](https://github.com/postalserver/postal) - The amazing open-source mail delivery platform
- [Resend](https://resend.com) & [Mailgun](https://www.mailgun.com) - Inspiration for the API design
- All our [contributors](https://github.com/yourusername/seentics-email/graphs/contributors)

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 💖 Support

If you find this project useful, please consider:

- ⭐ Starring the repository
- 🐛 Reporting bugs
- 💡 Suggesting new features
- 🔀 Submitting pull requests
- 📢 Sharing with others

---

**Built with ❤️ by the open-source community**
# seentics-email
