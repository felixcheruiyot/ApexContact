<template>
  <div class="h-screen bg-bg-page flex flex-col overflow-hidden">
    <!-- Top bar -->
    <div class="flex items-center justify-between px-4 py-3 bg-bg-surface border-b border-white/5">
      <RouterLink :to="`/commentary/${eventId}`" class="flex items-center gap-1.5 text-text-muted hover:text-white transition-colors">
        <ArrowLeft class="w-4 h-4" /> Back
      </RouterLink>
      <div class="flex items-center gap-2">
        <span v-if="event" class="text-white font-semibold text-sm truncate max-w-[200px]">{{ event.title }}</span>
        <span class="flex items-center gap-1 px-2 py-0.5 rounded-full bg-accent-orange/20 text-accent-orange text-xs font-bold">
          <span class="w-1.5 h-1.5 rounded-full bg-accent-orange animate-pulse" />
          Live
        </span>
      </div>
      <!-- Host controls -->
      <div v-if="isHost" class="flex items-center gap-2">
        <button
          @click="endRoom"
          class="px-3 py-1.5 rounded-lg bg-red-500/20 text-red-400 text-xs font-medium hover:bg-red-500/30 transition-colors"
        >
          End Room
        </button>
      </div>
      <div v-else class="w-20" />
    </div>

    <!-- Main layout -->
    <div class="flex-1 min-h-0 flex flex-col md:flex-row overflow-hidden">
      <!-- Audio Room (left / top on mobile) -->
      <div class="flex-1 min-h-0 md:flex-none md:w-[60%] p-4 flex flex-col gap-4 overflow-y-auto">
        <AudioRoom
          v-if="livekitToken && livekitUrl"
          :livekitUrl="livekitUrl"
          :token="livekitToken"
          :myRole="myRole"
          :isHost="isHost"
          :myUserId="auth.user?.id"
          :speakerIds="speakerIds"
          :showVideo="event?.event_type === 'audio_video'"
          @leave="leave"
          @raiseHand="raiseHand"
          @grantMic="grantMic"
          @revokeMic="revokeMic"
          @connected="roomConnected = true"
          @error="audioError = $event"
          class="flex-1 min-h-[280px]"
        />
        <div v-else class="flex-1 min-h-[280px] bg-bg-surface rounded-xl border border-white/5 flex items-center justify-center">
          <div class="text-center">
            <div class="w-10 h-10 rounded-full border-2 border-accent-orange border-t-transparent animate-spin mx-auto mb-3" />
            <p class="text-text-muted text-sm">Getting your room token…</p>
          </div>
        </div>

        <!-- Reaction bar -->
        <ReactionBar ref="reactionBarEl" @reaction="handleReaction" />
      </div>

      <!-- Chat Panel (right / bottom on mobile) -->
      <div class="shrink-0 h-[260px] md:h-auto md:flex-none md:w-[40%] md:min-h-0 border-t md:border-t-0 md:border-l border-white/5 flex flex-col overflow-hidden">
        <ChatPanel
          ref="chatPanelEl"
          :eventId="eventId"
          :isLive="true"
          :messages="chatMessages"
          :token="authToken"
          :myUserId="auth.user?.id"
          class="flex-1"
          @chatEvent="handleChatEvent"
          @reaction="handleIncomingReaction"
        />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter, RouterLink } from 'vue-router'
import { ArrowLeft } from 'lucide-vue-next'
import { useCommentaryStore } from '@/stores/commentary'
import { useAuthStore } from '@/stores/auth'
import { commentaryApi } from '@/api/commentary'
import AudioRoom from '@/components/commentary/AudioRoom.vue'
import ChatPanel from '@/components/commentary/ChatPanel.vue'
import ReactionBar from '@/components/commentary/ReactionBar.vue'
import type { LobbyMessage, ChatEvent } from '@/types'

const route = useRoute()
const router = useRouter()
const store = useCommentaryStore()
const auth = useAuthStore()

const eventId = route.params.id as string
const event = computed(() => store.current)
const livekitToken = computed(() => store.livekitToken)
const livekitUrl = computed(() => store.livekitUrl ?? '')
const myRole = computed(() => store.myRole)
const isHost = computed(() => myRole.value === 'host')
const authToken = computed(() => auth.token ?? '')

const chatMessages = ref<LobbyMessage[]>([])
const roomConnected = ref(false)
const audioError = ref('')
const speakerIds = ref<string[]>([])
const reactionBarEl = ref<InstanceType<typeof ReactionBar> | null>(null)
const chatPanelEl = ref<InstanceType<typeof ChatPanel> | null>(null)

onMounted(async () => {
  // Ensure event detail is loaded
  if (!store.current) {
    await store.fetchDetail(eventId)
  }
  // Restore role + nickname from backend (handles page refresh / new deploy
  // where the Pinia store resets to its defaults).
  await store.checkMe(eventId)
  // Fetch LiveKit token
  try {
    await store.fetchToken(eventId)
  } catch (e: any) {
    audioError.value = e.response?.data?.error ?? 'Could not get room token'
    // Redirect if not authorized
    if (e.response?.status === 403) {
      router.replace(`/commentary/${eventId}`)
    }
  }
})

function handleChatEvent(evt: ChatEvent) {
  if (evt.type === 'message') {
    chatMessages.value.push({
      id: Date.now().toString(),
      event_id: eventId,
      user_id: evt.user_id ?? '',
      nickname: evt.nickname,
      content: evt.content ?? '',
      message_type: 'text',
      created_at: evt.created_at,
    })
  } else if (evt.type === 'speaker_granted' && evt.user_id) {
    if (!speakerIds.value.includes(evt.user_id)) {
      speakerIds.value.push(evt.user_id)
    }
    // If the current user was granted the mic, update role then re-fetch
    // token so AudioRoom reconnects with canPublish=true
    if (evt.user_id === auth.user?.id) {
      store.myRole = 'speaker'
      store.fetchToken(eventId).catch(() => {})
    }
  } else if (evt.type === 'speaker_revoked' && evt.user_id) {
    speakerIds.value = speakerIds.value.filter(id => id !== evt.user_id)
    if (evt.user_id === auth.user?.id) {
      store.myRole = 'listener'
    }
  }
}

async function grantMic(userId: string) {
  try {
    await commentaryApi.updateParticipant(eventId, userId, 'speaker')
    if (!speakerIds.value.includes(userId)) {
      speakerIds.value.push(userId)
    }
  } catch { /* the WS event will sync state if it partially succeeded */ }
}

async function revokeMic(userId: string) {
  try {
    await commentaryApi.updateParticipant(eventId, userId, 'listener')
    speakerIds.value = speakerIds.value.filter(id => id !== userId)
  } catch { /* ignore */ }
}

function handleReaction(emoji: string) {
  // Triggered when THIS user clicks a reaction button
  chatPanelEl.value?.sendRawReaction(emoji)
}

function handleIncomingReaction(emoji: string) {
  // Received from a remote user via WebSocket → animate
  reactionBarEl.value?.spawnFloat(emoji)
}

function raiseHand() {
  const msg = store.myNickname
    ? `✋ ${store.myNickname} wants to speak`
    : '✋ Someone wants to speak'
  chatPanelEl.value?.sendMessage(msg)
}

async function leave() {
  router.push(`/commentary/${eventId}`)
}

async function endRoom() {
  if (!confirm('End this room for everyone?')) return
  try {
    await store.endRoom(eventId)
    router.push(`/commentary/${eventId}`)
  } catch (e: any) {
    alert(e.response?.data?.error ?? 'Failed to end room')
  }
}
</script>
