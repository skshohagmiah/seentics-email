# Email Management Platform

A comprehensive email management platform similar to Resend/Mailgun, built with Next.js, Go, and Postal server.

## Architecture

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

## License

MIT

## Support

For issues and questions, please open an issue on GitHub.
# seentics-email
