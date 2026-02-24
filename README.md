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