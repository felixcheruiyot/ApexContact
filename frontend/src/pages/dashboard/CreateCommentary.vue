<template>
  <div class="max-w-2xl mx-auto py-10 px-4">
    <div class="mb-8">
      <h1 class="text-2xl font-bold text-white">Host a Commentary Lobby</h1>
      <p class="text-text-muted mt-1 text-sm">Create an audio discussion room — no approval needed. Go live when you're ready.</p>
    </div>

    <form @submit.prevent="submit" class="space-y-6">
      <!-- Title -->
      <div>
        <label class="form-label">Title *</label>
        <input
          v-model="form.title"
          type="text"
          maxlength="100"
          placeholder="e.g. Ringside Commentary — Nairobi vs Lagos"
          class="form-input"
          required
        />
      </div>

      <!-- Teaser hook -->
      <div>
        <label class="form-label">Teaser Hook</label>
        <input
          v-model="form.teaser_hook"
          type="text"
          maxlength="150"
          placeholder="One-liner that gets fans excited"
          class="form-input"
        />
      </div>

      <!-- Description -->
      <div>
        <label class="form-label">Description</label>
        <textarea
          v-model="form.description"
          rows="3"
          placeholder="Who's fighting? What's the storyline?"
          class="form-input resize-none"
        />
      </div>

      <!-- Sport type -->
      <div>
        <label class="form-label">Sport *</label>
        <div class="flex gap-3">
          <button
            type="button"
            v-for="s in ['boxing', 'racing']"
            :key="s"
            @click="form.sport_type = s"
            class="flex-1 py-2.5 rounded-lg border text-sm font-medium capitalize transition-all"
            :class="form.sport_type === s
              ? 'border-accent-orange bg-accent-orange/10 text-accent-orange'
              : 'border-white/10 text-text-muted hover:border-white/30 hover:text-white'"
          >
            {{ s }}
          </button>
        </div>
      </div>

      <!-- Scheduled at -->
      <div>
        <label class="form-label">Scheduled Date & Time *</label>
        <input
          v-model="form.scheduled_at"
          type="datetime-local"
          class="form-input"
          required
        />
      </div>

      <!-- Price -->
      <div>
        <label class="form-label">Entry Price (KES)</label>
        <div class="flex items-center gap-3">
          <input
            v-model.number="form.price"
            type="number"
            min="0"
            step="1"
            placeholder="0"
            class="form-input"
          />
          <span class="text-text-muted text-sm whitespace-nowrap">0 = Free</span>
        </div>
      </div>

      <!-- Thumbnail URL -->
      <div>
        <label class="form-label">Thumbnail URL</label>
        <input
          v-model="form.thumbnail_url"
          type="url"
          placeholder="https://…"
          class="form-input"
        />
        <div v-if="form.thumbnail_url" class="mt-3 rounded-lg overflow-hidden aspect-video max-w-xs">
          <img :src="form.thumbnail_url" class="w-full h-full object-cover" alt="Thumbnail preview" />
        </div>
      </div>

      <!-- Error -->
      <p v-if="error" class="text-error text-sm">{{ error }}</p>

      <!-- Submit -->
      <button
        type="submit"
        :disabled="submitting"
        class="w-full py-3 rounded-xl bg-accent-orange text-white font-bold text-base
               disabled:opacity-40 hover:bg-orange-500 transition-colors"
      >
        {{ submitting ? 'Creating…' : 'Create Lobby' }}
      </button>
    </form>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { commentaryApi } from '@/api/commentary'

const router = useRouter()

const form = ref({
  title: '',
  teaser_hook: '',
  description: '',
  sport_type: 'boxing',
  scheduled_at: '',
  price: 0,
  thumbnail_url: '',
})

const submitting = ref(false)
const error = ref('')

async function submit() {
  error.value = ''
  if (!form.value.title.trim()) {
    error.value = 'Title is required'
    return
  }
  if (!form.value.sport_type) {
    error.value = 'Sport type is required'
    return
  }
  if (!form.value.scheduled_at) {
    error.value = 'Scheduled date/time is required'
    return
  }

  submitting.value = true
  try {
    const scheduledAt = new Date(form.value.scheduled_at).toISOString()
    const res = await commentaryApi.create({
      ...form.value,
      scheduled_at: scheduledAt,
    })
    const created = res.data.data
    if (created) {
      router.push(`/commentary/${created.id}`)
    }
  } catch (e: any) {
    error.value = e.response?.data?.error ?? 'Failed to create lobby'
  } finally {
    submitting.value = false
  }
}
</script>

<style scoped>
.form-label {
  @apply block text-text-muted text-xs uppercase tracking-wider mb-2;
}
.form-input {
  @apply w-full bg-bg-elevated border border-white/10 rounded-lg px-4 py-2.5 text-white
         placeholder-text-muted focus:outline-none focus:border-accent-orange transition-colors;
}
</style>
