<template>
  <div class="max-w-2xl">
    <h1 class="text-white font-bold text-2xl mb-8">Edit Event</h1>

    <!-- Loading -->
    <div v-if="loadingEvent" class="space-y-4 animate-pulse">
      <div class="h-10 bg-white/5 rounded-lg" />
      <div class="h-10 bg-white/5 rounded-lg" />
      <div class="h-24 bg-white/5 rounded-lg" />
    </div>

    <!-- Locked state: event is not draft -->
    <div v-else-if="event && event.status !== 'draft'"
      class="card p-8 text-center space-y-3">
      <svg class="w-12 h-12 text-text-muted mx-auto" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5"
          d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
      </svg>
      <p class="text-white font-semibold">This event cannot be edited</p>
      <p class="text-text-muted text-sm">
        Only events in <span class="text-white font-medium">Draft</span> status can be edited.
        This event is currently <span class="font-medium capitalize" :class="statusTextClass">{{ event.status.replace('_', ' ') }}</span>.
      </p>
      <RouterLink to="/dashboard" class="btn-ghost inline-block mt-2 text-sm py-2 px-4">
        Back to Dashboard
      </RouterLink>
    </div>

    <!-- Edit form -->
    <form v-else-if="event" @submit.prevent="handleSubmit" class="space-y-6">
      <div>
        <label class="block text-text-muted text-sm mb-2">Event Title</label>
        <input v-model="form.title" type="text" placeholder="e.g. Sales Masterclass: Closing High-Ticket Deals 2026" class="input" required />
      </div>

      <div class="grid grid-cols-2 gap-4">
        <div>
          <label class="block text-text-muted text-sm mb-2">Category</label>
          <select v-model="form.sport_type" class="input appearance-none" required>
            <option value="">Select a category</option>
            <optgroup label="Knowledge & Skills">
              <option value="sales">Sales & Negotiation</option>
              <option value="mentoring">Mentoring & Coaching</option>
              <option value="business">Business & Finance</option>
              <option value="education">Education & Workshops</option>
              <option value="visa">Visa & Migration</option>
              <option value="legal">Legal Consultations</option>
              <option value="fitness">Fitness & Wellness</option>
            </optgroup>
            <optgroup label="Entertainment">
              <option value="music">Music & Performances</option>
              <option value="gaming">Gaming & Esports</option>
              <option value="cooking">Cooking & Lifestyle</option>
              <option value="community">Faith & Community</option>
            </optgroup>
            <optgroup label="Other">
              <option value="other">Other</option>
            </optgroup>
          </select>
        </div>
        <div>
          <label class="block text-text-muted text-sm mb-2">Date & Time</label>
          <input v-model="form.scheduled_at" type="datetime-local" class="input" required />
        </div>
      </div>

      <div>
        <label class="block text-text-muted text-sm mb-2">Description</label>
        <textarea v-model="form.description" rows="4" placeholder="Describe the event..." class="input resize-none" />
      </div>

      <div>
        <label class="block text-text-muted text-sm mb-2">Thumbnail URL</label>
        <input v-model="form.thumbnail_url" type="url" placeholder="https://..." class="input" />
        <p class="text-text-muted text-xs mt-1">Link to the event banner image (1280×720 recommended)</p>
      </div>

      <div class="grid grid-cols-2 gap-4">
        <div>
          <label class="block text-text-muted text-sm mb-2">Ticket Price</label>
          <input v-model.number="form.price" type="number" min="0" step="1" placeholder="500" class="input" required />
        </div>
        <div>
          <label class="block text-text-muted text-sm mb-2">Currency</label>
          <select v-model="form.currency" class="input appearance-none">
            <option value="KES">KES (Kenyan Shilling)</option>
            <option value="USD">USD</option>
          </select>
        </div>
      </div>

      <!-- Visibility -->
      <div class="flex items-center justify-between p-4 rounded-lg bg-bg-surface border border-white/5">
        <div>
          <p class="text-white text-sm font-medium">Make this event public</p>
          <p class="text-text-muted text-xs mt-0.5">Public events appear on the Discover page. Private events are only accessible via direct link.</p>
        </div>
        <button
          type="button"
          @click="form.is_public = !form.is_public"
          class="relative w-11 h-6 rounded-full transition-colors duration-200 focus:outline-none"
          :class="form.is_public ? 'bg-accent-red' : 'bg-white/10'"
        >
          <span class="absolute top-0.5 left-0.5 w-5 h-5 bg-white rounded-full shadow transition-transform duration-200"
            :class="form.is_public ? 'translate-x-5' : 'translate-x-0'" />
        </button>
      </div>

      <div v-if="errorMsg" class="bg-status-error/10 border border-status-error/30 text-status-error
                                   text-sm rounded-lg px-4 py-3">
        {{ errorMsg }}
      </div>

      <div class="flex items-center gap-4">
        <button type="submit" class="btn-primary" :disabled="loading">
          {{ loading ? 'Saving…' : 'Save Changes' }}
        </button>
        <RouterLink to="/dashboard" class="btn-ghost">Cancel</RouterLink>
      </div>
    </form>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { RouterLink, useRouter, useRoute } from 'vue-router'
import { eventsApi } from '@/api/events'
import type { Event, SportType } from '@/types'

const router = useRouter()
const route = useRoute()
const eventId = route.params.eventId as string

const event = ref<Event | null>(null)
const loadingEvent = ref(true)
const errorMsg = ref('')
const loading = ref(false)

const form = ref<{
  title: string
  description: string
  sport_type: SportType | ''
  scheduled_at: string
  price: number
  currency: string
  thumbnail_url: string
  is_public: boolean
}>({
  title: '',
  description: '',
  sport_type: '',
  scheduled_at: '',
  price: 0,
  currency: 'KES',
  thumbnail_url: '',
  is_public: false,
})

const statusTextClass = computed(() => {
  switch (event.value?.status) {
    case 'pending_review': return 'text-blue-400'
    case 'live': return 'text-accent-red'
    case 'scheduled': return 'text-accent-orange'
    case 'completed': return 'text-status-success'
    case 'declined': return 'text-status-error'
    default: return 'text-text-muted'
  }
})

function toLocalInput(iso: string): string {
  const d = new Date(iso)
  const pad = (n: number) => String(n).padStart(2, '0')
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())}T${pad(d.getHours())}:${pad(d.getMinutes())}`
}

onMounted(async () => {
  try {
    const res = await eventsApi.get(eventId)
    event.value = res.data.data ?? null
    if (event.value && event.value.status === 'draft') {
      form.value = {
        title: event.value.title,
        description: event.value.description,
        sport_type: event.value.sport_type,
        scheduled_at: toLocalInput(event.value.scheduled_at),
        price: event.value.price,
        currency: event.value.currency,
        thumbnail_url: event.value.thumbnail_url,
        is_public: event.value.is_public,
      }
    }
  } catch {
    errorMsg.value = 'Failed to load event'
  } finally {
    loadingEvent.value = false
  }
})

async function handleSubmit() {
  if (!form.value.sport_type) return
  errorMsg.value = ''
  loading.value = true
  try {
    await eventsApi.update(eventId, {
      ...form.value,
      sport_type: form.value.sport_type as SportType,
      scheduled_at: new Date(form.value.scheduled_at).toISOString(),
      is_public: form.value.is_public,
    })
    router.push('/dashboard')
  } catch (e: any) {
    errorMsg.value = e.response?.data?.error ?? 'Failed to save event'
  } finally {
    loading.value = false
  }
}
</script>
