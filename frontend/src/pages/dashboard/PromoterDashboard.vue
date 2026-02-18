<template>
  <div class="space-y-8">
    <div class="flex items-center justify-between">
      <h1 class="text-white font-bold text-2xl">My Events</h1>
      <RouterLink to="/dashboard/create" class="btn-primary text-sm py-2 px-4">+ New Event</RouterLink>
    </div>

    <!-- Stats row -->
    <div class="grid grid-cols-1 sm:grid-cols-3 gap-4">
      <div class="card p-5">
        <p class="text-text-muted text-xs uppercase tracking-wider mb-2">Total Events</p>
        <p class="text-white font-bold text-3xl">{{ events.length }}</p>
      </div>
      <div class="card p-5">
        <p class="text-text-muted text-xs uppercase tracking-wider mb-2">Live Now</p>
        <p class="text-accent-red font-bold text-3xl">
          {{ events.filter(e => e.status === 'live').length }}
        </p>
      </div>
      <div class="card p-5">
        <p class="text-text-muted text-xs uppercase tracking-wider mb-2">Upcoming</p>
        <p class="text-white font-bold text-3xl">
          {{ events.filter(e => e.status === 'scheduled').length }}
        </p>
      </div>
    </div>

    <!-- Events table -->
    <div class="card overflow-hidden">
      <div class="overflow-x-auto">
        <table class="w-full">
          <thead>
            <tr class="border-b border-white/5">
              <th class="text-left text-text-muted text-xs uppercase tracking-wider px-6 py-4">Event</th>
              <th class="text-left text-text-muted text-xs uppercase tracking-wider px-6 py-4">Sport</th>
              <th class="text-left text-text-muted text-xs uppercase tracking-wider px-6 py-4">Date</th>
              <th class="text-left text-text-muted text-xs uppercase tracking-wider px-6 py-4">Status</th>
              <th class="text-left text-text-muted text-xs uppercase tracking-wider px-6 py-4">Price</th>
              <th class="px-6 py-4" />
            </tr>
          </thead>
          <tbody class="divide-y divide-white/5">
            <tr v-for="event in events" :key="event.id" class="hover:bg-bg-elevated transition-colors">
              <td class="px-6 py-4">
                <p class="text-white font-medium text-sm line-clamp-1">{{ event.title }}</p>
              </td>
              <td class="px-6 py-4">
                <span class="text-text-muted text-sm capitalize">{{ event.sport_type }}</span>
              </td>
              <td class="px-6 py-4">
                <span class="text-text-muted text-sm">{{ format(new Date(event.scheduled_at), 'MMM d, yyyy') }}</span>
              </td>
              <td class="px-6 py-4">
                <span v-if="event.status === 'live'" class="badge-live text-xs">Live</span>
                <span v-else class="badge-upcoming text-xs capitalize">{{ event.status }}</span>
              </td>
              <td class="px-6 py-4">
                <span class="text-white font-semibold text-sm">
                  {{ event.currency }} {{ event.price.toLocaleString() }}
                </span>
              </td>
              <td class="px-6 py-4">
                <div class="flex items-center gap-3">
                <button @click="showStreamKey(event.id)"
                  class="text-accent-orange text-sm hover:underline">
                  Stream Key
                </button>
                <RouterLink :to="`/dashboard/analytics/${event.id}`"
                  class="text-accent-red text-sm hover:underline">
                  Analytics
                </RouterLink>
              </div>
              </td>
            </tr>
          </tbody>
        </table>
    <!-- Stream Key Modal -->
    <Teleport to="body">
      <div v-if="streamKeyModal" class="fixed inset-0 z-50 flex items-center justify-center p-4">
        <div class="absolute inset-0 bg-black/70 backdrop-blur-sm" @click="streamKeyModal = null" />
        <div class="relative bg-bg-elevated border border-white/10 rounded-2xl w-full max-w-lg p-6 shadow-2xl">
          <h3 class="text-white font-bold text-lg mb-1">OBS / Streaming Setup</h3>
          <p class="text-text-muted text-sm mb-5">Configure these in OBS → Settings → Stream</p>

          <div class="space-y-4">
            <div>
              <label class="block text-text-muted text-xs uppercase tracking-wider mb-2">RTMP Server URL</label>
              <div class="flex items-center gap-2">
                <code class="flex-1 bg-bg-surface text-accent-orange text-sm px-4 py-3 rounded-lg font-mono truncate">
                  rtmp://{{ hostname }}:1935/live
                </code>
              </div>
            </div>
            <div>
              <label class="block text-text-muted text-xs uppercase tracking-wider mb-2">Stream Key</label>
              <div class="flex items-center gap-2">
                <code class="flex-1 bg-bg-surface text-white text-sm px-4 py-3 rounded-lg font-mono truncate">
                  {{ streamKeyModal.stream_key }}
                </code>
                <button @click="copyToClipboard(streamKeyModal!.stream_key)"
                  class="btn-ghost text-xs py-2 px-3 shrink-0">Copy</button>
              </div>
              <p class="text-status-error text-xs mt-2">Keep this secret — do not share it.</p>
            </div>
          </div>

          <div class="mt-5 pt-5 border-t border-white/5 text-text-muted text-xs space-y-1">
            <p>1. Open OBS → Settings → Stream → Custom</p>
            <p>2. Paste the Server URL and Stream Key above</p>
            <p>3. Click "Start Streaming" — your event will go <span class="text-accent-red font-bold">LIVE</span> automatically</p>
          </div>

          <button @click="streamKeyModal = null" class="btn-primary w-full mt-5 text-sm">Done</button>
        </div>
      </div>
    </Teleport>

        <div v-if="!events.length" class="text-center py-16">
          <p class="text-text-muted text-sm">No events yet.</p>
          <RouterLink to="/dashboard/create" class="btn-primary mt-4 inline-block text-sm py-2 px-4">
            Create your first event
          </RouterLink>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { RouterLink } from 'vue-router'
import { format } from 'date-fns'
import { eventsApi } from '@/api/events'
import client from '@/api/client'
import type { Event } from '@/types'

const events = ref<Event[]>([])
const streamKeyModal = ref<{ stream_key: string; rtmp_url: string; push_to: string } | null>(null)
const hostname = window.location.hostname

onMounted(async () => {
  const res = await eventsApi.myEvents()
  events.value = res.data.data ?? []
})

async function showStreamKey(eventId: string) {
  const res = await client.get(`/promoter/stream-key/${eventId}`)
  streamKeyModal.value = res.data.data
}

function copyToClipboard(text: string) {
  window.navigator.clipboard.writeText(text)
}
</script>
