# Go Authentication Service

## Features
- JWT-based Access/Refresh Tokens
- Token refresh & validation
- Email notifications for IP changes
- Docker support for local development

---

## Setup

### Prerequisites
- **Go 1.20+**, **Docker**, **PostgreSQL**

### Steps
1. Copy `.env`: `cp .env.example .env`
2. Update `.env` with your configuration
3. Run the app:
   - **Locally**: `go run cmd/main.go`
   - **With Docker**:
     ```bash
       docker-compose -f docker-compose.development.yml up -d
       docker-compose -f docker-compose.development.yml logs -f --tail 500
       docker-compose -f docker-compose.development.yml down
     ```

---

## API Endpoints

### `/auth/tokens`
- **POST**: Generates access/refresh tokens.
- **Headers**: `X-Forwarded-For`
- **Query**: `user_id=<uuid>`

### `/auth/refresh`
- **POST**: Refreshes tokens.
- **Headers**: `Authorization`, `X-Refresh-Token`

---

## Environment Variables
- `APP_PORT`, `APP_JWT_SECRET`, `DB_URL`
- SMTP settings: `SMTP_HOST`, `SMTP_PORT`, `SMTP_USER`, `SMTP_PASS`, `SMTP_FROM`

---

## Test & Development
- Run tests: `go test ./... -v`
- Hot-reload (Air): `air`