<template>
  <div class="space-y-8">

    <!-- Header -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-white font-bold text-2xl">Platform Overview</h1>
        <p class="text-text-muted text-sm mt-1">Real-time stats and live stream monitoring</p>
      </div>
      <button @click="refresh" class="btn-ghost text-sm py-2 px-4 flex items-center gap-2">
        <svg class="w-4 h-4" :class="{ 'animate-spin': loading }" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
            d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
        </svg>
        Refresh
      </button>
    </div>

    <!-- Error -->
    <div v-if="error" class="bg-status-error/10 border border-status-error/20 rounded-xl px-5 py-4 text-status-error text-sm">
      {{ error }}
    </div>

    <!-- Stats grid skeleton -->
    <div v-if="loading && !stats" class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-5 gap-4">
      <div v-for="i in 5" :key="i" class="card p-5 animate-pulse">
        <div class="h-3 bg-white/10 rounded w-3/4 mb-3" />
        <div class="h-7 bg-white/10 rounded w-1/2" />
      </div>
    </div>

    <!-- Stats grid -->
    <div v-else-if="stats" class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-5 gap-4">
      <div class="card p-5">
        <div class="flex items-center justify-between mb-3">
          <p class="text-text-muted text-xs uppercase tracking-wider">Total Users</p>
          <Users class="w-5 h-5 text-text-muted" />
        </div>
        <p class="text-white font-bold text-3xl">{{ stats.total_users.toLocaleString() }}</p>
      </div>

      <div class="card p-5">
        <div class="flex items-center justify-between mb-3">
          <p class="text-text-muted text-xs uppercase tracking-wider">Total Events</p>
          <Film class="w-5 h-5 text-text-muted" />
        </div>
        <p class="text-white font-bold text-3xl">{{ stats.total_events.toLocaleString() }}</p>
      </div>

      <div class="card p-5 border border-accent-red/20">
        <div class="flex items-center justify-between mb-3">
          <p class="text-text-muted text-xs uppercase tracking-wider">Live Now</p>
          <span class="w-2 h-2 rounded-full bg-accent-red animate-pulse" />
        </div>
        <p class="text-accent-red font-bold text-3xl">{{ stats.live_events }}</p>
      </div>

      <div class="card p-5">
        <div class="flex items-center justify-between mb-3">
          <p class="text-text-muted text-xs uppercase tracking-wider">Revenue</p>
          <DollarSign class="w-5 h-5 text-text-muted" />
        </div>
        <p class="text-status-success font-bold text-2xl">
          KES {{ stats.total_revenue.toLocaleString() }}
        </p>
      </div>

      <div class="card p-5" :class="stats.fraud_flags_open > 0 ? 'border border-status-warning/30' : ''">
        <div class="flex items-center justify-between mb-3">
          <p class="text-text-muted text-xs uppercase tracking-wider">Open Flags</p>
          <Shield class="w-5 h-5 text-text-muted" />
        </div>
        <p :class="stats.fraud_flags_open > 0 ? 'text-status-warning' : 'text-white'" class="font-bold text-3xl">
          {{ stats.fraud_flags_open }}
        </p>
      </div>
    </div>

    <!-- Live streams section -->
    <div>
      <div class="flex items-center justify-between mb-4">
        <div class="flex items-center gap-3">
          <h2 class="text-white font-semibold text-lg">Live Streams</h2>
          <span v-if="liveEvents.length" class="badge-live text-xs">
            <span class="w-1.5 h-1.5 bg-white rounded-full animate-pulse" />
            {{ liveEvents.length }} Live
          </span>
        </div>
        <RouterLink to="/admin/events" class="text-text-muted text-sm hover:text-white transition-colors">
          View all events →
        </RouterLink>
      </div>

      <div v-if="liveEvents.length" class="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-4">
        <div v-for="event in liveEvents" :key="event.id"
          class="card p-5 flex flex-col gap-4 border border-accent-red/10">
          <div class="flex items-start justify-between gap-3">
            <div class="flex-1 min-w-0">
              <div class="flex items-center gap-2 mb-1">
                <span class="badge-live text-xs">
                  <span class="w-1.5 h-1.5 bg-white rounded-full animate-pulse" />
                  Live
                </span>
                <span class="text-text-muted text-xs capitalize">{{ event.sport_type }}</span>
              </div>
              <h3 class="text-white font-semibold text-sm line-clamp-2">{{ event.title }}</h3>
            </div>
            <img
              :src="eventImage(event.thumbnail_url)"
              class="w-16 h-12 object-cover rounded-lg shrink-0"
              :alt="event.title"
              @error="onImageError"
            />
          </div>
          <div class="flex items-center justify-between">
            <span class="text-text-muted text-xs">
              KES {{ event.price.toLocaleString() }} / ticket
            </span>
            <RouterLink
              :to="{ name: 'watch', params: { eventId: event.id } }"
              class="btn-primary text-xs py-1.5 px-3"
            >
              Watch Free
            </RouterLink>
          </div>
        </div>
      </div>

      <div v-else-if="!loading" class="card p-10 text-center">
        <Radio class="w-10 h-10 mx-auto mb-3 text-text-muted" />
        <p class="text-text-muted text-sm">No events are live right now.</p>
      </div>

      <div v-else class="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-4">
        <div v-for="i in 3" :key="i" class="card p-5 animate-pulse">
          <div class="h-4 bg-white/10 rounded w-1/2 mb-2" />
          <div class="h-3 bg-white/10 rounded w-3/4 mb-4" />
          <div class="h-8 bg-white/10 rounded" />
        </div>
      </div>
    </div>

    <!-- Recent fraud flags -->
    <div v-if="recentFlags.length">
      <div class="flex items-center justify-between mb-4">
        <h2 class="text-white font-semibold text-lg flex items-center gap-3">
          Recent Fraud Activity
          <span class="bg-status-warning/20 text-status-warning text-xs font-bold px-2 py-0.5 rounded-full">
            {{ recentFlags.length }} open
          </span>
        </h2>
        <RouterLink to="/admin/fraud" class="text-text-muted text-sm hover:text-white transition-colors">
          View all →
        </RouterLink>
      </div>

      <div class="card overflow-hidden">
        <table class="w-full">
          <thead>
            <tr class="border-b border-white/5">
              <th class="text-left text-text-muted text-xs uppercase tracking-wider px-6 py-3">Reason</th>
              <th class="text-left text-text-muted text-xs uppercase tracking-wider px-6 py-3">User</th>
              <th class="text-left text-text-muted text-xs uppercase tracking-wider px-6 py-3">Detected</th>
              <th class="px-6 py-3" />
            </tr>
          </thead>
          <tbody class="divide-y divide-white/5">
            <tr v-for="flag in recentFlags" :key="flag.id" class="hover:bg-bg-elevated transition-colors">
              <td class="px-6 py-3">
                <span class="text-status-warning text-sm font-medium">{{ flag.reason }}</span>
              </td>
              <td class="px-6 py-3">
                <code class="text-text-muted text-xs bg-bg-surface px-2 py-0.5 rounded">
                  {{ flag.user_id.slice(0, 8) }}…
                </code>
              </td>
              <td class="px-6 py-3 text-text-muted text-sm">
                {{ formatDate(flag.detected_at) }}
              </td>
              <td class="px-6 py-3 text-right">
                <RouterLink to="/admin/fraud" class="text-accent-red text-xs hover:underline">
                  Review
                </RouterLink>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Quick actions -->
    <div class="flex gap-3 pt-2 flex-wrap">
      <RouterLink to="/admin/users" class="btn-ghost text-sm py-2 px-4 flex items-center gap-2">
        <Users class="w-4 h-4" /> Manage Users
      </RouterLink>
      <RouterLink to="/admin/events" class="btn-ghost text-sm py-2 px-4 flex items-center gap-2">
        <Film class="w-4 h-4" /> Manage Events
      </RouterLink>
      <RouterLink to="/admin/fraud" class="btn-primary text-sm py-2 px-4 flex items-center gap-2">
        <Shield class="w-4 h-4" /> Fraud Monitor
        <span v-if="stats?.fraud_flags_open"
          class="bg-white/20 text-white text-xs rounded-full px-1.5 py-0.5 leading-none">
          {{ stats.fraud_flags_open }}
        </span>
      </RouterLink>
    </div>

  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { RouterLink } from 'vue-router'
import { Users, Film, DollarSign, Shield, Radio } from 'lucide-vue-next'
import { format } from 'date-fns'
import { adminApi, type PlatformStats } from '@/api/admin'
import type { FraudFlag, Event } from '@/types'
import { eventImage, onImageError } from '@/utils/eventImage'

const stats = ref<PlatformStats | null>(null)
const liveEvents = ref<Event[]>([])
const recentFlags = ref<FraudFlag[]>([])
const loading = ref(false)
const error = ref<string | null>(null)

function formatDate(d: string) {
  return format(new Date(d), 'MMM d, h:mm a')
}

async function refresh() {
  loading.value = true
  error.value = null
  try {
    const [statsRes, eventsRes, fraudRes] = await Promise.all([
      adminApi.analytics(),
      adminApi.liveEvents(),
      adminApi.listFraudFlags(),
    ])
    stats.value = statsRes.data.data ?? null
    liveEvents.value = eventsRes.data.data ?? []
    recentFlags.value = (fraudRes.data.data ?? []).slice(0, 5)
  } catch (e: any) {
    error.value = e.response?.data?.error ?? 'Failed to load dashboard data'
  } finally {
    loading.value = false
  }
}

onMounted(refresh)
</script>
