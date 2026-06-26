# SpotSync - Smart Parking & EV Charging Reservation API

A clean architecture Go backend for managing parking zones and EV charging reservations with JWT authentication, role-based authorization, and concurrency-safe reservation creation.

## Tech Stack

- Go
- Echo
- GORM
- PostgreSQL / NeonDB
- JWT
- bcrypt
- Validator

## Architecture

- `dto`: Request and response objects used at the API boundary.
- `handler`: HTTP request handling, binding, validation, and response formatting.
- `service`: Business logic, authorization rules, password hashing, token generation, and DTO mapping.
- `repository`: Database access, GORM queries, transactions, preloads, and row locks.
- `models`: GORM database models.
- `middleware`: JWT authentication and role protection.
- `config`: Database connection setup.

## Environment Variables

```env
PORT=8080
DATABASE_URL=your_neon_postgresql_connection_string
JWT_SECRET=your_jwt_secret
```

## Local Setup

```bash
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

## Curl Testing Commands

Set the base URL:

```bash
BASE_URL=http://localhost:8080
```

1. Register admin

```bash
curl -X POST "$BASE_URL/api/v1/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Admin User",
    "email": "admin@spotsync.test",
    "password": "password123",
    "role": "admin"
  }'
```

2. Login admin

```bash
curl -X POST "$BASE_URL/api/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@spotsync.test",
    "password": "password123"
  }'
```

Save the token:

```bash
ADMIN_TOKEN=replace_with_admin_token
```

3. Create zone with admin token

```bash
curl -X POST "$BASE_URL/api/v1/zones" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -d '{
    "name": "Central Parking",
    "type": "ev_charging",
    "total_capacity": 10,
    "price_per_hour": 5.50
  }'
```

4. Get all zones

```bash
curl "$BASE_URL/api/v1/zones"
```

5. Register driver

```bash
curl -X POST "$BASE_URL/api/v1/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Driver User",
    "email": "driver@spotsync.test",
    "password": "password123",
    "role": "driver"
  }'
```

6. Login driver

```bash
curl -X POST "$BASE_URL/api/v1/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "driver@spotsync.test",
    "password": "password123"
  }'
```

Save the token:

```bash
DRIVER_TOKEN=replace_with_driver_token
```

7. Create reservation

```bash
curl -X POST "$BASE_URL/api/v1/reservations" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $DRIVER_TOKEN" \
  -d '{
    "zone_id": 1,
    "license_plate": "DHK-1234"
  }'
```

8. Get my reservations

```bash
curl "$BASE_URL/api/v1/reservations/my-reservations" \
  -H "Authorization: Bearer $DRIVER_TOKEN"
```

9. Cancel reservation

```bash
curl -X DELETE "$BASE_URL/api/v1/reservations/1" \
  -H "Authorization: Bearer $DRIVER_TOKEN"
```

10. Admin get all reservations

```bash
curl "$BASE_URL/api/v1/reservations" \
  -H "Authorization: Bearer $ADMIN_TOKEN"
```

## Concurrency Note

Reservation creation uses a GORM transaction and row-level lock with `FOR UPDATE` on the `parking_zones` row before checking active reservation count and creating the reservation. This prevents overbooking when multiple users try to reserve the last spot at the same time.

## Submission

GitHub Repo:

Live Deployment:

Interview Video:
