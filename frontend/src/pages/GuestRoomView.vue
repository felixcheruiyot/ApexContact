<template>
  <div class="min-h-[60vh]">
    <!-- Loading -->
    <div v-if="loading" class="flex items-center justify-center min-h-[60vh]">
      <div class="flex flex-col items-center gap-3">
        <div class="w-10 h-10 border-4 border-accent-orange border-t-transparent rounded-full animate-spin" />
        <p class="text-text-muted text-sm">Joining room…</p>
      </div>
    </div>

    <!-- Error -->
    <div v-else-if="error" class="flex items-center justify-center min-h-[60vh] px-4">
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

    <!-- Room — same layout as Watch.vue for audio/video events -->
    <template v-else-if="roomData">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6">

        <!-- Live room embed -->
        <div class="mb-6 bg-bg-elevated border border-white/10 rounded-xl overflow-hidden">
          <!-- Pre-join state -->
          <div v-if="!joined" class="p-10 flex flex-col items-center gap-5 text-center">
            <div class="w-16 h-16 rounded-full bg-accent-orange/15 flex items-center justify-center">
              <component :is="roomData.event_type === 'audio' ? Mic : Video" class="w-8 h-8 text-accent-orange" />
            </div>
            <div>
              <p class="text-white font-bold text-xl mb-1">You're joining a live session</p>
              <p class="text-text-muted text-sm">
                You'll join as a listener.
                {{ roomData.event_type === 'audio_video' ? 'You can see and hear the host.' : 'You can hear the host.' }}
              </p>
            </div>
            <button
              @click="joined = true"
              class="px-8 py-3 rounded-lg bg-accent-orange hover:bg-orange-500
                     text-white font-bold text-sm transition-all active:scale-95"
            >
              <component :is="roomData.event_type === 'audio' ? Mic : Video" class="w-4 h-4 inline-block mr-2" />
              Join as listener
            </button>
          </div>

          <!-- AudioRoom handles the LiveKit connection -->
          <AudioRoom
            v-else
            :token="roomData.token"
            :livekit-url="roomData.livekit_url"
            my-role="listener"
            :is-host="false"
          />
        </div>

        <!-- Info section below room -->
        <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
          <div class="lg:col-span-2">
            <div class="flex items-center gap-3 mb-3">
              <span class="w-2 h-2 rounded-full bg-accent-orange animate-pulse" />
              <span class="text-white font-semibold text-sm uppercase tracking-wider">
                Live {{ roomData.event_type === 'audio' ? 'Audio Room' : 'Audio + Video Room' }}
              </span>
            </div>
            <h1 class="text-white font-bold text-2xl mb-3">Guest Session</h1>
            <p class="text-text-muted leading-relaxed">
              This is a guest test session. Sign up to host unlimited sessions with ticket sales and M-Pesa payments.
            </p>
          </div>

          <!-- Sidebar -->
          <div class="card p-5 h-fit">
            <h3 class="text-white font-semibold mb-4">Session Info</h3>
            <dl class="space-y-3">
              <div>
                <dt class="text-text-muted text-xs uppercase tracking-wider mb-1">Status</dt>
                <dd>
                  <span class="inline-flex items-center gap-1.5 badge-live">
                    <span class="w-1.5 h-1.5 bg-white rounded-full animate-pulse" />
                    Live
                  </span>
                </dd>
              </div>
              <div>
                <dt class="text-text-muted text-xs uppercase tracking-wider mb-1">Format</dt>
                <dd class="text-white text-sm capitalize">
                  {{ roomData.event_type === 'audio' ? 'Audio Only' : 'Audio + Video' }}
                </dd>
              </div>
              <div>
                <dt class="text-text-muted text-xs uppercase tracking-wider mb-1">Host</dt>
                <dd class="text-white text-sm">Anonymous</dd>
              </div>
            </dl>
            <div class="mt-4 pt-4 border-t border-white/5 space-y-3">
              <p class="text-text-muted text-xs">Want to host your own session?</p>
              <RouterLink to="/try" class="btn-primary w-full text-center text-sm block py-2">
                Start a free stream
              </RouterLink>
              <RouterLink to="/register" class="btn-ghost w-full text-center text-sm block py-2">
                Create account
              </RouterLink>
            </div>
          </div>
        </div>

      </div>
    </template>
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
