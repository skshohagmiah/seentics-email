# Seentics Email - Postal Setup Guide

This guide will help you set up and configure the Postal mail server that powers the email delivery infrastructure.

## What is Postal?

Postal is a fully-featured open-source mail delivery platform for incoming & outgoing email. It handles:
- SMTP email sending
- Email tracking and analytics
- Bounce and complaint handling
- Webhook notifications for email events
- Domain and IP management

## Architecture

In our stack, Postal runs as a Docker container alongside:
- **Postal MySQL**: Database for Postal's data
- **Postal RabbitMQ**: Message queue for email processing
- **Backend API**: Our Go service that interfaces with Postal

## Initial Setup

### 1. Start the Services

```bash
docker-compose up -d
```

This starts all services including Postal, MySQL, and RabbitMQ.

### 2. Initialize Postal

Enter the Postal container:

```bash
docker exec -it seentics-email-postal bash
```

Initialize the Postal database:

```bash
postal initialize
```

### 3. Create Admin User

Create an admin user for the Postal web interface:

```bash
postal make-user
```

Follow the prompts to create your admin account. Save the credentials!

### 4. Create Organization

Organizations in Postal group mail servers together:

```bash
postal make-organization
```

Provide:
- Organization name (e.g., "Seentics")
- Permalink (e.g., "seentics")

### 5. Create Mail Server

Create your first mail server:

```bash
postal make-server
```

Provide:
- Server name (e.g., "Production")
- Permalink (e.g., "production")
- Organization (select the one you just created)

### 6. Generate API Key

After creating the mail server, generate an API key:

```bash
postal api-key create
```

**IMPORTANT**: Save this API key! Add it to your `.env` file:

```env
POSTAL_API_KEY=your-api-key-here
```

## Accessing Postal Web Interface

Postal provides a web interface for managing your mail servers:

1. Access at: `http://localhost:5000`
2. Login with the admin credentials you created
3. Navigate to your organization and mail server

## Configuring Domains

### 1. Add Domain in Postal

Via the web interface:
1. Go to your mail server
2. Click "Domains" → "Add Domain"
3. Enter your domain name (e.g., `yourdomain.com`)

### 2. Configure DNS Records

Postal will provide DNS records you need to add to your domain:

#### SPF Record
```
Type: TXT
Name: @
Value: v=spf1 include:postal.yourdomain.com ~all
```

#### DKIM Record
```
Type: TXT
Name: postal._domainkey
Value: [Provided by Postal]
```

#### DMARC Record
```
Type: TXT
Name: _dmarc
Value: v=DMARC1; p=none; rua=mailto:dmarc@yourdomain.com
```

#### MX Record (for receiving emails)
```
Type: MX
Name: @
Value: postal.yourdomain.com
Priority: 10
```

### 3. Verify Domain

After adding DNS records:
1. Wait for DNS propagation (can take up to 48 hours)
2. In Postal web interface, click "Verify" on your domain
3. Postal will check the DNS records

## Sending Your First Email

### Via API (using our backend)

```bash
curl -X POST http://localhost:8080/api/send \
  -H "X-API-Key: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "to": ["recipient@example.com"],
    "from": "sender@yourdomain.com",
    "subject": "Test Email",
    "html_body": "<h1>Hello!</h1><p>This is a test email.</p>"
  }'
```

### Via Postal SMTP

You can also use Postal's SMTP server directly:

- **Host**: `localhost` (or your Postal server address)
- **Port**: `25`
- **Username**: Your Postal SMTP username
- **Password**: Your Postal SMTP password

## Webhook Configuration

Postal can send webhooks for email events (delivered, bounced, opened, etc.):

1. In Postal web interface, go to your mail server
2. Click "Webhooks" → "Add Webhook"
3. Enter webhook URL: `http://backend:8080/webhooks/postal`
4. Select events to track
5. Save the webhook

Our backend automatically processes these webhooks and updates email statuses.

## Monitoring

### Postal Web Interface

- **Dashboard**: Overview of email activity
- **Message Queue**: See pending emails
- **Logs**: Detailed email logs
- **Statistics**: Delivery rates, bounces, etc.

### Backend API

Our backend provides additional analytics:
- Email delivery statistics
- Per-user email tracking
- API usage metrics

## Troubleshooting

### Emails Not Sending

1. **Check Postal logs**:
   ```bash
   docker logs seentics-email-postal
   ```

2. **Verify domain DNS records**:
   - Use online DNS checkers
   - Ensure all records are properly configured

3. **Check Postal queue**:
   - Access web interface
   - Look for stuck messages in queue

### SMTP Connection Issues

1. **Verify Postal is running**:
   ```bash
   docker ps | grep postal
   ```

2. **Check port accessibility**:
   ```bash
   telnet localhost 25
   ```

### API Authentication Errors

1. **Verify API key** in `.env` file
2. **Check API key permissions** in Postal web interface
3. **Ensure API key is for the correct mail server**

## Production Considerations

### Security

1. **Change default passwords** for MySQL and RabbitMQ
2. **Use strong API keys**
3. **Enable SSL/TLS** for SMTP connections
4. **Restrict network access** to Postal ports

### Performance

1. **Monitor resource usage**:
   - MySQL database size
   - RabbitMQ queue length
   - Postal container CPU/memory

2. **Scale horizontally** if needed:
   - Multiple Postal instances
   - Load balancer for SMTP

### Backup

1. **Backup Postal MySQL database** regularly
2. **Backup Postal configuration** files
3. **Document your DNS records**

## Additional Resources

- [Postal Documentation](https://docs.postalserver.io)
- [Postal GitHub](https://github.com/postalserver/postal)
- [Postal Discord Community](https://discord.postalserver.io)

## Support

For issues specific to our implementation:
1. Check backend logs: `docker logs seentics-email-backend`
2. Check Postal logs: `docker logs seentics-email-postal`
3. Review this documentation
4. Open an issue on GitHub
