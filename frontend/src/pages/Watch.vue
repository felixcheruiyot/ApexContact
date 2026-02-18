<template>
  <div class="min-h-screen bg-bg">
    <!-- Loading -->
    <div v-if="loading" class="flex items-center justify-center min-h-screen">
      <div class="text-center">
        <div class="w-12 h-12 border-4 border-accent-red border-t-transparent rounded-full animate-spin mx-auto mb-4" />
        <p class="text-text-muted">Verifying your ticket...</p>
      </div>
    </div>

    <!-- Error / Access denied -->
    <div v-else-if="error" class="flex items-center justify-center min-h-screen">
      <div class="text-center max-w-sm px-6">
        <p class="text-5xl mb-4">🔒</p>
        <h2 class="text-white font-bold text-xl mb-2">Access Denied</h2>
        <p class="text-text-muted text-sm mb-6">{{ error }}</p>
        <div class="flex flex-col gap-3">
          <RouterLink
            :to="{ name: 'event-detail', params: { id: route.params.eventId } }"
            class="btn-primary"
          >
            Back to Event
          </RouterLink>
          <RouterLink to="/" class="text-text-muted text-sm hover:text-white transition-colors">
            Browse All Events
          </RouterLink>
        </div>
      </div>
    </div>

    <!-- Player -->
    <template v-else-if="hlsUrl && event">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6">
        <VideoPlayer :src="hlsUrl" class="mb-6" />

        <!-- Event info (YouTube-style below player) -->
        <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
          <div class="lg:col-span-2">
            <div class="flex items-start justify-between gap-4 mb-4">
              <div>
                <div class="flex items-center gap-3 mb-2">
                  <span class="badge-live">
                    <span class="w-1.5 h-1.5 bg-white rounded-full animate-pulse" />
                    Live
                  </span>
                  <span class="text-text-muted text-sm capitalize">{{ event.sport_type }}</span>
                </div>
                <h1 class="text-white font-bold text-2xl">{{ event.title }}</h1>
              </div>
            </div>
            <p class="text-text-muted leading-relaxed">{{ event.description }}</p>
          </div>

          <!-- Stream info sidebar -->
          <div class="card p-5 h-fit">
            <h3 class="text-white font-semibold mb-4">Stream Info</h3>
            <dl class="space-y-3">
              <div>
                <dt class="text-text-muted text-xs uppercase tracking-wider mb-1">Status</dt>
                <dd><span class="badge-live">Live</span></dd>
              </div>
              <div>
                <dt class="text-text-muted text-xs uppercase tracking-wider mb-1">Sport</dt>
                <dd class="text-white text-sm capitalize">{{ event.sport_type }}</dd>
              </div>
              <div v-if="auth.isAdmin">
                <dt class="text-text-muted text-xs uppercase tracking-wider mb-1">Access</dt>
                <dd class="text-status-success text-sm font-semibold">Admin</dd>
              </div>
            </dl>
            <div class="mt-4 pt-4 border-t border-white/5">
              <p class="text-text-muted text-xs">
                {{ auth.isAdmin
                  ? 'Admin access — monitoring stream.'
                  : 'Your stream is secured and tied to this device. Do not share your link.' }}
              </p>
            </div>
          </div>
        </div>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, RouterLink, useRouter } from 'vue-router'
import { useEventsStore } from '@/stores/events'
import { useAuthStore } from '@/stores/auth'
import { paymentsApi } from '@/api/payments'
import VideoPlayer from '@/components/player/VideoPlayer.vue'

const route = useRoute()
const router = useRouter()
const eventsStore = useEventsStore()
const auth = useAuthStore()

const loading = ref(true)
const error = ref<string | null>(null)
const hlsUrl = ref<string | null>(null)
const event = computed(() => eventsStore.currentEvent)

onMounted(async () => {
  const eventId = route.params.eventId as string
  const token = route.query.token as string | undefined

  // Non-admin without a token has no business being here
  if (!token && !auth.isAdmin) {
    router.replace({
      name: 'event-detail',
      params: { id: eventId },
    })
    return
  }

  await eventsStore.fetchEvent(eventId)

  try {
    // Admins call without a token — backend detects admin role from JWT
    const res = await paymentsApi.getStreamToken(eventId, token)
    hlsUrl.value = res.data.data?.hls_url ?? null
    if (!hlsUrl.value) {
      error.value = 'Stream not yet available. The event may not have started.'
    }
  } catch (e: any) {
    const msg: string = e.response?.data?.error ?? ''
    if (e.response?.status === 403) {
      error.value = 'Your ticket is not valid for this event.'
    } else if (e.response?.status === 404) {
      error.value = 'Stream not yet available. The event may not have started.'
    } else {
      error.value = msg || 'Unable to access the stream. Please try again.'
    }
  } finally {
    loading.value = false
  }
})
</script>
