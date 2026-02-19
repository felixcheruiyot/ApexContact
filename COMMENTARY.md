# Commentary

Local man commentary is always entertaining and powerful. In addition to the features we have, I would like us to introduce commentary.

# Reality

User always meet in bars and other social places to discuss football and other sports. They chat and joke about all kind of things. There is a lot of discussion in the media and etal around games and predictions. This happen everyday. People assemble in match viewing places and talk about their favourite teams and etal. They are good in commentary too.

# Idea

Think about X spaces or any other where people have discussions about their favourite things. They discuss in their languages. We want to bring this to those who can't join the small social groups physically.

# Plan

Creat a section for commentaries or similar. Find a way you can show live commentaries or lobbies where conversations are happenning. Do a good design so that we don't loose what we have built. Anybody in the app can create a commentary lobby. They will watch football in their favourite places-  Just like events, it has to be schedules and users can subscribe. Commentaries will be audios, admin can choose who are the speakers. Just like spaces. Other users should be able to chat and post meme when the event is in progress. If it was a live event, those who did not attend are able to followup. This commentaries are like events but local man contribution of what is happenning. Using AI,or something similar, you can highlight the best comments and event sentiment. This is not a must for now if it will be too much. We should also think about users privacy. We can allow them to add nicknames or event auto-suggest to hide their real details. This should happen through phones and desktop (all devices) i.e the recoding and streaming. Note only audo streams.

# Impementation
 How does commentaries work and how can I access it                                                      
  ⎿  ⧉ Selected 1 lines from ../README.md in Visual Studio Code
                                                                                                          
● Commentary Lobbies — How it works                                                                     
                                                                                                          
  Commentary is an audio-only discussion room — like a sports bar watch party. Anyone can host, listeners 
  join and chat in real-time.                                                                             
                                                                                                          
  ---                                                                                                   
  User Roles in a Lobby

  ┌──────────┬────────────────────────────────────────────────────────────────┐
  │   Role   │                             Can do                             │
  ├──────────┼────────────────────────────────────────────────────────────────┤
  │ Host     │ Creates the lobby, starts/ends the room, grants speaker access │
  ├──────────┼────────────────────────────────────────────────────────────────┤
  │ Speaker  │ Speaks live via microphone                                     │
  ├──────────┼────────────────────────────────────────────────────────────────┤
  │ Listener │ Hears audio, sends chat messages, reacts with emojis           │
  └──────────┴────────────────────────────────────────────────────────────────┘

  ---
  Flow

  For hosts (creating a lobby):
  1. Go to /dashboard/commentary/create (requires login — any role can host)
  2. Fill in title, teaser hook, sport type, date/time, price (0 = free), thumbnail
  3. Submit → lobby is immediately scheduled (no admin review)
  4. Share the link; when ready, open your lobby page and click Start Room
  5. Click End Room when done

  For listeners (joining a lobby):
  1. See Commentary lobbies on the homepage (orange 🎙 section)
  2. Click a lobby card → goes to /commentary/:id
  3. Click Join Free or Buy Access → pick a nickname from suggestions or type your own
  4. If paid: M-Pesa STK Push flow runs first
  5. Once joined, click Join Live Room → /commentary/:id/room
  6. Hear audio, send chat messages, click emoji reactions (🔥 💯 👊 🥊 🏎️  😂 🤯 👀)

  ---
  Access points

  URL: /
  Who: Everyone
  What: Commentary section on homepage
  ────────────────────────────────────────
  URL: /commentary/:id
  Who: Everyone
  What: Lobby detail, join/buy button, chat replay after it ends
  ────────────────────────────────────────
  URL: /commentary/:id/room
  Who: Joined users only
  What: Live audio room + chat
  ────────────────────────────────────────
  URL: /dashboard/commentary/create
  Who: Logged-in users
  What: Create a new lobby

  ---
  After it ends

  Completed lobbies show a Chat Replay on their detail page — all persisted text messages in order. Emoji
  reactions are ephemeral and not saved.
