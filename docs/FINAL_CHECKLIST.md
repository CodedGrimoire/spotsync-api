# Final Checklist

## Code Quality

- [ ] Clean architecture followed
- [ ] Handler does not call DB
- [ ] Service contains business logic
- [ ] Repository contains GORM only
- [ ] DTOs used for requests/responses

## Auth

- [ ] Register works
- [ ] Login works
- [ ] Password hashed with bcrypt
- [ ] JWT contains user_id and role
- [ ] Missing token gives 401
- [ ] Non-admin admin route gives 403

## Parking Zones

- [ ] Admin can create zone
- [ ] Public can view zones
- [ ] available_spots is calculated dynamically
- [ ] Admin can update/delete zones

## Reservations

- [ ] Authenticated user can reserve
- [ ] User can view own reservations
- [ ] User can cancel own reservation
- [ ] Admin can view all reservations
- [ ] Driver cannot view all reservations
- [ ] Full zone returns 409
- [ ] Transaction + FOR UPDATE row lock is used

## Deployment

- [ ] .env not committed
- [ ] .env.example exists
- [ ] render.yaml exists
- [ ] DATABASE_URL set in platform
- [ ] JWT_SECRET set in platform
- [ ] /health works on deployed URL

## Submission

- [ ] README completed
- [ ] Interview notes prepared
- [ ] Deployment instructions added
- [ ] Public GitHub repo link
- [ ] Live deployment link
- [ ] Interview video link
- [ ] At least 10 meaningful commits
