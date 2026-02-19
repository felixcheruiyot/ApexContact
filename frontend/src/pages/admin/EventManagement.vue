<template>
  <div class="space-y-6">

    <!-- Header -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-white font-bold text-2xl">Event Management</h1>
        <p class="text-text-muted text-sm mt-1">Edit event details and manage statuses</p>
      </div>
      <button @click="loadEvents" class="btn-ghost text-sm py-2 px-4 flex items-center gap-2">
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

    <!-- Filters -->
    <div class="flex flex-wrap gap-3">
      <input
        v-model="search"
        type="text"
        placeholder="Search events..."
        class="input text-sm py-2 px-3 w-60"
      />
      <select v-model="filterStatus" class="input text-sm py-2 px-3">
        <option value="">All statuses</option>
        <option value="scheduled">Scheduled</option>
        <option value="live">Live</option>
        <option value="completed">Completed</option>
        <option value="cancelled">Cancelled</option>
      </select>
      <select v-model="filterSport" class="input text-sm py-2 px-3">
        <option value="">All sports</option>
        <option value="boxing">Boxing</option>
        <option value="racing">Racing</option>
      </select>
    </div>

    <!-- Skeleton -->
    <div v-if="loading && !events.length" class="card overflow-hidden animate-pulse">
      <div v-for="i in 5" :key="i" class="px-6 py-4 border-b border-white/5 flex gap-4">
        <div class="h-4 bg-white/10 rounded w-1/3" />
        <div class="h-4 bg-white/10 rounded w-1/6" />
        <div class="h-4 bg-white/10 rounded w-1/6" />
        <div class="h-4 bg-white/10 rounded w-1/8" />
      </div>
    </div>

    <!-- Table -->
    <div v-else-if="filtered.length" class="card overflow-hidden">
      <table class="w-full">
        <thead>
          <tr class="border-b border-white/5">
            <th class="text-left text-text-muted text-xs uppercase tracking-wider px-6 py-3">Event</th>
            <th class="text-left text-text-muted text-xs uppercase tracking-wider px-6 py-3">Sport</th>
            <th class="text-left text-text-muted text-xs uppercase tracking-wider px-6 py-3">Scheduled</th>
            <th class="text-left text-text-muted text-xs uppercase tracking-wider px-6 py-3">Price</th>
            <th class="text-left text-text-muted text-xs uppercase tracking-wider px-6 py-3">Status</th>
            <th class="px-6 py-3" />
          </tr>
        </thead>
        <tbody class="divide-y divide-white/5">
          <tr v-for="event in filtered" :key="event.id" class="hover:bg-surface-2 transition-colors">
            <td class="px-6 py-4 max-w-xs">
              <p class="text-white text-sm font-medium line-clamp-1">{{ event.title }}</p>
              <p class="text-text-muted text-xs mt-0.5 line-clamp-1">{{ event.description }}</p>
            </td>
            <td class="px-6 py-4">
              <span class="text-text-muted text-sm capitalize">{{ event.sport_type }}</span>
            </td>
            <td class="px-6 py-4 text-text-muted text-sm whitespace-nowrap">
              {{ formatDate(event.scheduled_at) }}
            </td>
            <td class="px-6 py-4 text-text-muted text-sm">
              {{ event.currency }} {{ event.price.toLocaleString() }}
            </td>
            <td class="px-6 py-4">
              <span :class="statusClass(event.status)" class="text-xs font-semibold px-2 py-0.5 rounded-full">
                {{ event.status }}
              </span>
            </td>
            <td class="px-6 py-4 text-right">
              <button @click="openEdit(event)" class="btn-ghost text-xs py-1.5 px-3">
                Edit
              </button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <div v-else-if="!loading" class="card p-12 text-center">
      <p class="text-3xl mb-3">🎬</p>
      <p class="text-text-muted text-sm">No events found.</p>
    </div>

    <!-- Edit Modal -->
    <Teleport to="body">
      <div v-if="editing" class="fixed inset-0 z-50 flex items-center justify-center p-4">
        <div class="absolute inset-0 bg-black/60 backdrop-blur-sm" @click="closeEdit" />
        <div class="relative card w-full max-w-lg p-6 space-y-5 max-h-[90vh] overflow-y-auto">

          <div class="flex items-center justify-between">
            <h2 class="text-white font-semibold text-lg">Edit Event</h2>
            <button @click="closeEdit" class="text-text-muted hover:text-white transition-colors">
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>

          <div v-if="saveError" class="bg-status-error/10 border border-status-error/20 rounded-lg px-4 py-3 text-status-error text-sm">
            {{ saveError }}
          </div>

          <div class="space-y-4">
            <!-- Title -->
            <div>
              <label class="block text-text-muted text-xs uppercase tracking-wider mb-1.5">Title</label>
              <input v-model="form.title" type="text" class="input w-full" placeholder="Event title" />
            </div>

            <!-- Description -->
            <div>
              <label class="block text-text-muted text-xs uppercase tracking-wider mb-1.5">Description</label>
              <textarea v-model="form.description" rows="3" class="input w-full resize-none" placeholder="Event description" />
            </div>

            <!-- Sport + Status row -->
            <div class="grid grid-cols-2 gap-4">
              <div>
                <label class="block text-text-muted text-xs uppercase tracking-wider mb-1.5">Sport</label>
                <select v-model="form.sport_type" class="input w-full">
                  <option value="boxing">Boxing</option>
                  <option value="racing">Racing</option>
                </select>
              </div>
              <div>
                <label class="block text-text-muted text-xs uppercase tracking-wider mb-1.5">Status</label>
                <select v-model="form.status" class="input w-full">
                  <option value="scheduled">Scheduled</option>
                  <option value="live">Live</option>
                  <option value="completed">Completed</option>
                  <option value="cancelled">Cancelled</option>
                </select>
              </div>
            </div>

            <!-- Scheduled At -->
            <div>
              <label class="block text-text-muted text-xs uppercase tracking-wider mb-1.5">Scheduled At</label>
              <input v-model="form.scheduled_at" type="datetime-local" class="input w-full" />
            </div>

            <!-- Price + Currency row -->
            <div class="grid grid-cols-2 gap-4">
              <div>
                <label class="block text-text-muted text-xs uppercase tracking-wider mb-1.5">Price</label>
                <input v-model.number="form.price" type="number" min="0" step="1" class="input w-full" placeholder="0" />
              </div>
              <div>
                <label class="block text-text-muted text-xs uppercase tracking-wider mb-1.5">Currency</label>
                <input v-model="form.currency" type="text" class="input w-full" placeholder="KES" maxlength="3" />
              </div>
            </div>

            <!-- Thumbnail URL -->
            <div>
              <label class="block text-text-muted text-xs uppercase tracking-wider mb-1.5">Thumbnail URL</label>
              <input v-model="form.thumbnail_url" type="url" class="input w-full" placeholder="https://..." />
            </div>
          </div>

          <!-- Actions -->
          <div class="flex gap-3 pt-2">
            <button @click="closeEdit" class="btn-ghost flex-1 py-2.5">Cancel</button>
            <button @click="saveEvent" :disabled="saving" class="btn-primary flex-1 py-2.5 flex items-center justify-center gap-2">
              <svg v-if="saving" class="w-4 h-4 animate-spin" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                  d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
              </svg>
              {{ saving ? 'Saving…' : 'Save Changes' }}
            </button>
          </div>
        </div>
      </div>
    </Teleport>

  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { format } from 'date-fns'
import { adminApi } from '@/api/admin'
import type { Event, EventStatus, SportType } from '@/types'

const events = ref<Event[]>([])
const loading = ref(false)
const error = ref<string | null>(null)

const search = ref('')
const filterStatus = ref('')
const filterSport = ref('')

const editing = ref<Event | null>(null)
const saving = ref(false)
const saveError = ref<string | null>(null)

interface EditForm {
  title: string
  description: string
  sport_type: SportType
  scheduled_at: string
  price: number
  currency: string
  thumbnail_url: string
  status: EventStatus
}

const form = ref<EditForm>({
  title: '',
  description: '',
  sport_type: 'boxing',
  scheduled_at: '',
  price: 0,
  currency: 'KES',
  thumbnail_url: '',
  status: 'scheduled',
})

const filtered = computed(() => {
  return events.value.filter(e => {
    const matchSearch = !search.value ||
      e.title.toLowerCase().includes(search.value.toLowerCase()) ||
      e.description.toLowerCase().includes(search.value.toLowerCase())
    const matchStatus = !filterStatus.value || e.status === filterStatus.value
    const matchSport = !filterSport.value || e.sport_type === filterSport.value
    return matchSearch && matchStatus && matchSport
  })
})

function formatDate(d: string) {
  return format(new Date(d), 'MMM d, yyyy h:mm a')
}

function statusClass(status: EventStatus) {
  switch (status) {
    case 'live': return 'bg-accent-red/20 text-accent-red'
    case 'scheduled': return 'bg-accent-orange/20 text-accent-orange'
    case 'completed': return 'bg-status-success/20 text-status-success'
    case 'cancelled': return 'bg-white/10 text-text-muted'
  }
}

// Convert ISO datetime to datetime-local input format (YYYY-MM-DDTHH:mm)
function toLocalInput(iso: string) {
  const d = new Date(iso)
  const pad = (n: number) => String(n).padStart(2, '0')
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())}T${pad(d.getHours())}:${pad(d.getMinutes())}`
}

async function loadEvents() {
  loading.value = true
  error.value = null
  try {
    const res = await adminApi.listAllEvents()
    events.value = res.data.data ?? []
  } catch (e: any) {
    error.value = e.response?.data?.error ?? 'Failed to load events'
  } finally {
    loading.value = false
  }
}

function openEdit(event: Event) {
  editing.value = event
  saveError.value = null
  form.value = {
    title: event.title,
    description: event.description,
    sport_type: event.sport_type,
    scheduled_at: toLocalInput(event.scheduled_at),
    price: event.price,
    currency: event.currency,
    thumbnail_url: event.thumbnail_url,
    status: event.status,
  }
}

function closeEdit() {
  editing.value = null
  saveError.value = null
}

async function saveEvent() {
  if (!editing.value) return
  saving.value = true
  saveError.value = null
  try {
    await adminApi.updateEvent(editing.value.id, {
      ...form.value,
      scheduled_at: new Date(form.value.scheduled_at).toISOString(),
    })
    // Patch local list so the table updates immediately
    const idx = events.value.findIndex(e => e.id === editing.value!.id)
    if (idx !== -1) {
      events.value[idx] = {
        ...events.value[idx],
        ...form.value,
        scheduled_at: new Date(form.value.scheduled_at).toISOString(),
      }
    }
    closeEdit()
  } catch (e: any) {
    saveError.value = e.response?.data?.error ?? 'Failed to save event'
  } finally {
    saving.value = false
  }
}

onMounted(loadEvents)
</script>
