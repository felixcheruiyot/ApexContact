<template>
  <div class="min-h-screen bg-bg flex flex-col">

    <!-- Minimal header -->
    <header class="border-b border-white/5 px-6 py-4 flex items-center justify-between">
      <RouterLink to="/" class="font-display text-xl tracking-widest text-white">
        LIVE <span class="text-accent-red">·</span> STREAMIFY
      </RouterLink>
      <div class="flex items-center gap-3">
        <RouterLink
          to="/login"
          class="text-text-muted hover:text-white text-sm transition-colors"
        >
          Sign in
        </RouterLink>
        <RouterLink
          to="/register"
          class="px-4 py-2 rounded-lg border border-white/20 hover:border-white/40
                 text-white text-sm font-medium transition-all hover:bg-white/5"
        >
          Create account
        </RouterLink>
      </div>
    </header>

    <!-- Main content -->
    <main class="flex-1 flex items-start justify-center px-4 pt-12 pb-20">
      <div class="w-full max-w-2xl">

        <!-- ── PHASE 1: Setup form ──────────────────────────────────── -->
        <div v-if="phase === 'setup'" class="animate-fade-in space-y-8">
          <div>
            <h1 class="font-display text-5xl sm:text-6xl uppercase tracking-wide text-white leading-none mb-3">
              Test Your Stream
            </h1>
            <p class="text-text-muted text-base leading-relaxed">
              Get a real stream key and a shareable viewer link in seconds.
              No account. No payment. 5 minutes free.
            </p>
          </div>

          <div class="bg-bg-elevated border border-white/10 rounded-xl p-8 space-y-6">
            <div>
              <label class="block text-white text-sm font-medium mb-2">
                What will you be streaming about?
                <span class="text-text-muted font-normal ml-1">(optional)</span>
              </label>
              <input
                v-model="streamTitle"
                type="text"
                placeholder="e.g. Sales masterclass, Visa Q&A, Cooking demo…"
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
              :disabled="loading"
              class="w-full flex items-center justify-center gap-2 px-6 py-4 rounded-lg
                     bg-accent-red hover:bg-accent-red-hover text-white font-bold text-base
                     transition-all duration-200 active:scale-95 disabled:opacity-60 disabled:cursor-not-allowed"
            >
              <svg v-if="loading" class="w-5 h-5 animate-spin" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"/>
              </svg>
              <svg v-else class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                  d="M15 10l4.553-2.276A1 1 0 0121 8.677v6.646a1 1 0 01-1.447.894L15 14M5 18h8a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v8a2 2 0 002 2z" />
              </svg>
              {{ loading ? 'Setting up your stream…' : 'Start Test Stream' }}
            </button>
          </div>

          <div class="flex items-center gap-6 text-text-muted text-sm">
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
              Works with OBS
            </span>
            <span class="flex items-center gap-2">
              <svg class="w-4 h-4 text-status-success" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
              </svg>
              No credit card
            </span>
          </div>
        </div>

        <!-- ── PHASE 2: Active stream ──────────────────────────────── -->
        <div v-if="phase === 'active'" class="animate-fade-in space-y-6">

          <!-- Timer + status bar -->
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-2.5">
              <span class="w-2 h-2 rounded-full bg-accent-red animate-pulse" />
              <span class="text-white font-semibold text-sm uppercase tracking-wider">Test Stream Active</span>
            </div>
            <div
              class="font-display text-2xl tracking-wider"
              :class="secondsLeft < 120 ? 'text-status-warning' : 'text-white'"
            >
              {{ formattedTime }}
            </div>
          </div>

          <!-- Soft warning when 2 min left -->
          <div
            v-if="secondsLeft < 120 && secondsLeft > 0"
            class="flex items-center gap-3 bg-status-warning/10 border border-status-warning/30
                   rounded-lg px-4 py-3 text-status-warning text-sm"
          >
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

          <!-- Stream credentials (broadcaster) -->
          <div class="bg-bg-elevated border border-white/10 rounded-xl p-6 space-y-4">
            <h2 class="text-white font-bold text-base mb-1">Stream with OBS or any RTMP software</h2>
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
                  <button
                    @click="copy(guestStream?.rtmp_url ?? '', 'rtmp')"
                    class="shrink-0 px-3 py-2.5 rounded-lg border border-white/15 hover:border-white/30
                           text-text-muted hover:text-white text-xs font-medium transition-all"
                  >
                    {{ copied === 'rtmp' ? 'Copied!' : 'Copy' }}
                  </button>
                </div>
              </div>

              <div>
                <label class="text-text-muted text-xs font-medium uppercase tracking-wider mb-1.5 block">
                  Stream Key
                </label>
                <div class="flex items-center gap-2">
                  <code class="flex-1 bg-bg border border-white/10 rounded-lg px-4 py-2.5 text-white text-sm font-mono truncate">
                    {{ showKey ? guestStream?.stream_key : '••••••••••••' }}
                  </code>
                  <button
                    @click="showKey = !showKey"
                    class="shrink-0 px-3 py-2.5 rounded-lg border border-white/15 hover:border-white/30
                           text-text-muted hover:text-white text-xs font-medium transition-all"
                  >
                    {{ showKey ? 'Hide' : 'Show' }}
                  </button>
                  <button
                    @click="copy(guestStream?.stream_key ?? '', 'key')"
                    class="shrink-0 px-3 py-2.5 rounded-lg border border-white/15 hover:border-white/30
                           text-text-muted hover:text-white text-xs font-medium transition-all"
                  >
                    {{ copied === 'key' ? 'Copied!' : 'Copy' }}
                  </button>
                </div>
              </div>
            </div>
          </div>

          <!-- Viewer share link -->
          <div class="bg-bg-elevated border border-white/10 rounded-xl p-6 space-y-3">
            <div>
              <h2 class="text-white font-bold text-base mb-1">Share with your audience</h2>
              <p class="text-text-muted text-sm">They can watch right now — no account needed on their end.</p>
            </div>
            <div class="flex items-center gap-2">
              <code class="flex-1 bg-bg border border-white/10 rounded-lg px-4 py-2.5 text-white text-sm font-mono truncate">
                {{ guestStream?.viewer_url }}
              </code>
              <button
                @click="copy(guestStream?.viewer_url ?? '', 'viewer')"
                class="shrink-0 px-4 py-2.5 rounded-lg bg-accent-red hover:bg-accent-red-hover
                       text-white text-xs font-semibold transition-all"
              >
                {{ copied === 'viewer' ? 'Copied!' : 'Copy link' }}
              </button>
            </div>
          </div>

          <!-- End test early -->
          <button
            @click="phase = 'expired'"
            class="text-text-muted hover:text-white text-sm transition-colors underline"
          >
            End test stream
          </button>
        </div>

        <!-- ── PHASE 3: Expired / upsell ──────────────────────────── -->
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
              If you could stream and your audience could watch, the platform works for you.
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
                <span class="text-white">Anti-piracy — no sharing stream links outside your paid audience</span>
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
import { ref, computed, onUnmounted } from 'vue'
import { RouterLink } from 'vue-router'
import { createGuestStream, type GuestStream } from '@/api/stream'

type Phase = 'setup' | 'active' | 'expired'

const phase = ref<Phase>('setup')
const streamTitle = ref('')
const loading = ref(false)
const error = ref('')
const guestStream = ref<GuestStream | null>(null)
const secondsLeft = ref(0)
const copied = ref<string | null>(null)
const showKey = ref(false)

let countdown: ReturnType<typeof setInterval> | null = null

const formattedTime = computed(() => {
  const m = Math.floor(secondsLeft.value / 60)
  const s = secondsLeft.value % 60
  return `${String(m).padStart(2, '0')}:${String(s).padStart(2, '0')}`
})

async function startStream() {
  if (loading.value) return
  loading.value = true
  error.value = ''

  try {
    const stream = await createGuestStream(streamTitle.value || undefined)
    guestStream.value = stream
    secondsLeft.value = stream.time_limit_seconds
    phase.value = 'active'
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
  guestStream.value = null
  streamTitle.value = ''
  error.value = ''
  secondsLeft.value = 0
  phase.value = 'setup'
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
