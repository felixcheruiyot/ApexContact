<template>
  <div class="max-w-2xl mx-auto px-4 sm:px-6 py-10">
    <div class="mb-8">
      <RouterLink to="/" class="text-text-muted hover:text-white text-sm transition-colors">← Back to home</RouterLink>
      <h1 class="text-2xl font-bold text-white mt-4">Host a Commentary Lobby</h1>
      <p class="text-text-muted mt-1 text-sm">
        Create an audio discussion room — no approval needed. Go live when you're ready.
      </p>
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
            {{ s === 'boxing' ? '🥊 Boxing' : '🏎️ Racing' }}
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
          <span class="text-text-muted text-sm whitespace-nowrap shrink-0">
            {{ form.price === 0 ? '✓ Free entry' : `Listeners pay KES ${form.price}` }}
          </span>
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
        <div v-if="form.thumbnail_url" class="mt-3 rounded-xl overflow-hidden aspect-video max-w-xs">
          <img :src="form.thumbnail_url" class="w-full h-full object-cover" alt="Preview" />
        </div>
      </div>

      <!-- Info box -->
      <div class="bg-accent-orange/10 border border-accent-orange/20 rounded-xl px-4 py-3 text-sm text-accent-orange">
        <p class="font-semibold mb-1">How it works after you create</p>
        <ol class="list-decimal list-inside space-y-1 text-accent-orange/80">
          <li>Share your lobby link with fans</li>
          <li>When ready, open the lobby and press <strong>Start Room</strong></li>
          <li>Fans join the live audio room and chat in real time</li>
          <li>Press <strong>End Room</strong> when you're done — chat is saved for replay</li>
        </ol>
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
        {{ submitting ? 'Creating…' : '🎙 Create Lobby' }}
      </button>
    </form>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { RouterLink, useRouter } from 'vue-router'
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
  if (!form.value.title.trim()) { error.value = 'Title is required'; return }
  if (!form.value.scheduled_at) { error.value = 'Scheduled date & time is required'; return }

  submitting.value = true
  try {
    const res = await commentaryApi.create({
      ...form.value,
      scheduled_at: new Date(form.value.scheduled_at).toISOString(),
    })
    const created = res.data.data
    if (created) router.push(`/commentary/${created.id}`)
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
