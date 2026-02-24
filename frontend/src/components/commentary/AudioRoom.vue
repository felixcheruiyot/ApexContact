<template>
  <div class="flex flex-col h-full bg-bg-surface rounded-xl border border-white/5 overflow-hidden">
    <!-- Header -->
    <div class="px-4 py-3 border-b border-white/5 flex items-center justify-between">
      <div class="flex items-center gap-2">
        <span class="text-white font-semibold text-sm">{{ showVideo ? 'Video Room' : 'Audio Room' }}</span>
        <span v-if="connected" class="w-1.5 h-1.5 rounded-full bg-success animate-pulse" />
      </div>
      <span class="text-text-muted text-xs">{{ participants.length }} in room</span>
    </div>

    <!-- Pre-join screen -->
    <div v-if="!joined" class="flex flex-col items-center justify-center flex-1 gap-4 p-6">
      <component :is="showVideo ? Video : Mic" class="w-10 h-10 text-text-muted" />
      <p class="text-white font-medium text-sm">Ready to join?</p>
      <p class="text-text-muted text-xs text-center max-w-xs">
        {{ canPublish
          ? (showVideo ? 'Your microphone and camera will be enabled.' : 'Your microphone will be enabled.')
          : 'You will join as a listener.' }}
      </p>
      <button
        @click="join"
        class="px-6 py-2.5 rounded-lg bg-accent-orange text-white text-sm font-semibold hover:bg-orange-500 transition-colors"
      >
        Join {{ showVideo ? 'Video' : 'Audio' }}
      </button>
    </div>

    <!-- Connecting -->
    <div v-else-if="!connected && !connectionError" class="flex flex-col items-center justify-center flex-1 gap-4">
      <div class="w-12 h-12 rounded-full border-2 border-accent-orange border-t-transparent animate-spin" />
      <p class="text-text-muted text-sm">Connecting…</p>
    </div>

    <!-- Error -->
    <div v-else-if="connectionError" class="flex flex-col items-center justify-center flex-1 gap-3 p-6">
      <MicOff class="w-8 h-8 text-text-muted" />
      <p class="text-white font-medium text-sm">Connection unavailable</p>
      <p class="text-text-muted text-xs text-center max-w-xs">{{ connectionError }}</p>
      <button @click="retry"
        class="mt-2 px-4 py-2 rounded-lg bg-accent-orange text-white text-sm font-medium hover:bg-orange-500 transition-colors">
        Retry
      </button>
    </div>

    <!-- Connected: participants -->
    <div v-else class="flex-1 overflow-hidden flex flex-col">

      <!-- Video grid (audio_video mode) -->
      <div v-if="showVideo && videoParticipants.length" class="grid grid-cols-2 gap-2 p-3 flex-1 overflow-auto">
        <div
          v-for="vp in videoParticipants"
          :key="vp.identity"
          class="relative rounded-lg overflow-hidden bg-black aspect-video"
        >
          <video
            v-if="vp.videoEl"
            :ref="el => attachVideo(el as HTMLVideoElement, vp.identity)"
            autoplay
            playsinline
            muted
            class="w-full h-full object-cover"
          />
          <div v-else class="w-full h-full flex items-center justify-center bg-bg-elevated text-2xl font-bold text-white">
            {{ initials(vp.identity) }}
          </div>
          <div class="absolute bottom-1 left-2 text-white text-xs font-semibold drop-shadow">
            {{ vp.identity }}
            <span v-if="vp.isSpeaking" class="ml-1 text-accent-orange">▶</span>
          </div>
        </div>
      </div>

      <!-- Audio-only: avatar grid -->
      <div v-else class="grid grid-cols-3 sm:grid-cols-4 gap-4 p-4 overflow-y-auto flex-1">
        <div
          v-for="p in participants"
          :key="p.identity"
          class="flex flex-col items-center gap-2"
        >
          <div
            class="relative w-14 h-14 rounded-full flex items-center justify-center text-xl font-bold"
            :class="[
              p.isSpeaking ? 'ring-2 ring-accent-orange ring-offset-2 ring-offset-bg-surface' : 'ring-1 ring-white/10',
              'bg-bg-elevated transition-all duration-200'
            ]"
          >
            {{ initials(p.identity) }}
            <div v-if="p.isSpeaking" class="absolute inset-0 rounded-full bg-accent-orange/20 animate-pulse" />
          </div>
          <span class="text-white text-xs text-center truncate max-w-[72px]" :title="p.identity">
            {{ p.identity }}
          </span>
          <span class="text-text-muted">
            <MicOff v-if="p.isMuted" class="w-3 h-3" />
            <Mic v-else class="w-3 h-3" />
          </span>
          <!-- Host: grant/revoke mic -->
          <template v-if="isHost && p.identity !== myUserId">
            <button
              v-if="(speakerIds ?? []).includes(p.identity)"
              @click="$emit('revokeMic', p.identity)"
              class="px-2 py-0.5 rounded text-[10px] font-medium bg-red-500/20 text-red-400 hover:bg-red-500/30 transition-colors"
            >Revoke Mic</button>
            <button
              v-else
              @click="$emit('grantMic', p.identity)"
              class="px-2 py-0.5 rounded text-[10px] font-medium bg-success/20 text-success hover:bg-success/30 transition-colors"
            >Grant Mic</button>
          </template>
        </div>
      </div>

      <!-- Controls -->
      <div class="px-4 py-3 border-t border-white/5 flex items-center justify-center gap-3 flex-wrap">
        <!-- Mute/unmute -->
        <button
          v-if="canPublish"
          @click="toggleMute"
          class="flex items-center gap-2 px-4 py-2 rounded-lg text-sm font-medium transition-colors"
          :class="muted ? 'bg-red-500/20 text-red-400 hover:bg-red-500/30' : 'bg-bg-elevated text-white hover:bg-white/10'"
        >
          <MicOff v-if="muted" class="w-4 h-4" /><Mic v-else class="w-4 h-4" />
          {{ muted ? 'Unmute' : 'Mute' }}
        </button>

        <!-- Camera toggle (audio_video only) -->
        <button
          v-if="canPublish && showVideo"
          @click="toggleCamera"
          class="flex items-center gap-2 px-4 py-2 rounded-lg text-sm font-medium transition-colors"
          :class="cameraOff ? 'bg-red-500/20 text-red-400 hover:bg-red-500/30' : 'bg-bg-elevated text-white hover:bg-white/10'"
        >
          <VideoOff v-if="cameraOff" class="w-4 h-4" /><Video v-else class="w-4 h-4" />
          {{ cameraOff ? 'Start Camera' : 'Stop Camera' }}
        </button>

        <!-- Raise hand -->
        <button
          v-if="!canPublish"
          @click="$emit('raiseHand')"
          class="flex items-center gap-2 px-4 py-2 rounded-lg bg-bg-elevated text-white text-sm font-medium hover:bg-white/10 transition-colors"
        >
          <Hand class="w-4 h-4" /> Raise Hand
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
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onUnmounted, watch, nextTick } from 'vue'
import { Mic, MicOff, Hand, Video, VideoOff } from 'lucide-vue-next'
import {
  Room,
  RoomEvent,
  Track,
  RemoteParticipant,
} from 'livekit-client'

interface ParticipantInfo {
  identity: string
  isSpeaking: boolean
  isMuted: boolean
  videoEl: boolean
}

const props = defineProps<{
  livekitUrl: string
  token: string
  myRole: 'host' | 'speaker' | 'listener'
  isHost?: boolean
  myUserId?: string
  speakerIds?: string[]
  showVideo?: boolean
}>()

const emit = defineEmits<{
  (e: 'raiseHand'): void
  (e: 'leave'): void
  (e: 'connected'): void
  (e: 'error', msg: string): void
  (e: 'grantMic', userId: string): void
  (e: 'revokeMic', userId: string): void
}>()

let room: Room | null = null
const audioElements: HTMLAudioElement[] = []

const joined = ref(false)
const connected = ref(false)
const connectionError = ref('')
const participants = ref<ParticipantInfo[]>([])
const muted = ref(false)
const cameraOff = ref(false)

const canPublish = computed(() => props.myRole === 'host' || props.myRole === 'speaker')
const videoParticipants = computed(() => participants.value.filter(p => p.videoEl))

function initials(identity: string) {
  return identity.slice(0, 2).toUpperCase()
}

// Attach a video track's stream to a <video> element
function attachVideo(el: HTMLVideoElement | null, identity: string) {
  if (!el) return
  if (!room) return
  const isLocal = room.localParticipant.identity === identity
  if (isLocal) {
    const pub = [...room.localParticipant.videoTrackPublications.values()].find(p => p.track)
    if (pub?.track) {
      const stream = new MediaStream([pub.track.mediaStreamTrack])
      el.srcObject = stream
    }
  } else {
    room.remoteParticipants.forEach((p) => {
      if (p.identity === identity) {
        p.videoTrackPublications.forEach(pub => {
          if (pub.track) {
            const stream = new MediaStream([pub.track.mediaStreamTrack])
            el.srcObject = stream
          }
        })
      }
    })
  }
}

function buildParticipantList() {
  if (!room) return
  const list: ParticipantInfo[] = []

  const local = room.localParticipant
  const localHasVideo = local.videoTrackPublications.size > 0
  list.push({
    identity: local.identity,
    isSpeaking: local.isSpeaking,
    isMuted: !local.isMicrophoneEnabled,
    videoEl: localHasVideo,
  })

  room.remoteParticipants.forEach((p: RemoteParticipant) => {
    const hasVideo = p.videoTrackPublications.size > 0
    list.push({
      identity: p.identity,
      isSpeaking: p.isSpeaking,
      isMuted: !p.audioTrackPublications.size,
      videoEl: hasVideo,
    })
  })

  participants.value = list
}

async function join() {
  joined.value = true
  await connect()
}

function cleanupAudio() {
  audioElements.forEach(el => { el.srcObject = null; el.remove() })
  audioElements.length = 0
}

async function connect() {
  connectionError.value = ''
  try {
    room = new Room()
    await room.connect(props.livekitUrl, props.token)
    await room.startAudio()

    if (room.localParticipant.permissions?.canPublish) {
      await room.localParticipant.setMicrophoneEnabled(true)
      if (props.showVideo) {
        await room.localParticipant.setCameraEnabled(true)
        cameraOff.value = false
      }
    }

    connected.value = true
    buildParticipantList()
    emit('connected')

    room.on(RoomEvent.TrackSubscribed, (track) => {
      if (track.kind === Track.Kind.Audio) {
        const el = track.attach() as HTMLAudioElement
        el.setAttribute('playsinline', '')
        document.body.appendChild(el)
        audioElements.push(el)
      }
      buildParticipantList()
      if (track.kind === Track.Kind.Video) {
        nextTick(() => buildParticipantList())
      }
    })

    room.on(RoomEvent.TrackUnsubscribed, (track) => {
      track.detach()
      buildParticipantList()
    })

    room.on(RoomEvent.ParticipantConnected, buildParticipantList)
    room.on(RoomEvent.ParticipantDisconnected, buildParticipantList)
    room.on(RoomEvent.ActiveSpeakersChanged, buildParticipantList)
    room.on(RoomEvent.TrackPublished, buildParticipantList)
    room.on(RoomEvent.TrackUnpublished, buildParticipantList)
    room.on(RoomEvent.LocalTrackPublished, buildParticipantList)
    room.on(RoomEvent.LocalTrackUnpublished, buildParticipantList)
  } catch (err: any) {
    connectionError.value = err?.message ?? 'Could not connect to the room.'
    emit('error', connectionError.value)
  }
}

async function toggleMute() {
  if (!room) return
  muted.value = !muted.value
  await room.localParticipant.setMicrophoneEnabled(!muted.value)
  buildParticipantList()
}

async function toggleCamera() {
  if (!room) return
  cameraOff.value = !cameraOff.value
  await room.localParticipant.setCameraEnabled(!cameraOff.value)
  buildParticipantList()
}

function retry() {
  cleanupAudio()
  room?.disconnect()
  connect()
}

watch(() => props.token, async (newToken, oldToken) => {
  if (newToken && oldToken && newToken !== oldToken && connected.value) {
    cleanupAudio()
    room?.disconnect()
    connected.value = false
    await connect()
  }
})

onUnmounted(() => {
  cleanupAudio()
  room?.disconnect()
})
</script>
