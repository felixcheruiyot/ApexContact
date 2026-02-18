<template>
  <div class="relative w-full aspect-video bg-black rounded-xl overflow-hidden">
    <video ref="videoEl" class="w-full h-full" controls playsinline />
    <!-- Loading overlay -->
    <div v-if="loading"
      class="absolute inset-0 flex items-center justify-center bg-black/70">
      <div class="flex flex-col items-center gap-3">
        <div class="w-10 h-10 border-4 border-accent-red border-t-transparent rounded-full animate-spin" />
        <p class="text-text-muted text-sm">Loading stream...</p>
      </div>
    </div>
    <!-- Error overlay -->
    <div v-if="playerError"
      class="absolute inset-0 flex items-center justify-center bg-black/80">
      <div class="text-center max-w-sm px-6">
        <p class="text-status-error text-lg font-semibold mb-2">Stream unavailable</p>
        <p class="text-text-muted text-sm">{{ playerError }}</p>
        <button @click="retry" class="btn-primary mt-4 text-sm py-2 px-4">Retry</button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch } from 'vue'
import Hls from 'hls.js'

const props = defineProps<{ src: string }>()

const videoEl = ref<HTMLVideoElement | null>(null)
const loading = ref(true)
const playerError = ref<string | null>(null)
let hls: Hls | null = null

function initPlayer(src: string) {
  if (!videoEl.value) return
  loading.value = true
  playerError.value = null

  if (Hls.isSupported()) {
    hls = new Hls({ enableWorker: true, lowLatencyMode: true })
    hls.loadSource(src)
    hls.attachMedia(videoEl.value)

    hls.on(Hls.Events.MANIFEST_PARSED, () => {
      loading.value = false
      videoEl.value?.play()
    })

    hls.on(Hls.Events.ERROR, (_event, data) => {
      if (data.fatal) {
        loading.value = false
        playerError.value = 'The stream has encountered an error. Please try again.'
      }
    })
  } else if (videoEl.value.canPlayType('application/vnd.apple.mpegurl')) {
    // Native HLS (Safari)
    videoEl.value.src = src
    videoEl.value.addEventListener('loadedmetadata', () => {
      loading.value = false
      videoEl.value?.play()
    })
  } else {
    playerError.value = 'Your browser does not support HLS streaming.'
    loading.value = false
  }
}

function destroyPlayer() {
  hls?.destroy()
  hls = null
}

function retry() {
  destroyPlayer()
  initPlayer(props.src)
}

onMounted(() => initPlayer(props.src))
onUnmounted(() => destroyPlayer())
watch(() => props.src, (newSrc) => { destroyPlayer(); initPlayer(newSrc) })
</script>
