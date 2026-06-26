# SpotSync - Smart Parking & EV Charging Reservation API

## Live Links

- GitHub Repository: https://github.com/CodedGrimoire/spotsync-api
- Live Deployment: https://spotsync-api-7vn1.onrender.com/
- Interview Video: https://drive.google.com/drive/folders/1F8RXZiVEiiGrZRTr3qnE9o48nfN-81VL?usp=sharing

## Overview

SpotSync is a clean architecture Go backend API for managing parking zones and EV charging reservations in busy locations such as airports and malls. It supports JWT authentication, role-based authorization, dynamic parking availability calculation, and concurrency-safe reservation creation.

## Features

- User registration and login
- JWT authentication
- Role-based access: driver and admin
- Admin parking zone management
- Public parking zone browsing
- Dynamic available spot calculation
- Authenticated reservation creation
- My reservations view
- Reservation cancellation
- Admin view of all reservations
- Concurrency-safe reservation creation using transaction and row-level locking

## Tech Stack

- Go
- Echo
- GORM
- PostgreSQL / NeonDB
- JWT
- bcrypt
- go-playground/validator
- godotenv

## Clean Architecture

SpotSync separates HTTP, business logic, and database access into clear layers. Requests flow through the application like this:

```text
Client
  -> Echo Router
  -> Middleware
  -> Handler
  -> Service
  -> Repository
  -> GORM Models
  -> PostgreSQL / NeonDB
```

Responses flow back in the reverse direction. Handlers never call GORM directly, repositories never return HTTP responses, and services contain the main business rules.

- `dto`: Request and response objects used at the API boundary.
- `handler`: HTTP request handling, request binding, validation, service calls, and JSON responses.
- `service`: Business logic such as auth rules, role checks, reservation rules, password hashing, JWT generation, and DTO mapping.
- `repository`: Database access only, including GORM queries, transactions, preloads, and row locks.
- `models`: GORM database models for users, parking zones, and reservations.
- `middleware`: JWT authentication and role-based route protection.
- `config`: Database connection and connection pool setup.
- `utils`: Shared helpers such as standard JSON responses and validation setup.
- `routes`: Route registration and dependency wiring.

## Environment Variables

```env
PORT=8080
DATABASE_URL=postgresql://username:password@host/database?sslmode=require
JWT_SECRET=replace_with_long_random_secret
```

Do not commit real `.env` values. Keep `DATABASE_URL` and `JWT_SECRET` secret.

## Local Setup

```bash
git clone <your-repository-url>
cd spotsync-api
go mod tidy
go run main.go
```

Health check:

```bash
curl http://localhost:8080/health
```

## API Endpoints

### Authentication

```text
POST /api/v1/auth/register
POST /api/v1/auth/login
```

### Parking Zones

```text
POST /api/v1/zones              admin only
GET /api/v1/zones               public
GET /api/v1/zones/:id           public
PUT /api/v1/zones/:id           admin only
DELETE /api/v1/zones/:id        admin only
```

### Reservations

```text
POST /api/v1/reservations                    authenticated
GET /api/v1/reservations/my-reservations     authenticated
DELETE /api/v1/reservations/:id              authenticated
GET /api/v1/reservations                     admin only
```

## Testing

Run the end-to-end API verification script:

```bash
chmod +x scripts/test_api.sh
./scripts/test_api.sh
```

For a deployed URL:

```bash
BASE_URL=https://your-deployed-url ./scripts/test_api.sh
```

Use [docs/FINAL_CHECKLIST.md](docs/FINAL_CHECKLIST.md) before submission to verify code quality, auth, parking zones, reservations, deployment, and submission requirements.

## Concurrency Safety

Reservation creation uses a GORM transaction. Inside the transaction, the selected `parking_zones` row is locked with `FOR UPDATE` before counting active reservations and creating a new reservation. This prevents two simultaneous requests from both reserving the final available spot.

## Deployment

### Render

1. Push code to GitHub.
2. Create a new Web Service on Render.
3. Connect the GitHub repo.
4. Build command:

```bash
go build -o app .
```

5. Start command:

```bash
./app
```

6. Add required environment variables:

```env
PORT=8080
DATABASE_URL=<your NeonDB connection string>
JWT_SECRET=<your long random secret>
```

7. Deploy.
8. Test:

```bash
curl https://your-app-name.onrender.com/health
```

### Railway

1. Create a new Railway project.
2. Deploy from the GitHub repo.
3. Add required environment variables:

```env
DATABASE_URL=<your NeonDB connection string>
JWT_SECRET=<your long random secret>
```

4. Railway may provide `PORT` automatically. The app reads `PORT` from the environment.
5. Test:

```bash
curl https://your-railway-domain/health
```

## Submission Checklist

- [ ] Public GitHub repo
- [ ] Live deployment URL
- [ ] Interview video URL
- [ ] At least 10 meaningful commits
- [ ] `.env` not committed
