# ApexContact — Technical Reference

## Overview

ApexContact is a scalable live streaming platform for boxing and car racing events. Promoters
schedule events; viewers pay to watch via a paywall. The platform enforces strict anti-piracy
controls so that only paid, verified users can stream on a single device.

---

## Tech Stack

| Layer           | Technology                                  |
|-----------------|---------------------------------------------|
| Backend API     | Go (Fiber v2) — concurrency + performance   |
| Frontend        | Vue 3 + TypeScript + Vite                   |
| Styling         | Tailwind CSS (dark theme)                   |
| State Mgmt      | Pinia                                       |
| Database        | PostgreSQL 16                               |
| Cache / Sessions| Redis 7                                     |
| Streaming       | Nginx-RTMP (ingest) → HLS (delivery)        |
| Payments        | IntaSend (M-Pesa STK Push)                  |
| Reverse Proxy   | Nginx                                       |
| Deployment      | Docker Compose                              |

---

## User Roles

1. **Viewer** — registers, pays per event, watches stream on a single locked device
2. **Promoter** — creates/manages events, sets pricing, views analytics and revenue split
3. **Admin** — platform-wide control: fraud monitoring, user management, analytics
4. **Broadcaster** — receives RTMP stream key, pushes live feed to media server

---

## Design System

### Colour Palette (Dark / Sports theme)
```
Background:   #0A0A0F  (near-black)
Surface:      #141418  (card backgrounds)
Surface-2:    #1E1E26  (elevated surfaces)
Accent-Red:   #E8002D  (primary CTA, boxing)
Accent-Orange:#FF6B00  (secondary, racing)
Text-Primary: #FFFFFF
Text-Muted:   #A0A0B0
Success:      #00C853
Warning:      #FFB300
Error:        #FF1744
```

### UI Direction
- Netflix-style event grid on homepage (large hero banner + card rows)
- YouTube-style video player page (fullscreen HLS player, event info below)
- Countdown timers on upcoming events
- Smooth dark-mode-first design throughout

---

## Project Structure

```
ApexContact/
├── backend/                    # Go API server
│   ├── cmd/api/main.go
│   ├── internal/
│   │   ├── config/             # App config (env vars)
│   │   ├── api/
│   │   │   ├── router.go
│   │   │   ├── middleware/     # auth, anti-piracy, rate-limit, CORS
│   │   │   └── handlers/       # auth, events, stream, payments, users, admin
│   │   ├── domain/             # Structs: User, Event, Subscription, StreamToken, Payment
│   │   ├── repository/         # DB queries (pgx)
│   │   ├── service/            # Business logic
│   │   └── integrations/
│   │       ├── intasend/       # M-Pesa STK Push client
│   │       └── media/          # Media server API client
│   ├── migrations/             # SQL migrations (numbered)
│   ├── Dockerfile
│   ├── go.mod
│   └── .env.example
│
├── frontend/                   # Vue 3 SPA
│   ├── src/
│   │   ├── api/                # Axios API clients
│   │   ├── assets/styles/      # Tailwind + CSS variables
│   │   ├── components/
│   │   │   ├── ui/             # Button, Badge, Modal, Spinner
│   │   │   ├── layout/         # AppHeader, AppFooter, Sidebar
│   │   │   ├── events/         # EventCard, EventGrid, EventHero, CountdownTimer
│   │   │   ├── player/         # VideoPlayer, PlayerOverlay
│   │   │   └── payment/        # MpesaModal, PaymentStatus
│   │   ├── layouts/            # DefaultLayout, AuthLayout, DashboardLayout
│   │   ├── pages/
│   │   │   ├── Home.vue
│   │   │   ├── EventDetail.vue
│   │   │   ├── Watch.vue       # Protected HLS player page
│   │   │   ├── auth/           # Login, Register
│   │   │   ├── dashboard/      # Promoter: CreateEvent, Analytics
│   │   │   └── admin/          # AdminDashboard, FraudMonitor, UserManagement
│   │   ├── stores/             # Pinia: auth, events, player, payment
│   │   ├── router/             # Vue Router with route guards
│   │   ├── composables/        # useStream, usePayment, useFingerprint
│   │   └── types/              # TypeScript interfaces
│   ├── Dockerfile
│   ├── package.json
│   ├── vite.config.ts
│   ├── tailwind.config.js
│   └── tsconfig.json
│
├── media-server/               # Nginx-RTMP for stream ingest → HLS
│   ├── nginx.conf
│   └── Dockerfile
│
├── nginx/                      # Reverse proxy / SSL termination
│   └── nginx.conf
│
├── docker-compose.yml          # Local development
├── docker-compose.prod.yml     # Production overrides
├── .env.example
└── CLAUDE.md
```

---

## Docker Compose Services

| Service        | Image / Build       | Ports          | Purpose                          |
|----------------|---------------------|----------------|----------------------------------|
| `postgres`     | postgres:16-alpine  | 5432           | Primary database                 |
| `redis`        | redis:7-alpine      | 6379           | Sessions, tokens, rate-limiting  |
| `media-server` | ./media-server      | 1935, 8888     | RTMP ingest, HLS segment output  |
| `backend`      | ./backend           | 8000           | Go REST API                      |
| `frontend`     | ./frontend          | 3000           | Vue 3 SPA (served by Nginx)      |
| `nginx`        | nginx:alpine        | 80, 443        | Reverse proxy for all services   |

---

## API Structure (Backend)

### Auth
- `POST /api/v1/auth/register`
- `POST /api/v1/auth/login`
- `POST /api/v1/auth/refresh`
- `POST /api/v1/auth/logout`

### Events
- `GET  /api/v1/events`             — list upcoming/live events
- `GET  /api/v1/events/:id`         — event detail
- `POST /api/v1/events`             — create event (Promoter)
- `PUT  /api/v1/events/:id`         — update event (Promoter)

### Streaming
- `POST /api/v1/stream/:eventId/subscribe` — pay + get stream token
- `GET  /api/v1/stream/:eventId/token`     — validate and fetch HLS URL
- `POST /api/v1/stream/ingest/callback`    — media server webhook

### Payments
- `POST /api/v1/payments/initiate`  — trigger M-Pesa STK Push
- `POST /api/v1/payments/callback`  — IntaSend webhook
- `GET  /api/v1/payments/status/:id`

### Promoter Dashboard
- `GET  /api/v1/promoter/events`
- `GET  /api/v1/promoter/analytics/:eventId`
- `GET  /api/v1/promoter/revenue`

### Admin
- `GET  /api/v1/admin/users`
- `GET  /api/v1/admin/fraud`
- `POST /api/v1/admin/users/:id/lock`
- `GET  /api/v1/admin/analytics`

---

## Anti-Piracy Strategy

1. **Unique stream token** — generated per subscription (UUID + event + user + device hash)
2. **Device fingerprint** — collected in browser (canvas, fonts, screen), stored on token
3. **IP lock** — first stream start locks token to that IP; new IP triggers re-verification
4. **Single-session enforcement** — Redis tracks active session per token; second open invalidates both
5. **VPN/Proxy detection** — IP reputation API check on subscribe + stream start
6. **Token expiry** — tokens expire when event ends (no offline replay)
7. **Signed HLS URLs** — media server checks token signature on each segment request

---

## Streaming Flow

```
Broadcaster  →  RTMP push  →  Nginx-RTMP (media-server:1935)
                                        ↓
                              HLS segments written to /tmp/hls
                                        ↓
                     Backend issues signed HLS URL to viewer
                                        ↓
              Viewer browser  →  hls.js player  →  fetches .m3u8 + .ts segments
```

---

## Payment Flow (M-Pesa STK Push)

```
1. Viewer clicks "Buy Ticket"
2. Frontend opens MpesaModal → collects phone number
3. POST /api/v1/payments/initiate → backend calls IntaSend API
4. IntaSend sends STK Push to user's phone
5. User enters M-Pesa PIN
6. IntaSend calls POST /api/v1/payments/callback (webhook)
7. Backend confirms payment → creates Subscription + StreamToken
8. Frontend polls GET /api/v1/payments/status/:id
9. On success → redirect to Watch page with token
```

---

## Database Schema (Key Tables)

- `users` — id, email, password_hash, role, phone, created_at
- `events` — id, promoter_id, title, description, sport_type, scheduled_at, status, price, thumbnail_url, stream_key
- `subscriptions` — id, user_id, event_id, payment_id, stream_token, device_fingerprint, ip_lock, active_session_id, created_at
- `payments` — id, user_id, event_id, amount, currency, status, intasend_ref, created_at
- `stream_analytics` — id, event_id, viewer_count, peak_viewers, total_revenue, created_at
- `fraud_flags` — id, user_id, subscription_id, reason, detected_at, resolved

---

## Development Conventions

- Backend: handlers call services; services call repositories. No DB calls in handlers.
- Frontend: pages use stores + composables. No raw API calls in components.
- All API responses follow `{ data, error, meta }` envelope.
- Migrations are numbered SQL files: `001_init.sql`, `002_add_analytics.sql`, etc.
- Environment variables defined in `.env.example`; never commit real secrets.

---

## currentDate
Today's date is 2026-02-18.
