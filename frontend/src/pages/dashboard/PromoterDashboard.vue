<template>
  <div class="space-y-8">
    <div class="flex items-center justify-between">
      <h1 class="text-white font-bold text-2xl">My Events</h1>
      <RouterLink to="/dashboard/create" class="btn-primary text-sm py-2 px-4">+ New Event</RouterLink>
    </div>

    <!-- Stats row -->
    <div class="grid grid-cols-2 sm:grid-cols-4 gap-4">
      <div class="card p-5">
        <p class="text-text-muted text-xs uppercase tracking-wider mb-2">Draft</p>
        <p class="text-white font-bold text-3xl">{{ counts.draft }}</p>
      </div>
      <div class="card p-5">
        <p class="text-text-muted text-xs uppercase tracking-wider mb-2">Pending Review</p>
        <p class="text-accent-orange font-bold text-3xl">{{ counts.pending_review }}</p>
      </div>
      <div class="card p-5">
        <p class="text-text-muted text-xs uppercase tracking-wider mb-2">Live Now</p>
        <p class="text-accent-red font-bold text-3xl">{{ counts.live }}</p>
      </div>
      <div class="card p-5">
        <p class="text-text-muted text-xs uppercase tracking-wider mb-2">Scheduled</p>
        <p class="text-white font-bold text-3xl">{{ counts.scheduled }}</p>
      </div>
    </div>

    <!-- Events table -->
    <div class="card overflow-hidden">
      <div class="overflow-x-auto">
        <table class="w-full min-w-[600px]">
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
          <tbody>
            <template v-for="event in events" :key="event.id">
              <!-- Review note banner row -->
              <tr v-if="event.status === 'draft' && event.review_note"
                class="bg-accent-orange/5 border-b border-accent-orange/10">
                <td colspan="6" class="px-6 py-2">
                  <div class="flex items-start gap-2 text-sm">
                    <svg class="w-4 h-4 text-accent-orange mt-0.5 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                        d="M12 9v2m0 4h.01M10.29 3.86L1.82 18a2 2 0 001.71 3h16.94a2 2 0 001.71-3L13.71 3.86a2 2 0 00-3.42 0z" />
                    </svg>
                    <span>
                      <span class="font-semibold text-accent-orange">Edits Requested: </span>
                      <span class="text-text-muted">{{ event.review_note }}</span>
                    </span>
                  </div>
                </td>
              </tr>
              <tr v-if="event.status === 'declined'"
                class="bg-status-error/5 border-b border-status-error/10">
                <td colspan="6" class="px-6 py-2">
                  <div class="flex items-start gap-2 text-sm">
                    <svg class="w-4 h-4 text-status-error mt-0.5 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                        d="M18.364 18.364A9 9 0 005.636 5.636m12.728 12.728A9 9 0 015.636 5.636m12.728 12.728L5.636 5.636" />
                    </svg>
                    <span>
                      <span class="font-semibold text-status-error">Declined: </span>
                      <span class="text-text-muted">{{ event.review_note }}</span>
                    </span>
                  </div>
                </td>
              </tr>
              <!-- Main data row -->
              <tr class="border-b border-white/5 hover:bg-bg-elevated transition-colors">
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
                  <span :class="statusBadgeClass(event.status)" class="text-xs font-semibold px-2 py-0.5 rounded-full capitalize">
                    {{ statusLabel(event.status) }}
                  </span>
                </td>
                <td class="px-6 py-4">
                  <span class="text-white font-semibold text-sm">
                    {{ event.currency }} {{ event.price.toLocaleString() }}
                  </span>
                </td>
                <td class="px-6 py-4">
                  <div class="flex items-center gap-3 justify-end">
                    <!-- Draft actions -->
                    <template v-if="event.status === 'draft'">
                      <RouterLink :to="`/dashboard/edit/${event.id}`"
                        class="text-text-muted text-sm hover:text-white transition-colors">
                        Edit
                      </RouterLink>
                      <button @click="submitEvent(event.id)"
                        :disabled="submitting === event.id"
                        class="text-accent-orange text-sm hover:underline disabled:opacity-50">
                        {{ submitting === event.id ? 'Submitting…' : (event.review_note ? 'Resubmit' : 'Submit for Review') }}
                      </button>
                    </template>

                    <!-- Pending review: no actions -->
                    <template v-else-if="event.status === 'pending_review'">
                      <span class="text-text-muted text-xs italic">Under Review</span>
                    </template>

                    <!-- Declined: no actions -->
                    <template v-else-if="event.status === 'declined'">
                      <!-- no actions -->
                    </template>

                    <!-- Active events -->
                    <template v-else>
                      <button @click="showStreamKey(event.id)"
                        class="text-accent-orange text-sm hover:underline">
                        Stream Key
                      </button>
                      <RouterLink :to="`/dashboard/analytics/${event.id}`"
                        class="text-accent-red text-sm hover:underline">
                        Analytics
                      </RouterLink>
                    </template>
                  </div>
                </td>
              </tr>
            </template>
          </tbody>
        </table>

        <div v-if="!events.length" class="text-center py-16">
          <p class="text-text-muted text-sm">No events yet.</p>
          <RouterLink to="/dashboard/create" class="btn-primary mt-4 inline-block text-sm py-2 px-4">
            Create your first event
          </RouterLink>
        </div>
      </div>
    </div>

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
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { RouterLink } from 'vue-router'
import { format } from 'date-fns'
import { eventsApi } from '@/api/events'
import client from '@/api/client'
import type { Event, EventStatus } from '@/types'

const events = ref<Event[]>([])
const streamKeyModal = ref<{ stream_key: string; rtmp_url: string; push_to: string } | null>(null)
const submitting = ref<string | null>(null)
const hostname = window.location.hostname

const counts = computed(() => ({
  draft: events.value.filter(e => e.status === 'draft').length,
  pending_review: events.value.filter(e => e.status === 'pending_review').length,
  live: events.value.filter(e => e.status === 'live').length,
  scheduled: events.value.filter(e => e.status === 'scheduled').length,
}))

onMounted(async () => {
  const res = await eventsApi.myEvents()
  events.value = res.data.data ?? []
})

async function submitEvent(eventId: string) {
  submitting.value = eventId
  try {
    await eventsApi.submit(eventId)
    const idx = events.value.findIndex(e => e.id === eventId)
    if (idx !== -1) {
      events.value[idx] = { ...events.value[idx], status: 'pending_review', review_note: '' }
    }
  } catch (e: any) {
    alert(e.response?.data?.error ?? 'Failed to submit event')
  } finally {
    submitting.value = null
  }
}

async function showStreamKey(eventId: string) {
  const res = await client.get(`/promoter/stream-key/${eventId}`)
  streamKeyModal.value = res.data.data
}

function copyToClipboard(text: string) {
  window.navigator.clipboard.writeText(text)
}

function statusLabel(status: EventStatus): string {
  if (status === 'pending_review') return 'Pending Review'
  return status
}

function statusBadgeClass(status: EventStatus): string {
  switch (status) {
    case 'live': return 'bg-accent-red/20 text-accent-red'
    case 'scheduled': return 'bg-accent-orange/20 text-accent-orange'
    case 'pending_review': return 'bg-blue-500/20 text-blue-400'
    case 'draft': return 'bg-white/10 text-text-muted'
    case 'completed': return 'bg-status-success/20 text-status-success'
    case 'declined': return 'bg-status-error/20 text-status-error'
    default: return 'bg-white/10 text-text-muted'
  }
}
</script>
