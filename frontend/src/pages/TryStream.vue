<template>
  <div class="min-h-screen bg-bg flex flex-col">

    <!-- Minimal header -->
    <header class="border-b border-white/5 px-4 sm:px-6 py-3 sm:py-4 flex items-center justify-between gap-3">
      <RouterLink to="/" class="flex items-center gap-2 shrink-0 group min-w-0">
        <img src="@/assets/logo.svg" alt="Live Streamify" class="w-7 h-7 shrink-0" />
        <span class="font-display text-lg sm:text-xl tracking-widest text-white group-hover:text-white transition-colors truncate">
          LIVE <span class="text-accent-red">STREAMIFY</span>
        </span>
      </RouterLink>
      <div class="flex items-center gap-2 sm:gap-3 shrink-0">
        <RouterLink
          to="/login"
          class="text-text-muted hover:text-white text-sm transition-colors whitespace-nowrap"
        >
          Sign in
        </RouterLink>
        <RouterLink
          to="/register"
          class="px-3 sm:px-4 py-1.5 sm:py-2 rounded-lg border border-white/20 hover:border-white/40
                 text-white text-xs sm:text-sm font-medium transition-all hover:bg-white/5 whitespace-nowrap"
        >
          <span class="sm:hidden">Sign up</span>
          <span class="hidden sm:inline">Create account</span>
        </RouterLink>
      </div>
    </header>

    <!-- Main content -->
    <main class="flex-1 flex items-start justify-center px-4 pt-8 sm:pt-12 pb-20">
      <div class="w-full max-w-2xl">

        <!-- ── PHASE: pick ──────────────────────────────────────────── -->
        <div v-if="phase === 'pick'" class="animate-fade-in space-y-8">
          <div>
            <h1 class="font-display text-5xl sm:text-6xl uppercase tracking-wide text-white leading-none mb-3">
              How do you want<br />to stream?
            </h1>
            <p class="text-text-muted text-base leading-relaxed">
              Pick a mode and we'll get you live in seconds. No account needed.
            </p>
          </div>

          <div class="space-y-3">
            <button
              v-for="mode in pickModes"
              :key="mode.type"
              @click="selectedMode = mode.type"
              class="w-full text-left rounded-xl border p-5 transition-all duration-200 group"
              :class="selectedMode === mode.type
                ? 'border-accent-red bg-accent-red/5'
                : 'border-white/10 bg-bg-elevated hover:border-white/25'"
            >
              <div class="flex items-center gap-4">
                <div
                  class="w-11 h-11 rounded-lg flex items-center justify-center shrink-0"
                  :class="mode.iconBg"
                >
                  <component :is="mode.icon" class="w-5 h-5" :class="mode.iconColor" />
                </div>
                <div class="flex-1 min-w-0">
                  <div class="flex items-center gap-2 mb-0.5">
                    <span class="text-white font-semibold text-sm">{{ mode.label }}</span>
                    <span class="text-xs font-semibold px-1.5 py-0.5 rounded" :class="mode.badgeClass">
                      {{ mode.badge }}
                    </span>
                  </div>
                  <p class="text-text-muted text-xs leading-relaxed">{{ mode.description }}</p>
                </div>
                <div
                  class="w-5 h-5 rounded-full border-2 flex items-center justify-center shrink-0 transition-colors"
                  :class="selectedMode === mode.type ? 'border-accent-red bg-accent-red' : 'border-white/20'"
                >
                  <svg v-if="selectedMode === mode.type" class="w-3 h-3 text-white" fill="currentColor" viewBox="0 0 24 24">
                    <circle cx="12" cy="12" r="6" />
                  </svg>
                </div>
              </div>
            </button>
          </div>

          <div class="bg-bg-elevated border border-white/10 rounded-xl p-6 space-y-5">
            <div>
              <label class="block text-white text-sm font-medium mb-2">
                What will you be streaming about?
                <span class="text-text-muted font-normal ml-1">(optional)</span>
              </label>
              <input
                v-model="streamTitle"
                type="text"
                placeholder="e.g. Sales masterclass, Visa Q&A, Live music session…"
                class="input"
                :disabled="loading"
                @keydown.enter="startStream"
              />
            </div>

            <div v-if="error" class="flex items-start gap-3 bg-status-error/10 border border-status-error/30
                                      rounded-lg p-4 text-status-error text-sm">
              <svg class="w-4 h-4 shrink-0 mt-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                  d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              {{ error }}
            </div>

            <button
              @click="startStream"
              :disabled="loading || !selectedMode"
              class="w-full flex items-center justify-center gap-2 px-6 py-4 rounded-lg
                     bg-accent-red hover:bg-accent-red-hover text-white font-bold text-base
                     transition-all duration-200 active:scale-95 disabled:opacity-60 disabled:cursor-not-allowed"
            >
              <svg v-if="loading" class="w-5 h-5 animate-spin" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"/>
              </svg>
              <component v-else :is="currentModeIcon" class="w-5 h-5" />
              {{ loading ? 'Setting up your stream…' : 'Start Test Stream' }}
            </button>
          </div>

          <div class="flex flex-wrap items-center gap-x-6 gap-y-2 text-text-muted text-sm">
            <span class="flex items-center gap-2">
              <svg class="w-4 h-4 text-status-success" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
              </svg>
              5 minutes free
            </span>
            <span class="flex items-center gap-2">
              <svg class="w-4 h-4 text-status-success" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
              </svg>
              No account needed
            </span>
            <span class="flex items-center gap-2">
              <svg class="w-4 h-4 text-status-success" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
              </svg>
              No credit card
            </span>
          </div>
        </div>

        <!-- ── PHASE: active (OBS mode) ───────────────────────────── -->
        <div v-if="phase === 'active' && selectedMode === 'commercial'" class="animate-fade-in space-y-6">
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-2.5">
              <span class="w-2 h-2 rounded-full bg-accent-red animate-pulse" />
              <span class="text-white font-semibold text-sm uppercase tracking-wider">OBS Stream Active</span>
            </div>
            <div class="font-display text-2xl tracking-wider" :class="secondsLeft < 120 ? 'text-status-warning' : 'text-white'">
              {{ formattedTime }}
            </div>
          </div>

          <div v-if="secondsLeft < 120 && secondsLeft > 0"
               class="flex items-center gap-3 bg-status-warning/10 border border-status-warning/30
                      rounded-lg px-4 py-3 text-status-warning text-sm">
            <svg class="w-4 h-4 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            Less than 2 minutes left.
            <RouterLink to="/register" class="font-semibold underline hover:no-underline ml-1">
              Create a free account
            </RouterLink>
            to keep streaming.
          </div>

          <!-- RTMP credentials -->
          <div class="bg-bg-elevated border border-white/10 rounded-xl p-6 space-y-4">
            <h2 class="text-white font-bold text-base mb-1">Stream with OBS or any RTMP encoder</h2>
            <p class="text-text-muted text-sm mb-4">
              Open OBS → Settings → Stream → Select "Custom", then paste the values below.
            </p>
            <div class="space-y-3">
              <div>
                <label class="text-text-muted text-xs font-medium uppercase tracking-wider mb-1.5 block">
                  Server (RTMP URL)
                </label>
                <div class="flex items-center gap-2">
                  <code class="flex-1 bg-bg border border-white/10 rounded-lg px-4 py-2.5 text-white text-sm font-mono truncate">
                    {{ guestStream?.rtmp_url }}
                  </code>
                  <button @click="copy(guestStream?.rtmp_url ?? '', 'rtmp')"
                    class="shrink-0 px-3 py-2.5 rounded-lg border border-white/15 hover:border-white/30
                           text-text-muted hover:text-white text-xs font-medium transition-all">
                    {{ copied === 'rtmp' ? 'Copied!' : 'Copy' }}
                  </button>
                </div>
              </div>
              <div>
                <label class="text-text-muted text-xs font-medium uppercase tracking-wider mb-1.5 block">
                  Stream Key
                </label>
                <div class="flex items-center gap-2">
                  <code class="flex-1 bg-bg border border-white/10 rounded-lg px-3 sm:px-4 py-2.5 text-white text-sm font-mono truncate min-w-0">
                    {{ showKey ? guestStream?.stream_key : '••••••••••••' }}
                  </code>
                  <div class="flex items-center gap-1.5 shrink-0">
                    <button @click="showKey = !showKey"
                      class="px-3 py-2.5 rounded-lg border border-white/15 hover:border-white/30
                             text-text-muted hover:text-white text-xs font-medium transition-all whitespace-nowrap">
                      {{ showKey ? 'Hide' : 'Show' }}
                    </button>
                    <button @click="copy(guestStream?.stream_key ?? '', 'key')"
                      class="px-3 py-2.5 rounded-lg border border-white/15 hover:border-white/30
                             text-text-muted hover:text-white text-xs font-medium transition-all whitespace-nowrap">
                      {{ copied === 'key' ? 'Copied!' : 'Copy' }}
                    </button>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- Viewer link -->
          <div class="bg-bg-elevated border border-white/10 rounded-xl p-6 space-y-3">
            <div>
              <h2 class="text-white font-bold text-base mb-1">Share with your audience</h2>
              <p class="text-text-muted text-sm">They can watch right now — no account needed on their end.</p>
            </div>
            <div class="flex items-center gap-2">
              <code class="flex-1 bg-bg border border-white/10 rounded-lg px-4 py-2.5 text-white text-sm font-mono truncate">
                {{ guestStream?.viewer_url }}
              </code>
              <button @click="copy(guestStream?.viewer_url ?? '', 'viewer')"
                class="shrink-0 px-4 py-2.5 rounded-lg bg-accent-red hover:bg-accent-red-hover
                       text-white text-xs font-semibold transition-all">
                {{ copied === 'viewer' ? 'Copied!' : 'Copy link' }}
              </button>
            </div>
          </div>

          <button @click="clearSession(); phase = 'expired'" class="text-text-muted hover:text-white text-sm transition-colors underline">
            End test stream
          </button>
        </div>

        <!-- ── PHASE: active (audio / audio_video mode) ───────────── -->
        <div v-if="phase === 'active' && (selectedMode === 'audio' || selectedMode === 'audio_video')" class="animate-fade-in space-y-6">
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-2.5">
              <span class="w-2 h-2 rounded-full bg-accent-orange animate-pulse" />
              <span class="text-white font-semibold text-sm uppercase tracking-wider">
                {{ selectedMode === 'audio' ? 'Audio Room Active' : 'Audio + Video Room Active' }}
              </span>
            </div>
            <div class="font-display text-2xl tracking-wider" :class="secondsLeft < 120 ? 'text-status-warning' : 'text-white'">
              {{ formattedTime }}
            </div>
          </div>

          <div v-if="secondsLeft < 120 && secondsLeft > 0"
               class="flex items-center gap-3 bg-status-warning/10 border border-status-warning/30
                      rounded-lg px-4 py-3 text-status-warning text-sm">
            <svg class="w-4 h-4 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            Less than 2 minutes left.
            <RouterLink to="/register" class="font-semibold underline hover:no-underline ml-1">
              Create a free account
            </RouterLink>
            to keep streaming.
          </div>

          <!-- LiveKit room info -->
          <div class="bg-bg-elevated border border-white/10 rounded-xl p-6 space-y-4">
            <div class="flex items-center gap-3 mb-2">
              <div class="w-10 h-10 rounded-lg bg-accent-orange/15 flex items-center justify-center">
                <component :is="selectedMode === 'audio' ? Mic : Video" class="w-5 h-5 text-accent-orange" />
              </div>
              <div>
                <h2 class="text-white font-bold text-base">Your live room is ready</h2>
                <p class="text-text-muted text-xs">
                  {{ selectedMode === 'audio' ? 'Mic-only interactive session' : 'Camera + mic interactive session' }}
                </p>
              </div>
            </div>

            <div class="bg-bg border border-white/8 rounded-lg p-4">
              <p class="text-text-muted text-xs uppercase tracking-wider mb-2 font-medium">Room name</p>
              <code class="text-white text-sm font-mono">{{ guestRoom?.room_name }}</code>
            </div>

            <a
              :href="guestRoom?.viewer_url"
              target="_blank"
              class="inline-flex items-center gap-2 px-5 py-2.5 rounded-lg bg-accent-orange hover:bg-orange-500
                     text-white text-sm font-semibold transition-all"
            >
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                  d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14" />
              </svg>
              Open my room
            </a>

            <p class="text-text-muted text-xs">
              Your browser will ask for microphone{{ selectedMode === 'audio_video' ? ' and camera' : '' }} permission.
              The room opens in a new tab.
            </p>
          </div>

          <!-- Viewer share link -->
          <div class="bg-bg-elevated border border-white/10 rounded-xl p-6 space-y-3">
            <div>
              <h2 class="text-white font-bold text-base mb-1">Share with your audience</h2>
              <p class="text-text-muted text-sm">They join as listeners — no account needed.</p>
            </div>
            <div class="flex items-center gap-2">
              <code class="flex-1 bg-bg border border-white/10 rounded-lg px-4 py-2.5 text-white text-sm font-mono truncate">
                {{ guestRoom?.viewer_url }}
              </code>
              <button @click="copy(guestRoom?.viewer_url ?? '', 'viewer')"
                class="shrink-0 px-4 py-2.5 rounded-lg bg-accent-red hover:bg-accent-red-hover
                       text-white text-xs font-semibold transition-all">
                {{ copied === 'viewer' ? 'Copied!' : 'Copy link' }}
              </button>
            </div>
          </div>

          <button @click="clearSession(); phase = 'expired'" class="text-text-muted hover:text-white text-sm transition-colors underline">
            End test room
          </button>
        </div>

        <!-- ── PHASE: expired ─────────────────────────────────────── -->
        <div v-if="phase === 'expired'" class="animate-fade-in space-y-8">
          <div>
            <div class="inline-flex items-center gap-2 bg-white/5 border border-white/10 text-text-muted
                        text-xs font-semibold px-3 py-1.5 rounded mb-6">
              Test stream ended
            </div>
            <h1 class="font-display text-5xl sm:text-6xl uppercase tracking-wide text-white leading-none mb-4">
              Did it work?
            </h1>
            <p class="text-text-muted text-base leading-relaxed max-w-md">
              If your stream ran and your audience could watch, the platform works for you.
              Create a free account to go live without limits and start charging.
            </p>
          </div>

          <div class="bg-bg-elevated border border-white/10 rounded-xl p-8 space-y-6">
            <div class="space-y-3">
              <div class="flex items-start gap-3 text-sm">
                <svg class="w-5 h-5 text-status-success shrink-0 mt-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
                </svg>
                <span class="text-white">Unlimited stream time</span>
              </div>
              <div class="flex items-start gap-3 text-sm">
                <svg class="w-5 h-5 text-status-success shrink-0 mt-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
                </svg>
                <span class="text-white">Set a ticket price — viewers pay via M-Pesa</span>
              </div>
              <div class="flex items-start gap-3 text-sm">
                <svg class="w-5 h-5 text-status-success shrink-0 mt-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
                </svg>
                <span class="text-white">70% payout to you after every session</span>
              </div>
              <div class="flex items-start gap-3 text-sm">
                <svg class="w-5 h-5 text-status-success shrink-0 mt-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
                </svg>
                <span class="text-white">Anti-piracy — only paid audience can watch</span>
              </div>
            </div>

            <div class="flex flex-col sm:flex-row gap-3 pt-2">
              <RouterLink
                to="/register"
                class="flex-1 text-center px-6 py-3.5 rounded-lg bg-accent-red
                       hover:bg-accent-red-hover text-white font-bold text-sm
                       transition-all duration-200 active:scale-95"
              >
                Create Free Account
              </RouterLink>
              <button
                @click="resetAndTryAgain"
                class="flex-1 text-center px-6 py-3.5 rounded-lg border border-white/20
                       hover:border-white/40 text-white text-sm font-medium
                       transition-all hover:bg-white/5"
              >
                Try again
              </button>
            </div>
          </div>
        </div>

      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { RouterLink } from 'vue-router'
import { Mic, Video, Monitor } from 'lucide-vue-next'
import { createGuestStream, createGuestRoom, type GuestStream, type GuestRoom } from '@/api/stream'

type Phase = 'pick' | 'active' | 'expired'
type ModeType = 'audio_video' | 'audio' | 'commercial'

const STORAGE_KEY = 'apex_guest_session'

const phase = ref<Phase>('pick')
const selectedMode = ref<ModeType>('audio_video')
const streamTitle = ref('')
const loading = ref(false)
const error = ref('')
const guestStream = ref<GuestStream | null>(null)
const guestRoom = ref<GuestRoom | null>(null)
const secondsLeft = ref(0)
const copied = ref<string | null>(null)
const showKey = ref(false)

let countdown: ReturnType<typeof setInterval> | null = null

const pickModes = [
  {
    type: 'audio_video' as ModeType,
    label: 'Phone / PC — Audio + Video',
    badge: 'Interactive',
    badgeClass: 'bg-accent-orange/20 text-accent-orange',
    description: 'Stream from your browser or phone camera. Viewers join as listeners in a live interactive room.',
    icon: Video,
    iconBg: 'bg-accent-orange/15',
    iconColor: 'text-accent-orange',
  },
  {
    type: 'audio' as ModeType,
    label: 'Phone / PC — Audio Only',
    badge: 'Audio room',
    badgeClass: 'bg-accent-red/20 text-accent-red',
    description: 'Mic-only session — like a podcast room. No camera needed. Great for discussions and panels.',
    icon: Mic,
    iconBg: 'bg-accent-red/15',
    iconColor: 'text-accent-red',
  },
  {
    type: 'commercial' as ModeType,
    label: 'OBS / RTMP Encoder',
    badge: 'Professional',
    badgeClass: 'bg-white/10 text-text-muted',
    description: 'Use OBS Studio or any RTMP encoder. High-quality one-to-many broadcast, watched via HLS.',
    icon: Monitor,
    iconBg: 'bg-white/10',
    iconColor: 'text-text-muted',
  },
]

const currentModeIcon = computed(() => {
  if (selectedMode.value === 'audio') return Mic
  if (selectedMode.value === 'audio_video') return Video
  return Monitor
})

const formattedTime = computed(() => {
  const m = Math.floor(secondsLeft.value / 60)
  const s = secondsLeft.value % 60
  return `${String(m).padStart(2, '0')}:${String(s).padStart(2, '0')}`
})

function saveSession() {
  const expiresAt = guestStream.value?.expires_at ?? guestRoom.value?.expires_at
  localStorage.setItem(STORAGE_KEY, JSON.stringify({
    mode: selectedMode.value,
    expires_at: expiresAt,
    stream: guestStream.value,
    room: guestRoom.value,
  }))
}

function clearSession() {
  localStorage.removeItem(STORAGE_KEY)
}

onMounted(() => {
  try {
    const raw = localStorage.getItem(STORAGE_KEY)
    if (!raw) return
    const saved = JSON.parse(raw)
    if (!saved.expires_at) return

    const remaining = Math.floor((new Date(saved.expires_at).getTime() - Date.now()) / 1000)

    if (remaining > 0) {
      selectedMode.value = saved.mode
      guestStream.value = saved.stream ?? null
      guestRoom.value = saved.room ?? null
      secondsLeft.value = remaining
      phase.value = 'active'
      startCountdown()
    } else {
      clearSession()
      phase.value = 'expired'
    }
  } catch {
    // ignore corrupt storage
  }
})

async function startStream() {
  if (loading.value || !selectedMode.value) return
  loading.value = true
  error.value = ''

  try {
    if (selectedMode.value === 'commercial') {
      const stream = await createGuestStream(streamTitle.value || undefined)
      guestStream.value = stream
      secondsLeft.value = stream.time_limit_seconds
    } else {
      const room = await createGuestRoom(selectedMode.value, streamTitle.value || undefined)
      guestRoom.value = room
      secondsLeft.value = room.time_limit_seconds
    }
    phase.value = 'active'
    saveSession()
    startCountdown()
  } catch (err: any) {
    const msg = err?.response?.data?.error
    error.value = msg || 'Could not start the stream. Please try again.'
  } finally {
    loading.value = false
  }
}

function startCountdown() {
  clearCountdown()
  countdown = setInterval(() => {
    if (secondsLeft.value <= 1) {
      clearCountdown()
      clearSession()
      phase.value = 'expired'
    } else {
      secondsLeft.value--
    }
  }, 1000)
}

function clearCountdown() {
  if (countdown !== null) {
    clearInterval(countdown)
    countdown = null
  }
}

function resetAndTryAgain() {
  clearSession()
  guestStream.value = null
  guestRoom.value = null
  streamTitle.value = ''
  error.value = ''
  secondsLeft.value = 0
  phase.value = 'pick'
}

async function copy(text: string, key: string) {
  try {
    await navigator.clipboard.writeText(text)
    copied.value = key
    setTimeout(() => { copied.value = null }, 2000)
  } catch {
    // clipboard not available
  }
}

onUnmounted(clearCountdown)
</script>
