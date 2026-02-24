<template>
  <div class="min-h-screen bg-bg flex flex-col">

    <!-- Minimal header -->
    <header class="border-b border-white/5 px-6 py-4 flex items-center justify-between">
      <RouterLink to="/" class="font-display text-xl tracking-widest text-white/80 hover:text-white transition-colors">
        LIVE <span class="text-accent-red">·</span> STREAMIFY
      </RouterLink>
      <RouterLink
        to="/register"
        class="px-4 py-2 rounded-lg bg-accent-red hover:bg-accent-red-hover
               text-white text-sm font-semibold transition-all"
      >
        Start your own
      </RouterLink>
    </header>

    <!-- Loading -->
    <div v-if="loading" class="flex-1 flex items-center justify-center">
      <div class="flex flex-col items-center gap-3">
        <div class="w-10 h-10 border-4 border-accent-orange border-t-transparent rounded-full animate-spin" />
        <p class="text-text-muted text-sm">Joining room…</p>
      </div>
    </div>

    <!-- Error -->
    <div v-else-if="error" class="flex-1 flex items-center justify-center px-4">
      <div class="text-center max-w-sm space-y-4">
        <div class="w-16 h-16 rounded-full bg-white/5 flex items-center justify-center mx-auto">
          <Mic class="w-8 h-8 text-text-muted" />
        </div>
        <div>
          <p class="text-white font-semibold text-lg mb-1">Room not available</p>
          <p class="text-text-muted text-sm">{{ error }}</p>
        </div>
        <RouterLink
          to="/try"
          class="inline-block px-5 py-2.5 rounded-lg border border-white/20 hover:border-white/40
                 text-white text-sm font-medium transition-all hover:bg-white/5"
        >
          Start your own test stream
        </RouterLink>
      </div>
    </div>

    <!-- Room ready — show in commentary/AudioRoom component -->
    <div v-else-if="roomData" class="flex-1 flex flex-col items-center justify-center px-4 py-10">
      <div class="w-full max-w-xl space-y-6">
        <div class="flex items-center gap-3">
          <span class="w-2 h-2 rounded-full bg-accent-orange animate-pulse" />
          <span class="text-white font-semibold text-sm uppercase tracking-wider">
            Guest {{ roomData.event_type === 'audio' ? 'Audio Room' : 'Audio + Video Room' }}
          </span>
        </div>

        <div class="bg-bg-elevated border border-white/10 rounded-xl p-8 text-center space-y-4">
          <div class="w-14 h-14 rounded-full bg-accent-orange/15 flex items-center justify-center mx-auto">
            <component :is="roomData.event_type === 'audio' ? Mic : Video" class="w-7 h-7 text-accent-orange" />
          </div>
          <div>
            <p class="text-white font-bold text-lg mb-1">You're joining a live session</p>
            <p class="text-text-muted text-sm leading-relaxed">
              You'll join as a listener.
              {{ roomData.event_type === 'audio_video' ? 'You can see and hear the host.' : 'You can hear the host.' }}
            </p>
          </div>

          <!-- The AudioRoom component handles the actual LiveKit connection -->
          <AudioRoom
            v-if="joined"
            :token="roomData.token"
            :livekit-url="roomData.livekit_url"
            my-role="listener"
            :is-host="false"
          />

          <button
            v-else
            @click="joined = true"
            class="w-full px-6 py-3.5 rounded-lg bg-accent-orange hover:bg-orange-500
                   text-white font-bold text-sm transition-all active:scale-95"
          >
            <component :is="roomData.event_type === 'audio' ? Mic : Video" class="w-4 h-4 inline-block mr-2" />
            Join as listener
          </button>
        </div>

        <p class="text-center text-text-muted text-xs">
          This is a guest session — it expires automatically.
          <RouterLink to="/try" class="text-white underline hover:no-underline ml-1">Start your own</RouterLink>
        </p>
      </div>
    </div>

  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { RouterLink, useRoute } from 'vue-router'
import { Mic, Video } from 'lucide-vue-next'
import { getGuestRoomToken } from '@/api/stream'
import AudioRoom from '@/components/commentary/AudioRoom.vue'

const route = useRoute()
const loading = ref(true)
const error = ref<string | null>(null)
const joined = ref(false)
const roomData = ref<{
  token: string
  room_name: string
  event_type: string
  livekit_url: string
} | null>(null)

onMounted(async () => {
  try {
    const guestId = route.params.guestId as string
    roomData.value = await getGuestRoomToken(guestId)
  } catch (err: any) {
    const msg = err?.response?.data?.error
    error.value = msg || 'This room link has expired or the host has ended the session.'
  } finally {
    loading.value = false
  }
})
</script>
