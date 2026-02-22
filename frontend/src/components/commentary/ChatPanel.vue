<template>
  <div class="flex flex-col h-full bg-bg-surface rounded-xl border border-white/5 overflow-hidden">
    <!-- Header -->
    <div class="px-4 py-3 border-b border-white/5 flex items-center gap-2">
      <span class="text-white font-semibold text-sm">Live Chat</span>
      <span v-if="isLive" class="w-1.5 h-1.5 rounded-full bg-accent-orange animate-pulse" />
      <span v-else class="text-text-muted text-xs">(Replay)</span>
    </div>

    <!-- Messages -->
    <div ref="messagesEl" class="flex-1 overflow-y-auto px-4 py-3 space-y-3 scroll-smooth">
      <div v-if="!messages.length && !isLive" class="text-center text-text-muted text-sm py-8">
        No messages yet.
      </div>
      <div
        v-for="msg in messages"
        :key="msg.id ?? msg.created_at"
        class="flex flex-col"
        :class="{'items-end': msg.user_id === myUserId}"
      >
        <span class="text-accent-orange text-xs font-semibold mb-0.5">{{ msg.nickname }}</span>
        <div
          class="max-w-[85%] px-3 py-2 rounded-xl text-sm text-white"
          :class="msg.user_id === myUserId
            ? 'bg-accent-orange/20 rounded-tr-none'
            : 'bg-bg-elevated rounded-tl-none'"
        >
          {{ msg.content }}
        </div>
        <span class="text-text-muted text-[10px] mt-0.5">
          {{ formatTime(msg.created_at) }}
        </span>
      </div>

      <!-- System events (joined/left) -->
      <div
        v-for="(evt, i) in systemEvents"
        :key="`sys-${i}`"
        class="text-center text-text-muted text-xs py-0.5"
      >
        <span class="font-medium text-white/60">{{ evt.nickname }}</span>
        {{ evt.type === 'joined' ? ' joined the room' : ' left the room' }}
      </div>
    </div>

    <!-- Input (only when live) -->
    <div v-if="isLive" class="px-4 py-3 border-t border-white/5">
      <div class="flex gap-2">
        <input
          v-model="inputText"
          type="text"
          maxlength="500"
          placeholder="Say something…"
          class="flex-1 bg-bg-elevated border border-white/10 rounded-lg px-3 py-2 text-white text-sm
                 placeholder-text-muted focus:outline-none focus:border-accent-orange transition-colors"
          @keyup.enter="send"
        />
        <button
          @click="send"
          :disabled="!inputText.trim()"
          class="px-4 py-2 rounded-lg bg-accent-orange text-white text-sm font-semibold
                 disabled:opacity-40 hover:bg-orange-500 transition-colors"
        >
          Send
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, nextTick, onUnmounted } from 'vue'
import { format } from 'date-fns'
import type { LobbyMessage, ChatEvent } from '@/types'

const props = defineProps<{
  eventId: string
  isLive: boolean
  messages: LobbyMessage[]
  token?: string
  myUserId?: string
}>()

const emit = defineEmits<{
  (e: 'send', content: string): void
  (e: 'reaction', emoji: string): void
  (e: 'chatEvent', event: ChatEvent): void
}>()

const inputText = ref('')
const messagesEl = ref<HTMLElement | null>(null)
const systemEvents = ref<ChatEvent[]>([])

let ws: WebSocket | null = null

// Connect WebSocket when live
watch(() => [props.isLive, props.token], ([isLive, token]) => {
  if (isLive && token) {
    connectWs(token as string)
  } else if (!isLive && ws) {
    ws.close()
    ws = null
  }
}, { immediate: true })

function connectWs(token: string) {
  const wsProtocol = window.location.protocol === 'https:' ? 'wss' : 'ws'
  const wsBase = `${wsProtocol}://${window.location.host}`
  const url = `${wsBase}/ws/commentary/${props.eventId}/chat?token=${encodeURIComponent(token)}`
  ws = new WebSocket(url)

  ws.onmessage = (ev) => {
    try {
      const evt: ChatEvent = JSON.parse(ev.data)
      if (evt.type === 'message') {
        emit('chatEvent', evt)
      } else if (evt.type === 'reaction') {
        emit('reaction', evt.content ?? '')
      } else if (evt.type === 'joined' || evt.type === 'left') {
        systemEvents.value.push(evt)
        // Cap system events in memory
        if (systemEvents.value.length > 50) systemEvents.value.shift()
      } else {
        // speaker_granted / speaker_revoked
        emit('chatEvent', evt)
      }
    } catch { /* ignore malformed */ }
  }
}

function send() {
  const content = inputText.value.trim()
  if (!content) return
  emit('send', content)
  if (ws?.readyState === WebSocket.OPEN) {
    ws.send(JSON.stringify({ type: 'message', content }))
  }
  inputText.value = ''
}

// Auto-scroll to bottom on new messages
watch(() => props.messages.length, () => {
  nextTick(() => {
    if (messagesEl.value) {
      messagesEl.value.scrollTop = messagesEl.value.scrollHeight
    }
  })
})

function formatTime(ts: string) {
  return format(new Date(ts), 'HH:mm')
}

onUnmounted(() => {
  ws?.close()
})

defineExpose({
  ws: () => ws,
  sendRawReaction: (emoji: string) => {
    if (ws?.readyState === WebSocket.OPEN) {
      ws.send(JSON.stringify({ type: 'reaction', content: emoji }))
    }
  },
  sendMessage: (content: string) => {
    if (ws?.readyState === WebSocket.OPEN) {
      ws.send(JSON.stringify({ type: 'message', content }))
      return true
    }
    return false
  },
})
</script>
