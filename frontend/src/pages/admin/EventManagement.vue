<template>
  <div class="space-y-6">

    <!-- Header -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-white font-bold text-2xl">Event Management</h1>
        <p class="text-text-muted text-sm mt-1">Review, approve, and manage all events</p>
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

    <!-- Tabs -->
    <div class="flex gap-1 overflow-x-auto border-b border-white/5 pb-0 scrollbar-none">
      <button
        v-for="tab in tabs" :key="tab.value"
        @click="activeTab = tab.value"
        :class="[
          'px-4 py-2 text-sm font-medium rounded-t-lg transition-colors relative shrink-0',
          activeTab === tab.value
            ? 'text-white bg-bg-elevated border-b-2 border-accent-red'
            : 'text-text-muted hover:text-white'
        ]"
      >
        {{ tab.label }}
        <span v-if="tab.value === 'pending_review' && pendingCount > 0"
          class="ml-1.5 inline-flex items-center justify-center w-5 h-5 text-xs font-bold rounded-full bg-accent-red text-white">
          {{ pendingCount }}
        </span>
      </button>
    </div>

    <!-- Search -->
    <div class="flex flex-col sm:flex-row flex-wrap gap-3">
      <input
        v-model="search"
        type="text"
        placeholder="Search events..."
        class="input text-sm py-2 px-3 w-full sm:w-60"
      />
      <select v-model="filterSport" class="input text-sm py-2 px-3 w-full sm:w-auto">
        <option value="">All categories</option>
        <option value="sales">Sales & Negotiation</option>
        <option value="mentoring">Mentoring & Coaching</option>
        <option value="business">Business & Finance</option>
        <option value="education">Education & Workshops</option>
        <option value="visa">Visa & Migration</option>
        <option value="legal">Legal Consultations</option>
        <option value="fitness">Fitness & Wellness</option>
        <option value="music">Music & Performances</option>
        <option value="gaming">Gaming & Esports</option>
        <option value="cooking">Cooking & Lifestyle</option>
        <option value="community">Faith & Community</option>
        <option value="other">Other</option>
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
      <div class="overflow-x-auto">
      <table class="w-full min-w-[640px]">
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
              <p v-if="event.review_note" class="text-text-muted text-xs mt-1 italic line-clamp-1">
                Note: {{ event.review_note }}
              </p>
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
              <span :class="statusClass(event.status)" class="text-xs font-semibold px-2 py-0.5 rounded-full capitalize">
                {{ event.status.replace('_', ' ') }}
              </span>
            </td>
            <td class="px-6 py-4 text-right">
              <!-- Pending review actions -->
              <div v-if="event.status === 'pending_review'" class="flex items-center gap-2 justify-end">
                <button @click="approve(event)"
                  :disabled="actioning === event.id"
                  class="text-xs font-semibold px-3 py-1.5 rounded-lg bg-status-success/10 text-status-success
                         hover:bg-status-success/20 transition-colors disabled:opacity-50">
                  Approve
                </button>
                <button @click="openReasonModal(event, 'request_edits')"
                  :disabled="actioning === event.id"
                  class="text-xs font-semibold px-3 py-1.5 rounded-lg bg-accent-orange/10 text-accent-orange
                         hover:bg-accent-orange/20 transition-colors disabled:opacity-50">
                  Request Edits
                </button>
                <button @click="openReasonModal(event, 'decline')"
                  :disabled="actioning === event.id"
                  class="text-xs font-semibold px-3 py-1.5 rounded-lg bg-status-error/10 text-status-error
                         hover:bg-status-error/20 transition-colors disabled:opacity-50">
                  Decline
                </button>
              </div>
              <!-- Other events: edit -->
              <button v-else @click="openEdit(event)" class="btn-ghost text-xs py-1.5 px-3">
                Edit
              </button>
            </td>
          </tr>
        </tbody>
      </table>
      </div>
    </div>

    <div v-else-if="!loading" class="card p-12 text-center">
      <p class="text-3xl mb-3">🎬</p>
      <p class="text-text-muted text-sm">No events found.</p>
    </div>

    <!-- Reason Modal (Request Edits / Decline) -->
    <Teleport to="body">
      <div v-if="reasonModal" class="fixed inset-0 z-50 flex items-center justify-center p-4">
        <div class="absolute inset-0 bg-black/60 backdrop-blur-sm" @click="closeReasonModal" />
        <div class="relative card w-full max-w-md p-6 space-y-4">
          <div class="flex items-center justify-between">
            <h2 class="text-white font-semibold text-lg">
              {{ reasonModal.action === 'decline' ? 'Decline Event' : 'Request Edits' }}
            </h2>
            <button @click="closeReasonModal" class="text-text-muted hover:text-white transition-colors">
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>

          <p class="text-text-muted text-sm">
            {{ reasonModal.action === 'decline'
              ? 'This action is final. The promoter will not be able to edit or resubmit this event.'
              : 'Describe the changes the promoter needs to make before the event can be approved.' }}
          </p>

          <div v-if="reasonError" class="bg-status-error/10 border border-status-error/20 rounded-lg px-4 py-3 text-status-error text-sm">
            {{ reasonError }}
          </div>

          <div>
            <label class="block text-text-muted text-xs uppercase tracking-wider mb-1.5">Reason *</label>
            <textarea
              v-model="reasonText"
              rows="4"
              class="input w-full resize-none"
              :placeholder="reasonModal.action === 'decline' ? 'Why is this event being declined?' : 'What needs to be changed?'"
            />
          </div>

          <div class="flex gap-3 pt-1">
            <button @click="closeReasonModal" class="btn-ghost flex-1 py-2.5">Cancel</button>
            <button
              @click="submitReason"
              :disabled="actioning === reasonModal.event.id"
              :class="[
                'flex-1 py-2.5 rounded-lg font-semibold text-sm transition-colors disabled:opacity-50',
                reasonModal.action === 'decline'
                  ? 'bg-status-error/20 text-status-error hover:bg-status-error/30'
                  : 'bg-accent-orange/20 text-accent-orange hover:bg-accent-orange/30'
              ]"
            >
              {{ actioning === reasonModal.event.id ? 'Submitting…' : (reasonModal.action === 'decline' ? 'Decline Event' : 'Send to Promoter') }}
            </button>
          </div>
        </div>
      </div>
    </Teleport>

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
            <div>
              <label class="block text-text-muted text-xs uppercase tracking-wider mb-1.5">Title</label>
              <input v-model="form.title" type="text" class="input w-full" placeholder="Event title" />
            </div>
            <div>
              <label class="block text-text-muted text-xs uppercase tracking-wider mb-1.5">Description</label>
              <textarea v-model="form.description" rows="3" class="input w-full resize-none" placeholder="Event description" />
            </div>
            <div class="grid grid-cols-2 gap-4">
              <div>
                <label class="block text-text-muted text-xs uppercase tracking-wider mb-1.5">Category</label>
                <select v-model="form.sport_type" class="input w-full">
                  <option value="sales">Sales & Negotiation</option>
                  <option value="mentoring">Mentoring & Coaching</option>
                  <option value="business">Business & Finance</option>
                  <option value="education">Education & Workshops</option>
                  <option value="visa">Visa & Migration</option>
                  <option value="legal">Legal Consultations</option>
                  <option value="fitness">Fitness & Wellness</option>
                  <option value="music">Music & Performances</option>
                  <option value="gaming">Gaming & Esports</option>
                  <option value="cooking">Cooking & Lifestyle</option>
                  <option value="community">Faith & Community</option>
                  <option value="other">Other</option>
                </select>
              </div>
              <div>
                <label class="block text-text-muted text-xs uppercase tracking-wider mb-1.5">Status</label>
                <select v-model="form.status" class="input w-full">
                  <option value="draft">Draft</option>
                  <option value="pending_review">Pending Review</option>
                  <option value="scheduled">Scheduled</option>
                  <option value="live">Live</option>
                  <option value="completed">Completed</option>
                  <option value="cancelled">Cancelled</option>
                  <option value="declined">Declined</option>
                </select>
              </div>
            </div>
            <div>
              <label class="block text-text-muted text-xs uppercase tracking-wider mb-1.5">Scheduled At</label>
              <input v-model="form.scheduled_at" type="datetime-local" class="input w-full" />
            </div>
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
            <div>
              <label class="block text-text-muted text-xs uppercase tracking-wider mb-1.5">Thumbnail URL</label>
              <input v-model="form.thumbnail_url" type="url" class="input w-full" placeholder="https://..." />
            </div>
          </div>

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
const actioning = ref<string | null>(null)

const search = ref('')
const filterSport = ref('')

const activeTab = ref<EventStatus | 'all'>('all')

const tabs: { label: string; value: EventStatus | 'all' }[] = [
  { label: 'All', value: 'all' },
  { label: 'Pending Review', value: 'pending_review' },
  { label: 'Draft', value: 'draft' },
  { label: 'Scheduled', value: 'scheduled' },
  { label: 'Live', value: 'live' },
  { label: 'Completed', value: 'completed' },
  { label: 'Declined', value: 'declined' },
]

const pendingCount = computed(() =>
  events.value.filter(e => e.status === 'pending_review').length,
)

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
  sport_type: 'mentoring',
  scheduled_at: '',
  price: 0,
  currency: 'KES',
  thumbnail_url: '',
  status: 'scheduled',
})

// Reason modal
const reasonModal = ref<{ event: Event; action: 'request_edits' | 'decline' } | null>(null)
const reasonText = ref('')
const reasonError = ref<string | null>(null)

const filtered = computed(() => {
  return events.value.filter(e => {
    const matchTab = activeTab.value === 'all' || e.status === activeTab.value
    const matchSearch = !search.value ||
      e.title.toLowerCase().includes(search.value.toLowerCase()) ||
      e.description.toLowerCase().includes(search.value.toLowerCase())
    const matchSport = !filterSport.value || e.sport_type === filterSport.value
    return matchTab && matchSearch && matchSport
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
    case 'pending_review': return 'bg-blue-500/20 text-blue-400'
    case 'draft': return 'bg-white/10 text-text-muted'
    case 'declined': return 'bg-status-error/20 text-status-error'
    default: return 'bg-white/10 text-text-muted'
  }
}

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

async function approve(event: Event) {
  actioning.value = event.id
  try {
    await adminApi.approveEvent(event.id)
    const idx = events.value.findIndex(e => e.id === event.id)
    if (idx !== -1) events.value[idx] = { ...events.value[idx], status: 'scheduled', review_note: '' }
  } catch (e: any) {
    error.value = e.response?.data?.error ?? 'Failed to approve event'
  } finally {
    actioning.value = null
  }
}

function openReasonModal(event: Event, action: 'request_edits' | 'decline') {
  reasonModal.value = { event, action }
  reasonText.value = ''
  reasonError.value = null
}

function closeReasonModal() {
  reasonModal.value = null
  reasonText.value = ''
  reasonError.value = null
}

async function submitReason() {
  if (!reasonModal.value) return
  if (!reasonText.value.trim()) {
    reasonError.value = 'Reason is required'
    return
  }
  const { event, action } = reasonModal.value
  actioning.value = event.id
  reasonError.value = null
  try {
    if (action === 'request_edits') {
      await adminApi.requestEdits(event.id, reasonText.value.trim())
      const idx = events.value.findIndex(e => e.id === event.id)
      if (idx !== -1) events.value[idx] = { ...events.value[idx], status: 'draft', review_note: reasonText.value.trim() }
    } else {
      await adminApi.declineEvent(event.id, reasonText.value.trim())
      const idx = events.value.findIndex(e => e.id === event.id)
      if (idx !== -1) events.value[idx] = { ...events.value[idx], status: 'declined', review_note: reasonText.value.trim() }
    }
    closeReasonModal()
  } catch (e: any) {
    reasonError.value = e.response?.data?.error ?? 'Action failed'
  } finally {
    actioning.value = null
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
