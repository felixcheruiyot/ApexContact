<template>
  <div class="max-w-4xl mx-auto px-4 sm:px-6 py-10">
    <!-- Loading -->
    <div v-if="store.loading" class="space-y-4 animate-pulse">
      <div class="h-64 bg-bg-elevated rounded-xl" />
      <div class="h-8 bg-bg-elevated rounded w-1/2" />
      <div class="h-4 bg-bg-elevated rounded w-1/3" />
    </div>

    <div v-else-if="!event">
      <p class="text-text-muted text-center py-20">Lobby not found.</p>
    </div>

    <template v-else>
      <!-- Thumbnail -->
      <div class="relative rounded-xl overflow-hidden aspect-video mb-8">
        <img
          :src="event.thumbnail_url || '/placeholder-event.jpg'"
          :alt="event.title"
          class="w-full h-full object-cover"
        />
        <div class="absolute inset-0 bg-gradient-to-t from-black/80 via-transparent to-transparent" />
        <div class="absolute bottom-6 left-6 right-6">
          <div class="flex items-center gap-2 mb-2">
            <span v-if="event.status === 'live'" class="flex items-center gap-1.5 px-3 py-1 rounded-full bg-accent-orange text-white text-xs font-bold uppercase">
              <span class="w-1.5 h-1.5 bg-white rounded-full animate-pulse" />
              Live Now
            </span>
            <span class="bg-black/50 backdrop-blur-sm text-white text-xs px-2 py-1 rounded-full uppercase tracking-wider">
              {{ event.sport_type }}
            </span>
            <span class="bg-accent-orange/20 text-accent-orange text-xs px-2 py-1 rounded-full font-medium">
              🎙 Commentary
            </span>
          </div>
          <h1 class="text-white text-2xl md:text-3xl font-bold font-display leading-tight">
            {{ event.title }}
          </h1>
          <p v-if="event.teaser_hook" class="text-white/80 mt-1 text-sm">{{ event.teaser_hook }}</p>
        </div>
      </div>

      <!-- Meta row -->
      <div class="flex flex-wrap items-center gap-6 mb-6 text-sm text-text-muted">
        <span class="flex items-center gap-1.5">
          📅 {{ formattedDate }}
        </span>
        <span class="flex items-center gap-1.5">
          👥 {{ store.participantCount }} joining
        </span>
        <span class="flex items-center gap-1.5 font-bold" :class="event.price === 0 ? 'text-success' : 'text-accent-orange'">
          {{ event.price === 0 ? 'Free' : `${event.currency} ${event.price.toLocaleString()}` }}
        </span>
        <span v-if="event.status === 'completed'" class="text-text-muted">Ended</span>
      </div>

      <!-- Description -->
      <p v-if="event.description" class="text-text-muted leading-relaxed mb-8">
        {{ event.description }}
      </p>

      <!-- CTA -->
      <div class="mb-12">
        <!-- Already joined & live -->
        <RouterLink
          v-if="hasJoined && event.status === 'live'"
          :to="`/commentary/${event.id}/room`"
          class="inline-flex items-center gap-2 px-6 py-3 rounded-xl bg-accent-orange text-white font-bold
                 hover:bg-orange-500 transition-colors text-base"
        >
          🎙 Join Live Room
        </RouterLink>

        <!-- Already joined, not live yet -->
        <div v-else-if="hasJoined && event.status === 'scheduled'" class="flex items-center gap-4">
          <span class="text-success font-medium">✓ You're in! Waiting for the host to start.</span>
        </div>

        <!-- Completed — replay only -->
        <div v-else-if="event.status === 'completed'" class="text-text-muted">
          This lobby has ended. See the chat replay below.
        </div>

        <!-- Join button -->
        <button
          v-else-if="event.status === 'scheduled' || event.status === 'live'"
          @click="openJoinModal"
          class="inline-flex items-center gap-2 px-6 py-3 rounded-xl font-bold text-base transition-colors"
          :class="event.price === 0
            ? 'bg-success text-white hover:bg-green-600'
            : 'bg-accent-orange text-white hover:bg-orange-500'"
        >
          {{ event.price === 0 ? '🎙 Join Free' : '🎫 Buy Access' }}
        </button>
      </div>

      <!-- Chat replay (completed) -->
      <div v-if="event.status === 'completed'" class="mt-8">
        <h2 class="text-white font-bold text-lg mb-4">Chat Replay</h2>
        <div v-if="store.loading" class="space-y-2">
          <div v-for="i in 5" :key="i" class="h-10 bg-bg-elevated rounded animate-pulse" />
        </div>
        <div v-else-if="store.messages.length" class="space-y-3">
          <div
            v-for="msg in store.messages"
            :key="msg.id"
            class="flex gap-3 items-start"
          >
            <span class="text-accent-orange text-sm font-semibold shrink-0 min-w-[80px]">{{ msg.nickname }}</span>
            <span class="text-white text-sm">{{ msg.content }}</span>
            <span class="text-text-muted text-xs shrink-0 ml-auto">{{ formatTime(msg.created_at) }}</span>
          </div>
        </div>
        <p v-else class="text-text-muted text-sm">No chat messages for this lobby.</p>
      </div>
    </template>

    <!-- Nickname Modal -->
    <NicknameModal
      v-if="showNicknameModal && event"
      :eventId="event.id"
      :isFree="event.price === 0"
      @close="showNicknameModal = false"
      @confirm="handleJoin"
    />

    <!-- Payment Modal (paid lobbies) -->
    <MpesaModal
      v-if="showPaymentModal && event"
      :event="event"
      @close="afterPaymentModalClose"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter, RouterLink } from 'vue-router'
import { format } from 'date-fns'
import { useCommentaryStore } from '@/stores/commentary'
import { useAuthStore } from '@/stores/auth'
import NicknameModal from '@/components/commentary/NicknameModal.vue'
import MpesaModal from '@/components/payment/MpesaModal.vue'

const route = useRoute()
const router = useRouter()
const store = useCommentaryStore()
const auth = useAuthStore()

const showNicknameModal = ref(false)
const showPaymentModal = ref(false)
const pendingNickname = ref('')
const hasJoined = ref(false)

const event = computed(() => store.current)

const formattedDate = computed(() =>
  event.value ? format(new Date(event.value.scheduled_at), "MMM d, yyyy · h:mm a") : ''
)

function formatTime(ts: string) {
  return format(new Date(ts), 'HH:mm')
}

onMounted(async () => {
  await store.fetchDetail(route.params.id as string)
  if (event.value?.status === 'completed') {
    store.fetchMessages(event.value.id)
  }
})

function openJoinModal() {
  if (!auth.isAuthenticated) {
    router.push({ name: 'login', query: { redirect: route.fullPath } })
    return
  }
  showNicknameModal.value = true
}

async function handleJoin(nickname: string) {
  if (!event.value) return
  showNicknameModal.value = false

  if (event.value.price > 0) {
    pendingNickname.value = nickname
    showPaymentModal.value = true
    return
  }

  try {
    await store.joinLobby(event.value.id, nickname)
    hasJoined.value = true
    if (event.value.status === 'live') {
      router.push(`/commentary/${event.value.id}/room`)
    }
  } catch (e: any) {
    alert(e.response?.data?.error ?? 'Failed to join lobby')
  }
}

// Called when the MpesaModal closes (either after payment success or cancel)
async function afterPaymentModalClose() {
  showPaymentModal.value = false
  if (!event.value || !pendingNickname.value) return
  // Attempt to join — will succeed if payment completed and subscription was created
  try {
    await store.joinLobby(event.value.id, pendingNickname.value)
    hasJoined.value = true
    if (event.value.status === 'live') {
      router.push(`/commentary/${event.value.id}/room`)
    }
  } catch {
    // Payment may not have completed — user can try again
  }
}
</script>
