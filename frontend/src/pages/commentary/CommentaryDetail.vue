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
          :src="eventImage(event.thumbnail_url)"
          :alt="event.title"
          class="w-full h-full object-cover"
          @error="onImageError"
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
            <span class="bg-accent-orange/20 text-accent-orange text-xs px-2 py-1 rounded-full font-medium flex items-center gap-1">
              <Mic class="w-3 h-3" /> Live Room
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
        <span class="flex items-center gap-1.5"><Calendar class="w-4 h-4" /> {{ formattedDate }}</span>
        <span class="flex items-center gap-1.5"><Users class="w-4 h-4" /> {{ store.participantCount }} joining</span>
        <span class="font-bold" :class="event.price === 0 ? 'text-success' : 'text-accent-orange'">
          {{ event.price === 0 ? 'Free' : `${event.currency} ${event.price.toLocaleString()}` }}
        </span>
        <span v-if="event.status === 'completed'" class="text-text-muted">Ended</span>
      </div>

      <!-- Description -->
      <p v-if="event.description" class="text-text-muted leading-relaxed mb-8">{{ event.description }}</p>

      <!-- ─── CTA block ──────────────────────────────────────────────────────── -->
      <div class="flex flex-wrap items-center gap-3 mb-12">

        <!-- ── CREATOR controls ── -->
        <template v-if="isCreator">
          <!-- Start Room -->
          <button
            v-if="event.status === 'scheduled'"
            @click="handleStart"
            :disabled="starting"
            class="inline-flex items-center gap-2 px-6 py-3 rounded-xl bg-accent-orange text-white
                   font-bold text-base hover:bg-orange-500 transition-colors disabled:opacity-50"
          >
            <Play v-if="!starting" class="w-4 h-4" />{{ starting ? 'Starting…' : 'Start Room' }}
          </button>

          <!-- Go to live room + End Room -->
          <template v-else-if="event.status === 'live'">
            <RouterLink
              :to="`/commentary/${event.id}/room`"
              class="inline-flex items-center gap-2 px-6 py-3 rounded-xl bg-accent-orange text-white
                     font-bold text-base hover:bg-orange-500 transition-colors"
            >
              <Mic class="w-4 h-4" /> Enter Your Room
            </RouterLink>
            <button
              @click="handleEnd"
              :disabled="ending"
              class="px-4 py-3 rounded-xl border border-red-500/40 text-red-400 text-sm font-medium
                     hover:bg-red-500/10 transition-colors disabled:opacity-50"
            >
              {{ ending ? 'Ending…' : 'End Room' }}
            </button>
          </template>

          <span v-else-if="event.status === 'completed'" class="text-text-muted text-sm">
            This lobby has ended.
          </span>
        </template>

        <!-- ── VIEWER / LISTENER controls ── -->
        <template v-else>
          <!-- Completed — no action -->
          <span v-if="event.status === 'completed'" class="text-text-muted text-sm">
            This lobby has ended. See the chat replay below.
          </span>

          <!-- Already joined & live → enter room -->
          <RouterLink
            v-else-if="hasJoined && event.status === 'live'"
            :to="`/commentary/${event.id}/room`"
            class="inline-flex items-center gap-2 px-6 py-3 rounded-xl bg-accent-orange text-white
                   font-bold text-base hover:bg-orange-500 transition-colors"
          >
            <Mic class="w-4 h-4" /> Enter Room
          </RouterLink>

          <!-- Already joined, waiting for host -->
          <div v-else-if="hasJoined && event.status === 'scheduled'"
            class="flex items-center gap-3 px-4 py-3 rounded-xl bg-success/10 border border-success/20">
            <Check class="w-5 h-5 text-success" />
            <div>
              <p class="text-success font-semibold text-sm">You're in as <strong>{{ store.myNickname }}</strong></p>
              <p class="text-text-muted text-xs">Waiting for the host to start the room.</p>
            </div>
          </div>

          <!-- Not yet joined — free -->
          <button
            v-else-if="event.price === 0 && (event.status === 'scheduled' || event.status === 'live')"
            @click="openJoinModal"
            class="inline-flex items-center gap-2 px-6 py-3 rounded-xl bg-success text-white
                   font-bold text-base hover:bg-green-600 transition-colors"
          >
            <Mic class="w-4 h-4" /> Join Free
          </button>

          <!-- Not yet joined — paid -->
          <button
            v-else-if="event.price > 0 && (event.status === 'scheduled' || event.status === 'live')"
            @click="openJoinModal"
            class="inline-flex items-center gap-2 px-6 py-3 rounded-xl bg-accent-orange text-white
                   font-bold text-base hover:bg-orange-500 transition-colors"
          >
            <Ticket class="w-4 h-4" /> Buy Access · {{ event.currency }} {{ event.price.toLocaleString() }}
          </button>
        </template>

      </div>
      <!-- ─────────────────────────────────────────────────────────────────────── -->

      <!-- Chat replay (completed) -->
      <div v-if="event.status === 'completed'" class="mt-8">
        <h2 class="text-white font-bold text-lg mb-4">Chat Replay</h2>
        <div v-if="store.loading" class="space-y-2">
          <div v-for="i in 5" :key="i" class="h-10 bg-bg-elevated rounded animate-pulse" />
        </div>
        <div v-else-if="store.messages.length" class="space-y-3">
          <div v-for="msg in store.messages" :key="msg.id" class="flex gap-3 items-start">
            <span class="text-accent-orange text-sm font-semibold shrink-0 min-w-[80px]">{{ msg.nickname }}</span>
            <span class="text-white text-sm">{{ msg.content }}</span>
            <span class="text-text-muted text-xs shrink-0 ml-auto">{{ formatTime(msg.created_at) }}</span>
          </div>
        </div>
        <p v-else class="text-text-muted text-sm">No chat messages recorded for this lobby.</p>
      </div>
    </template>

    <!-- Nickname Modal (free join) -->
    <NicknameModal
      v-if="showNicknameModal && event"
      :eventId="event.id"
      :isFree="event.price === 0"
      @close="showNicknameModal = false"
      @confirm="handleJoin"
    />

    <!-- Payment Modal (paid join) -->
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
import { Mic, Calendar, Users, Play, Check, Ticket } from 'lucide-vue-next'
import { format } from 'date-fns'
import { useCommentaryStore } from '@/stores/commentary'
import { useAuthStore } from '@/stores/auth'
import NicknameModal from '@/components/commentary/NicknameModal.vue'
import MpesaModal from '@/components/payment/MpesaModal.vue'
import { eventImage, onImageError } from '@/utils/eventImage'

const route = useRoute()
const router = useRouter()
const store = useCommentaryStore()
const auth = useAuthStore()

const showNicknameModal = ref(false)
const showPaymentModal = ref(false)
const pendingNickname = ref('')
const hasJoined = ref(false)
const starting = ref(false)
const ending = ref(false)

const event = computed(() => store.current)

// True when the logged-in user is the creator of this lobby
const isCreator = computed(() =>
  auth.isAuthenticated && !!event.value && auth.user?.id === event.value.promoter_id
)

const formattedDate = computed(() =>
  event.value ? format(new Date(event.value.scheduled_at), 'MMM d, yyyy · h:mm a') : ''
)

function formatTime(ts: string) {
  return format(new Date(ts), 'HH:mm')
}

onMounted(async () => {
  await store.fetchDetail(route.params.id as string)

  // Check whether the current user has already joined
  if (auth.isAuthenticated) {
    hasJoined.value = await store.checkMe(route.params.id as string)
  }

  if (event.value?.status === 'completed') {
    store.fetchMessages(event.value.id)
  }
})

// ── Creator actions ────────────────────────────────────────────────────────────

async function handleStart() {
  if (!event.value) return
  starting.value = true
  try {
    await store.startRoom(event.value.id)
    // Auto-join as host after starting
    if (!hasJoined.value) {
      const nickname = auth.user?.full_name?.split(' ')[0] ?? 'Host'
      await store.joinLobby(event.value.id, nickname)
      hasJoined.value = true
    }
  } catch (e: any) {
    alert(e.response?.data?.error ?? 'Failed to start room')
  } finally {
    starting.value = false
  }
}

async function handleEnd() {
  if (!event.value || !confirm('End this room for everyone?')) return
  ending.value = true
  try {
    await store.endRoom(event.value.id)
  } catch (e: any) {
    alert(e.response?.data?.error ?? 'Failed to end room')
  } finally {
    ending.value = false
  }
}

// ── Viewer join flow ───────────────────────────────────────────────────────────

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

  // Paid lobby → open M-Pesa first, then join after callback
  if (event.value.price > 0) {
    pendingNickname.value = nickname
    showPaymentModal.value = true
    return
  }

  // Free lobby → join directly
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

async function afterPaymentModalClose() {
  showPaymentModal.value = false
  if (!event.value || !pendingNickname.value) return
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
