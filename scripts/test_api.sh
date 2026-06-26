#!/usr/bin/env bash

set -e

BASE_URL="${BASE_URL:-http://localhost:8080}"

if ! command -v jq >/dev/null 2>&1; then
  echo "jq is required for this script. Install jq first."
  exit 1
fi

print_section() {
  echo
  echo "== $1 =="
}

pretty_print() {
  echo "$1" | jq .
}

post_json() {
  local url="$1"
  local body="$2"
  shift 2

  curl -sS -X POST "$url" \
    -H "Content-Type: application/json" \
    "$@" \
    -d "$body"
}

print_section "A. Health check"
HEALTH_RESPONSE=$(curl -sS "$BASE_URL/health")
pretty_print "$HEALTH_RESPONSE"

print_section "B. Register admin"
ADMIN_REGISTER_RESPONSE=$(post_json "$BASE_URL/api/v1/auth/register" '{
  "name": "Admin User",
  "email": "admin@spotsync.com",
  "password": "admin123",
  "role": "admin"
}')
pretty_print "$ADMIN_REGISTER_RESPONSE"

print_section "C. Login admin"
ADMIN_LOGIN_RESPONSE=$(post_json "$BASE_URL/api/v1/auth/login" '{
  "email": "admin@spotsync.com",
  "password": "admin123"
}')
pretty_print "$ADMIN_LOGIN_RESPONSE"
ADMIN_TOKEN=$(echo "$ADMIN_LOGIN_RESPONSE" | jq -r '.data.token // empty')
if [ -z "$ADMIN_TOKEN" ]; then
  echo "Failed to extract admin token."
  exit 1
fi

print_section "D. Create parking zone as admin"
ZONE_CREATE_RESPONSE=$(post_json "$BASE_URL/api/v1/zones" '{
  "name": "Terminal 1 EV Charging",
  "type": "ev_charging",
  "total_capacity": 2,
  "price_per_hour": 5.50
}' -H "Authorization: Bearer $ADMIN_TOKEN")
pretty_print "$ZONE_CREATE_RESPONSE"
ZONE_ID=$(echo "$ZONE_CREATE_RESPONSE" | jq -r '.data.id // empty')
if [ -z "$ZONE_ID" ]; then
  echo "Failed to extract zone ID."
  exit 1
fi

print_section "E. Get all zones public"
ZONES_RESPONSE=$(curl -sS "$BASE_URL/api/v1/zones")
pretty_print "$ZONES_RESPONSE"

print_section "F. Register driver"
DRIVER_REGISTER_RESPONSE=$(post_json "$BASE_URL/api/v1/auth/register" '{
  "name": "Driver One",
  "email": "driver1@spotsync.com",
  "password": "driver123",
  "role": "driver"
}')
pretty_print "$DRIVER_REGISTER_RESPONSE"

print_section "G. Login driver"
DRIVER_LOGIN_RESPONSE=$(post_json "$BASE_URL/api/v1/auth/login" '{
  "email": "driver1@spotsync.com",
  "password": "driver123"
}')
pretty_print "$DRIVER_LOGIN_RESPONSE"
DRIVER_TOKEN=$(echo "$DRIVER_LOGIN_RESPONSE" | jq -r '.data.token // empty')
if [ -z "$DRIVER_TOKEN" ]; then
  echo "Failed to extract driver token."
  exit 1
fi

print_section "H. Create reservation as driver"
RESERVATION_CREATE_RESPONSE=$(post_json "$BASE_URL/api/v1/reservations" "{
  \"zone_id\": $ZONE_ID,
  \"license_plate\": \"ABC-1234\"
}" -H "Authorization: Bearer $DRIVER_TOKEN")
pretty_print "$RESERVATION_CREATE_RESPONSE"
RESERVATION_ID=$(echo "$RESERVATION_CREATE_RESPONSE" | jq -r '.data.id // empty')
if [ -z "$RESERVATION_ID" ]; then
  echo "Failed to extract reservation ID."
  exit 1
fi

print_section "I. Get my reservations"
MY_RESERVATIONS_RESPONSE=$(curl -sS "$BASE_URL/api/v1/reservations/my-reservations" \
  -H "Authorization: Bearer $DRIVER_TOKEN")
pretty_print "$MY_RESERVATIONS_RESPONSE"

print_section "J. Driver tries admin route, expected 403 Forbidden"
DRIVER_ADMIN_RESPONSE=$(curl -sS -w "\nHTTP_STATUS:%{http_code}\n" "$BASE_URL/api/v1/reservations" \
  -H "Authorization: Bearer $DRIVER_TOKEN")
echo "$DRIVER_ADMIN_RESPONSE" | sed '/^HTTP_STATUS:/d' | jq .
DRIVER_ADMIN_STATUS=$(echo "$DRIVER_ADMIN_RESPONSE" | awk -F: '/^HTTP_STATUS:/ {print $2}')
echo "HTTP status: $DRIVER_ADMIN_STATUS"
if [ "$DRIVER_ADMIN_STATUS" != "403" ]; then
  echo "Expected 403 Forbidden, got $DRIVER_ADMIN_STATUS."
  exit 1
fi

print_section "K. Admin gets all reservations"
ALL_RESERVATIONS_RESPONSE=$(curl -sS "$BASE_URL/api/v1/reservations" \
  -H "Authorization: Bearer $ADMIN_TOKEN")
pretty_print "$ALL_RESERVATIONS_RESPONSE"

print_section "L. Cancel reservation"
CANCEL_RESPONSE=$(curl -sS -X DELETE "$BASE_URL/api/v1/reservations/$RESERVATION_ID" \
  -H "Authorization: Bearer $DRIVER_TOKEN")
pretty_print "$CANCEL_RESPONSE"

print_section "M. Get single zone again"
ZONE_RESPONSE=$(curl -sS "$BASE_URL/api/v1/zones/$ZONE_ID")
pretty_print "$ZONE_RESPONSE"
AVAILABLE_SPOTS=$(echo "$ZONE_RESPONSE" | jq -r '.data.available_spots // empty')
echo "available_spots after cancellation: $AVAILABLE_SPOTS"

echo
echo "API flow completed successfully."
