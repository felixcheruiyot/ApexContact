<template>
  <div class="space-y-6">
    <div class="flex items-center justify-between">
      <h1 class="text-white font-bold text-2xl">My Streams</h1>
    </div>
    <p class="text-text-muted text-sm">Events you've purchased a ticket for.</p>

    <!-- Loading -->
    <div v-if="loading" class="flex items-center justify-center py-20">
      <div class="w-8 h-8 border-2 border-accent-red border-t-transparent rounded-full animate-spin" />
    </div>

    <!-- Empty state -->
    <div v-else-if="!subscriptions.length" class="card p-12 text-center">
      <svg class="w-12 h-12 text-text-muted mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5"
          d="M15 10l4.553-2.276A1 1 0 0121 8.677v6.646a1 1 0 01-1.447.894L15 14M5 18h8a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v8a2 2 0 002 2z" />
      </svg>
      <p class="text-text-muted text-sm mb-4">You haven't purchased any stream tickets yet.</p>
      <RouterLink to="/discover" class="btn-primary text-sm py-2 px-5">Browse streams</RouterLink>
    </div>

    <!-- Subscriptions grid -->
    <div v-else class="grid sm:grid-cols-2 lg:grid-cols-3 gap-4">
      <div
        v-for="sub in subscriptions"
        :key="sub.subscription_id"
        class="card overflow-hidden group"
      >
        <!-- Thumbnail -->
        <div class="relative aspect-video bg-bg-elevated overflow-hidden">
          <img
            v-if="sub.thumbnail_url"
            :src="sub.thumbnail_url"
            :alt="sub.title"
            class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-300"
          />
          <div v-else class="w-full h-full flex items-center justify-center">
            <svg class="w-10 h-10 text-text-muted/40" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5"
                d="M15 10l4.553-2.276A1 1 0 0121 8.677v6.646a1 1 0 01-1.447.894L15 14M5 18h8a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v8a2 2 0 002 2z" />
            </svg>
          </div>
          <!-- Status badge -->
          <div class="absolute top-2 left-2">
            <span :class="statusBadge(sub.status)" class="text-xs font-bold px-2 py-0.5 rounded-full capitalize">
              {{ sub.status === 'live' ? '● Live' : sub.status }}
            </span>
          </div>
        </div>

        <!-- Info -->
        <div class="p-4 space-y-2">
          <p class="text-white font-semibold text-sm line-clamp-2">{{ sub.title }}</p>
          <p class="text-text-muted text-xs">{{ formatDate(sub.scheduled_at) }}</p>
          <p class="text-text-muted text-xs capitalize">{{ sub.event_type.replace('_', ' ') }}</p>
        </div>

        <!-- Watch button -->
        <div class="px-4 pb-4">
          <RouterLink
            v-if="sub.status === 'live' || sub.status === 'scheduled'"
            :to="watchLink(sub)"
            class="btn-primary w-full text-center text-sm py-2"
          >
            {{ sub.status === 'live' ? 'Watch Now' : 'View Event' }}
          </RouterLink>
          <span v-else class="text-text-muted text-xs">Event ended</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { RouterLink } from 'vue-router'
import { format } from 'date-fns'
import client from '@/api/client'

interface StreamSubscription {
  subscription_id: string
  event_id: string
  title: string
  thumbnail_url: string
  scheduled_at: string
  status: string
  event_type: string
  price: number
  currency: string
  created_at: string
}

const subscriptions = ref<StreamSubscription[]>([])
const loading = ref(true)

onMounted(async () => {
  try {
    const res = await client.get('/profile/subscriptions')
    subscriptions.value = res.data.data ?? []
  } catch {
    subscriptions.value = []
  } finally {
    loading.value = false
  }
})

function formatDate(dateStr: string) {
  return format(new Date(dateStr), 'MMM d, yyyy · h:mm a')
}

function statusBadge(status: string) {
  switch (status) {
    case 'live': return 'bg-accent-red text-white'
    case 'scheduled': return 'bg-accent-orange/20 text-accent-orange'
    case 'completed': return 'bg-white/10 text-text-muted'
    default: return 'bg-white/10 text-text-muted'
  }
}

function watchLink(sub: StreamSubscription) {
  if (sub.event_type === 'video') return `/watch/${sub.event_id}`
  return `/commentary/${sub.event_id}`
}
</script>
