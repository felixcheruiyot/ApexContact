<template>
  <div class="max-w-2xl">
    <h1 class="text-white font-bold text-2xl mb-8">Create New Event</h1>

    <form @submit.prevent="handleSubmit" class="space-y-6">
      <div>
        <label class="block text-text-muted text-sm mb-2">Event Title</label>
        <input v-model="form.title" type="text" placeholder="e.g. Nairobi Career Masterclass 2026" class="input" required />
      </div>

      <div class="grid grid-cols-2 gap-4">
        <div>
          <label class="block text-text-muted text-sm mb-2">Category</label>
          <select v-model="form.sport_type" class="input appearance-none" required>
            <option value="">Select a category</option>
            <optgroup label="Sports">
              <option value="boxing">🥊 Boxing</option>
              <option value="racing">🏎️ Motorsports & Racing</option>
              <option value="fitness">💪 Fitness & Wellness</option>
            </optgroup>
            <optgroup label="Knowledge & Consulting">
              <option value="mentoring">🏫 Mentoring & Coaching</option>
              <option value="education">📚 Education & Workshops</option>
              <option value="legal">⚖️ Legal Consultations</option>
              <option value="visa">🌍 Visa & Immigration</option>
              <option value="business">💼 Business & Finance</option>
            </optgroup>
            <optgroup label="Entertainment">
              <option value="music">🎵 Music & Performances</option>
              <option value="gaming">🎮 Gaming & Esports</option>
              <option value="cooking">🍳 Cooking & Lifestyle</option>
              <option value="community">🙏 Faith & Community</option>
            </optgroup>
            <optgroup label="Other">
              <option value="other">📌 Other</option>
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

      <div v-if="errorMsg" class="bg-status-error/10 border border-status-error/30 text-status-error
                                   text-sm rounded-lg px-4 py-3">
        {{ errorMsg }}
      </div>

      <div class="flex items-center gap-4">
        <button type="submit" class="btn-primary" :disabled="loading">
          {{ loading ? 'Creating...' : 'Create Event' }}
        </button>
        <RouterLink to="/dashboard" class="btn-ghost">Cancel</RouterLink>
      </div>
    </form>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { RouterLink, useRouter } from 'vue-router'
import { eventsApi } from '@/api/events'
import type { SportType } from '@/types'

const router = useRouter()
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
}>({
  title: '',
  description: '',
  sport_type: '',
  scheduled_at: '',
  price: 0,
  currency: 'KES',
  thumbnail_url: '',
})

async function handleSubmit() {
  if (!form.value.sport_type) return
  errorMsg.value = ''
  loading.value = true
  try {
    await eventsApi.create({
      ...form.value,
      sport_type: form.value.sport_type as SportType,
      scheduled_at: new Date(form.value.scheduled_at).toISOString(),
    })
    router.push('/dashboard')
  } catch (e: any) {
    errorMsg.value = e.response?.data?.error ?? 'Failed to create event'
  } finally {
    loading.value = false
  }
}
</script>
