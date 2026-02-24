<template>
  <div class="min-h-screen bg-black flex flex-col">

    <!-- Minimal header -->
    <header class="absolute top-0 left-0 right-0 z-10 px-6 py-4 flex items-center justify-between">
      <RouterLink to="/" class="font-display text-xl tracking-widest text-white/80 hover:text-white transition-colors">
        LIVE <span class="text-accent-red">·</span> STREAMIFY
      </RouterLink>
      <RouterLink
        to="/register"
        class="px-4 py-2 rounded-lg bg-accent-red hover:bg-accent-red-hover
               text-white text-sm font-semibold transition-all"
      >
        Start your own stream
      </RouterLink>
    </header>

    <!-- Player area -->
    <div class="flex-1 flex flex-col">
      <!-- Loading state -->
      <div v-if="loading" class="flex-1 flex items-center justify-center">
        <div class="flex flex-col items-center gap-3">
          <div class="w-10 h-10 border-4 border-accent-red border-t-transparent rounded-full animate-spin" />
          <p class="text-text-muted text-sm">Loading stream…</p>
        </div>
      </div>

      <!-- Error: stream not found or expired -->
      <div v-else-if="error" class="flex-1 flex items-center justify-center px-4">
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

      <!-- Player -->
      <div v-else-if="hlsUrl" class="flex flex-col flex-1">
        <!-- Full-width player -->
        <div class="w-full">
          <VideoPlayer :src="hlsUrl" />
        </div>

        <!-- Stream info bar -->
        <div class="bg-bg border-t border-white/5 px-6 py-4 flex items-center justify-between flex-wrap gap-4">
          <div class="flex items-center gap-2">
            <span class="w-2 h-2 rounded-full bg-accent-red animate-pulse" />
            <span class="text-white text-sm font-medium">Live Test Stream</span>
            <span class="text-text-muted text-xs ml-1">(guest session)</span>
          </div>
          <RouterLink
            to="/try"
            class="inline-flex items-center gap-2 px-4 py-2 rounded-lg bg-accent-red
                   hover:bg-accent-red-hover text-white text-sm font-semibold transition-all"
          >
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M15 10l4.553-2.276A1 1 0 0121 8.677v6.646a1 1 0 01-1.447.894L15 14M5 18h8a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v8a2 2 0 002 2z" />
            </svg>
            Start your own stream
          </RouterLink>
        </div>
      </div>
    </div>

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
