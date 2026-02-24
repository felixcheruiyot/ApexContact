docker compose exec postgres psql -U apex -d apexcontact -c \
    "UPDATE users SET role='admin' WHERE email='felix@intasend.com';"



Streaming instruction
How it all works end-to-end                                   
                                                                                                          
  Broadcaster (OBS setup)                                 
                                                                                                          
  1. Promoter goes to Dashboard → Stream Key button next to their event                                   
  2. In OBS: Settings → Stream → Service: Custom                                                          
    - Server: rtmp://your-host:1935/live                                                                  
    - Stream Key: (from dashboard)                                                                        
  3. Click Start Streaming — that's it

  What happens automatically

  OBS push → Nginx-RTMP (port 1935)
           → writes HLS chunks to /tmp/hls/<stream_key>/index.m3u8
           → POST /api/v1/stream/ingest/callback  (form: call=publish, name=<stream_key>)
           → backend sets event status = 'live'
           → event appears as LIVE on homepage

  Viewer path

  Viewer → pays via M-Pesa → gets stream_token
         → GET /api/v1/stream/:eventId/token?token=<stream_token>
         → backend returns HLS URL: /hls/<stream_key>/index.m3u8
         → Watch.vue loads hls.js → plays stream

  Bugs fixed in this session

  ┌────────────────┬──────────────────────────────────────────────────────────────────────────────┐
  │      Bug       │                                What was wrong                                │
  ├────────────────┼──────────────────────────────────────────────────────────────────────────────┤
  │ IngestCallback │ Nginx-RTMP sends form data (call=publish), not JSON (action)                 │
  ├────────────────┼──────────────────────────────────────────────────────────────────────────────┤
  │ HLS URL        │ Used internal Docker hostname — browsers can't reach it                      │
  ├────────────────┼──────────────────────────────────────────────────────────────────────────────┤
  │ Stream key     │ No way for promoter to retrieve it — added GET /promoter/stream-key/:eventId │
  └────────────────┴──────────────────────────────────────────────────────────────────────────────┘


# How to apply migrations

  docker exec -i livestreamify-postgres-1 psql -U apex -d livestreamify < backend/migrations/004_event_publishing_flow.sql

  docker exec -i livestreamify-postgres-1 psql -U apex -d livestreamify < backend/migrations/003_user_profile.sql

  docker exec -i livestreamify-postgres-1 psql -U apex -d livestreamify < backend/migrations/006_merge_roles.sql

  docker exec -i livestreamify-postgres-1 psql -U apex -d livestreamify < backend/migrations/010_withdrawals.sql
  docker exec -i livestreamify-postgres-1 psql -U apex -d livestreamify < backend/migrations/011_public_events.sql


# References

- https://www.geniussports.com/
- https://sportradar.com/media-tech/data-content/


What was implemented                                                                                                    
                                                                                                                          
  Flow (standard OAuth 2.0 Authorization Code)                                                                            
                                                                                                                          
  Login.vue → Google consent screen
           ← /auth/callback?code=...
  OAuthCallback.vue → POST /api/v1/auth/google/callback { code, redirect_uri }
                    ← { token, user }  (our own JWT)
           → store in localStorage, redirect home

  Backend

  ┌──────────────────┬─────────────────────────────────────────────────────────────────────────────────────────────────┐
  │       File       │                                             Change                                              │
  ├──────────────────┼─────────────────────────────────────────────────────────────────────────────────────────────────┤
  │ config/config.go │ Added GoogleClientID, GoogleClientSecret fields + getEnv calls                                  │
  ├──────────────────┼─────────────────────────────────────────────────────────────────────────────────────────────────┤
  │ handlers/auth.go │ Added GoogleCallback handler + exchangeGoogleCode + fetchGoogleUserInfo helpers (stdlib only —  │
  │                  │ no new deps)                                                                                    │
  ├──────────────────┼─────────────────────────────────────────────────────────────────────────────────────────────────┤
  │ api/router.go    │ POST /api/v1/auth/google/callback                                                               │
  ├──────────────────┼─────────────────────────────────────────────────────────────────────────────────────────────────┤
  │ .env.example     │ Documented GOOGLE_CLIENT_ID + GOOGLE_CLIENT_SECRET                                              │
  └──────────────────┴─────────────────────────────────────────────────────────────────────────────────────────────────┘

  Upsert logic: First Google sign-in creates the account with password_hash = 'GOOGLE_OAUTH' (a non-bcrypt string —
  email/password login will gracefully reject it). Subsequent sign-ins update email_verified = true and leave everything
  else intact.

  Frontend

  ┌──────────────────────────────────┬─────────────────────────────────────────────────────────────────────┐
  │               File               │                               Change                                │
  ├──────────────────────────────────┼─────────────────────────────────────────────────────────────────────┤
  │ src/api/auth.ts                  │ googleCallback(code, redirect_uri)                                  │
  ├──────────────────────────────────┼─────────────────────────────────────────────────────────────────────┤
  │ src/stores/auth.ts               │ loginWithGoogle(code, redirectUri)                                  │
  ├──────────────────────────────────┼─────────────────────────────────────────────────────────────────────┤
  │ src/pages/auth/OAuthCallback.vue │ New — handles /auth/callback, reads ?code=, calls store, redirects  │
  ├──────────────────────────────────┼─────────────────────────────────────────────────────────────────────┤
  │ src/router/index.ts              │ /auth/callback route (standalone, no AuthLayout chrome)             │
  ├──────────────────────────────────┼─────────────────────────────────────────────────────────────────────┤
  │ src/pages/auth/Login.vue         │ Button enabled, signInWithGoogle() builds consent URL and redirects │
  ├──────────────────────────────────┼─────────────────────────────────────────────────────────────────────┤
  │ frontend/.env.example            │ New — documents VITE_GOOGLE_CLIENT_ID                               │
  └──────────────────────────────────┴─────────────────────────────────────────────────────────────────────┘

  ---
  Keys needed

  Google Cloud Console (console.cloud.google.com)

  1. APIs & Services → Credentials → Create → OAuth 2.0 Client ID → Web application
  2. Authorised redirect URIs:
    - http://localhost:3000/auth/callback (dev)
    - https://yourdomain.com/auth/callback (prod)
  3. Copy the Client ID and Client Secret

  Backend .env

  GOOGLE_CLIENT_ID=….apps.googleusercontent.com
  GOOGLE_CLIENT_SECRET=GOCSPX-…

  Frontend .env.local

  VITE_GOOGLE_CLIENT_ID=….apps.googleusercontent.com   # same Client ID, safe to expose
  The GOOGLE_CLIENT_SECRET never goes to the frontend.