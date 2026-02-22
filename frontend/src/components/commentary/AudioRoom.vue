<template>
  <div class="flex flex-col h-full bg-bg-surface rounded-xl border border-white/5 overflow-hidden">
    <!-- Header -->
    <div class="px-4 py-3 border-b border-white/5 flex items-center justify-between">
      <div class="flex items-center gap-2">
        <span class="text-white font-semibold text-sm">Audio Room</span>
        <span v-if="connected" class="w-1.5 h-1.5 rounded-full bg-success animate-pulse" />
      </div>
      <span class="text-text-muted text-xs">{{ participants.length }} in room</span>
    </div>

    <!-- Participants grid -->
    <div class="flex-1 overflow-y-auto p-4">
      <!-- Pre-join: user must click to allow AudioContext + mic -->
      <div v-if="!joined" class="flex flex-col items-center justify-center h-full gap-4">
        <span class="text-4xl">🎙</span>
        <p class="text-white font-medium text-sm">Ready to join the audio room?</p>
        <p class="text-text-muted text-xs text-center max-w-xs">
          {{ canPublish ? 'Your microphone will be enabled.' : 'You will join as a listener.' }}
        </p>
        <button
          @click="join"
          class="px-6 py-2.5 rounded-lg bg-accent-orange text-white text-sm font-semibold hover:bg-orange-500 transition-colors"
        >
          Join Audio
        </button>
      </div>

      <div v-else-if="!connected && !connectionError" class="flex flex-col items-center justify-center h-full gap-4">
        <div class="w-12 h-12 rounded-full border-2 border-accent-orange border-t-transparent animate-spin" />
        <p class="text-text-muted text-sm">Connecting to audio room…</p>
      </div>

      <div v-else-if="connectionError" class="flex flex-col items-center justify-center h-full gap-3">
        <span class="text-3xl">🔇</span>
        <p class="text-white font-medium text-sm">Audio unavailable</p>
        <p class="text-text-muted text-xs text-center max-w-xs">{{ connectionError }}</p>
        <button @click="retry" class="mt-2 px-4 py-2 rounded-lg bg-accent-orange text-white text-sm font-medium hover:bg-orange-500 transition-colors">
          Retry
        </button>
      </div>

      <div v-else class="grid grid-cols-3 sm:grid-cols-4 gap-4">
        <div
          v-for="p in participants"
          :key="p.identity"
          class="flex flex-col items-center gap-2"
        >
          <!-- Avatar with speaking indicator -->
          <div
            class="relative w-14 h-14 rounded-full flex items-center justify-center text-xl font-bold"
            :class="[
              p.isSpeaking ? 'ring-2 ring-accent-orange ring-offset-2 ring-offset-bg-surface' : 'ring-1 ring-white/10',
              'bg-bg-elevated transition-all duration-200'
            ]"
          >
            {{ initials(p.identity) }}
            <!-- Speaking animation ring -->
            <div v-if="p.isSpeaking" class="absolute inset-0 rounded-full bg-accent-orange/20 animate-pulse" />
          </div>
          <span class="text-white text-xs text-center truncate max-w-[72px]" :title="p.identity">
            {{ p.identity }}
          </span>
          <!-- Mic status -->
          <span class="text-text-muted text-[10px]">
            {{ p.isMuted ? '🔇' : '🎙' }}
          </span>
          <!-- Host: grant/revoke mic button (not for self) -->
          <template v-if="isHost && p.identity !== myUserId">
            <button
              v-if="(speakerIds ?? []).includes(p.identity)"
              @click="$emit('revokeMic', p.identity)"
              class="px-2 py-0.5 rounded text-[10px] font-medium bg-red-500/20 text-red-400 hover:bg-red-500/30 transition-colors"
            >
              Revoke Mic
            </button>
            <button
              v-else
              @click="$emit('grantMic', p.identity)"
              class="px-2 py-0.5 rounded text-[10px] font-medium bg-success/20 text-success hover:bg-success/30 transition-colors"
            >
              Grant Mic
            </button>
          </template>
        </div>
      </div>
    </div>

    <!-- Controls -->
    <div v-if="connected" class="px-4 py-3 border-t border-white/5 flex items-center justify-center gap-3">
      <!-- Mute/unmute (speakers & host) -->
      <button
        v-if="canPublish"
        @click="toggleMute"
        class="flex items-center gap-2 px-4 py-2 rounded-lg text-sm font-medium transition-colors"
        :class="muted ? 'bg-red-500/20 text-red-400 hover:bg-red-500/30' : 'bg-bg-elevated text-white hover:bg-white/10'"
      >
        {{ muted ? '🔇 Unmute' : '🎙 Mute' }}
      </button>

      <!-- Raise hand (listeners) -->
      <button
        v-if="!canPublish"
        @click="$emit('raiseHand')"
        class="flex items-center gap-2 px-4 py-2 rounded-lg bg-bg-elevated text-white text-sm font-medium hover:bg-white/10 transition-colors"
      >
        ✋ Raise Hand
      </button>

      <!-- Leave -->
      <button
        @click="$emit('leave')"
        class="px-4 py-2 rounded-lg bg-red-500/20 text-red-400 text-sm font-medium hover:bg-red-500/30 transition-colors"
      >
        Leave
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onUnmounted, computed, watch } from 'vue'
import {
  Room,
  RoomEvent,
  RemoteParticipant,
} from 'livekit-client'

interface ParticipantInfo {
  identity: string
  isSpeaking: boolean
  isMuted: boolean
}

const props = defineProps<{
  livekitUrl: string
  token: string
  myRole: 'host' | 'speaker' | 'listener'
  isHost?: boolean
  myUserId?: string
  speakerIds?: string[]
}>()

const emit = defineEmits<{
  (e: 'raiseHand'): void
  (e: 'leave'): void
  (e: 'connected'): void
  (e: 'error', msg: string): void
  (e: 'grantMic', userId: string): void
  (e: 'revokeMic', userId: string): void
}>()

// Room is created lazily inside connect() so no AudioContext is created
// until the user explicitly clicks "Join Audio".
let room: Room | null = null

const joined = ref(false)
const connected = ref(false)
const connectionError = ref('')
const participants = ref<ParticipantInfo[]>([])
const muted = ref(false)

const canPublish = computed(() => props.myRole === 'host' || props.myRole === 'speaker')

function initials(identity: string) {
  return identity.slice(0, 2).toUpperCase()
}

function buildParticipantList() {
  if (!room) return
  const list: ParticipantInfo[] = []

  const local = room.localParticipant
  list.push({
    identity: local.identity,
    isSpeaking: local.isSpeaking,
    isMuted: local.isMicrophoneEnabled === false,
  })

  room.remoteParticipants.forEach((p: RemoteParticipant) => {
    list.push({
      identity: p.identity,
      isSpeaking: p.isSpeaking,
      isMuted: !p.audioTrackPublications.size,
    })
  })

  participants.value = list
}

// Called on explicit user click — satisfies browser autoplay + getUserMedia policies.
async function join() {
  joined.value = true
  await connect()
}

async function connect() {
  connectionError.value = ''
  try {
    room = new Room()
    await room.connect(props.livekitUrl, props.token)

    if (canPublish.value) {
      await room.localParticipant.setMicrophoneEnabled(true)
    }

    connected.value = true
    buildParticipantList()
    emit('connected')

    room.on(RoomEvent.ParticipantConnected, buildParticipantList)
    room.on(RoomEvent.ParticipantDisconnected, buildParticipantList)
    room.on(RoomEvent.ActiveSpeakersChanged, buildParticipantList)
    room.on(RoomEvent.TrackPublished, buildParticipantList)
    room.on(RoomEvent.TrackUnpublished, buildParticipantList)
  } catch (err: any) {
    connectionError.value = err?.message ?? 'Could not connect to the audio room.'
    emit('error', connectionError.value)
  }
}

async function toggleMute() {
  if (!room) return
  muted.value = !muted.value
  await room.localParticipant.setMicrophoneEnabled(!muted.value)
  buildParticipantList()
}

function retry() {
  room?.disconnect()
  connect()
}

// When the token changes (e.g. role upgraded from listener → speaker),
// reconnect to LiveKit so the new permissions take effect immediately.
watch(() => props.token, async (newToken, oldToken) => {
  if (newToken && oldToken && newToken !== oldToken && connected.value) {
    room?.disconnect()
    connected.value = false
    await connect()
  }
})

onUnmounted(() => {
  room?.disconnect()
})
</script>
