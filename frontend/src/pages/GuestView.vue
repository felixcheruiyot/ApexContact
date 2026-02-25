<template>
  <div class="min-h-[60vh]">
    <!-- Loading state -->
    <div v-if="loading" class="flex items-center justify-center min-h-[60vh]">
      <div class="flex flex-col items-center gap-3">
        <div class="w-10 h-10 border-4 border-accent-red border-t-transparent rounded-full animate-spin" />
        <p class="text-text-muted text-sm">Loading stream…</p>
      </div>
    </div>

    <!-- Error: stream not found or expired -->
    <div v-else-if="error" class="flex items-center justify-center min-h-[60vh] px-4">
      <div class="text-center max-w-sm space-y-4">
        <div class="w-16 h-16 rounded-full bg-white/5 flex items-center justify-center mx-auto">
          <svg class="w-8 h-8 text-text-muted" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5"
              d="M15 10l4.553-2.276A1 1 0 0121 8.677v6.646a1 1 0 01-1.447.894L15 14M5 18h8a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v8a2 2 0 002 2z" />
          </svg>
        </div>
        <div>
          <p class="text-white font-semibold text-lg mb-1">Stream not available</p>
          <p class="text-text-muted text-sm">{{ error }}</p>
        </div>
        <button
          @click="fetchStream"
          class="px-5 py-2.5 rounded-lg border border-white/20 hover:border-white/40
                 text-white text-sm font-medium transition-all hover:bg-white/5"
        >
          Try again
        </button>
      </div>
    </div>

    <!-- Player — same layout as Watch.vue -->
    <template v-else-if="hlsUrl">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6">
        <VideoPlayer :src="hlsUrl" class="mb-6" />

        <!-- Event info (YouTube-style below player) -->
        <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
          <div class="lg:col-span-2">
            <div class="flex items-start gap-4 mb-4">
              <div>
                <div class="flex items-center gap-3 mb-2">
                  <span class="badge-live">
                    <span class="w-1.5 h-1.5 bg-white rounded-full animate-pulse" />
                    Live
                  </span>
                </div>
                <h1 class="text-white font-bold text-2xl">Live Test Stream</h1>
              </div>
            </div>
            <p class="text-text-muted leading-relaxed">
              This is a guest test stream. Create a free account to host unlimited sessions and charge your audience via M-Pesa.
            </p>
          </div>

          <!-- Sidebar -->
          <div class="card p-5 h-fit">
            <h3 class="text-white font-semibold mb-4">Stream Info</h3>
            <dl class="space-y-3">
              <div>
                <dt class="text-text-muted text-xs uppercase tracking-wider mb-1">Status</dt>
                <dd><span class="badge-live">Live</span></dd>
              </div>
              <div>
                <dt class="text-text-muted text-xs uppercase tracking-wider mb-1">Host</dt>
                <dd class="text-white text-sm">Anonymous</dd>
              </div>
            </dl>
            <div class="mt-4 pt-4 border-t border-white/5 space-y-3">
              <p class="text-text-muted text-xs">Want to host your own stream?</p>
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
import client from '@/api/client'
import VideoPlayer from '@/components/player/VideoPlayer.vue'


const route = useRoute()
const hlsUrl = ref<string | null>(null)
const loading = ref(true)
const error = ref<string | null>(null)

async function fetchStream() {
  loading.value = true
  error.value = null
  try {
    const { data } = await client.get(`/stream/guest/${route.params.guestId}`)
    hlsUrl.value = data.data.hls_url
  } catch (err: any) {
    const msg = err?.response?.data?.error
    error.value = msg || 'This stream link has expired or the broadcaster hasn\'t started yet.'
  } finally {
    loading.value = false
  }
}

onMounted(fetchStream)
</script>
