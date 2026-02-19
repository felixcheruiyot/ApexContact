<template>
  <div class="min-h-screen bg-bg-page flex flex-col">
    <!-- Top bar -->
    <div class="flex items-center justify-between px-4 py-3 bg-bg-surface border-b border-white/5">
      <RouterLink :to="`/commentary/${eventId}`" class="text-text-muted hover:text-white transition-colors">
        ← Back
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
    <div class="flex-1 flex flex-col md:flex-row gap-0 overflow-hidden" style="height: calc(100vh - 57px)">
      <!-- Audio Room (left / top on mobile) -->
      <div class="flex-1 md:flex-none md:w-[60%] p-4 flex flex-col gap-4 overflow-y-auto">
        <AudioRoom
          v-if="livekitToken && livekitUrl"
          :livekitUrl="livekitUrl"
          :token="livekitToken"
          :myRole="myRole"
          @leave="leave"
          @raiseHand="raiseHand"
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
      <div class="md:w-[40%] border-t md:border-t-0 md:border-l border-white/5 flex flex-col overflow-hidden" style="min-height: 300px; max-height: 100%">
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
import { useCommentaryStore } from '@/stores/commentary'
import { useAuthStore } from '@/stores/auth'
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
const reactionBarEl = ref<InstanceType<typeof ReactionBar> | null>(null)
const chatPanelEl = ref<InstanceType<typeof ChatPanel> | null>(null)

onMounted(async () => {
  // Ensure event detail is loaded
  if (!store.current) {
    await store.fetchDetail(eventId)
  }
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
    // Push as a LobbyMessage shape for display
    chatMessages.value.push({
      id: Date.now().toString(),
      event_id: eventId,
      user_id: evt.user_id ?? '',
      nickname: evt.nickname,
      content: evt.content ?? '',
      message_type: 'text',
      created_at: evt.created_at,
    })
  }
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
  // TODO: could emit a chat message like "✋ [nickname] wants to speak"
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
