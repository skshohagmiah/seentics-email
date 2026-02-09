# Security Policy

## Supported Versions

We release patches for security vulnerabilities. Currently supported versions:

| Version | Supported          |
| ------- | ------------------ |
| 1.0.x   | :white_check_mark: |
| < 1.0   | :x:                |

## Reporting a Vulnerability

We take the security of Seentics Email seriously. If you believe you have found a security vulnerability, please report it to us as described below.

**Please do not report security vulnerabilities through public GitHub issues.**

Instead, please report them via email to: **security@seentics.com** (or create a private security advisory on GitHub)

You should receive a response within 48 hours. If for some reason you do not, please follow up via email to ensure we received your original message.

Please include the following information:

* Type of issue (e.g. buffer overflow, SQL injection, cross-site scripting, etc.)
* Full paths of source file(s) related to the manifestation of the issue
* The location of the affected source code (tag/branch/commit or direct URL)
* Any special configuration required to reproduce the issue
* Step-by-step instructions to reproduce the issue
* Proof-of-concept or exploit code (if possible)
* Impact of the issue, including how an attacker might exploit it

## Preferred Languages

We prefer all communications to be in English.

## Security Best Practices

When deploying Seentics Email:

1. **Change all default passwords** - Update PostgreSQL, Redis, MySQL, and RabbitMQ passwords
2. **Use strong JWT secrets** - Generate cryptographically secure random strings
3. **Enable HTTPS** - Use SSL/TLS certificates for all services
4. **Restrict network access** - Use firewalls to limit access to internal services
5. **Keep dependencies updated** - Regularly update Docker images and npm/Go packages
6. **Monitor logs** - Set up logging and monitoring for suspicious activity
7. **Regular backups** - Backup databases and configuration regularly
8. **API key rotation** - Encourage users to rotate API keys periodically

## Disclosure Policy

When we receive a security bug report, we will:

1. Confirm the problem and determine affected versions
2. Audit code to find any similar problems
3. Prepare fixes for all supported versions
4. Release new versions as soon as possible

## Comments on this Policy

If you have suggestions on how this process could be improved, please submit a pull request.
