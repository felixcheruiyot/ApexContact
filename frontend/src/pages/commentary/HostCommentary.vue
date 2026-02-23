<template>
  <div class="max-w-2xl mx-auto px-4 sm:px-6 py-10">
    <div class="mb-8">
      <RouterLink to="/" class="flex items-center gap-1.5 text-text-muted hover:text-white text-sm transition-colors"><ArrowLeft class="w-4 h-4" /> Back to home</RouterLink>
      <h1 class="text-2xl font-bold text-white mt-4">Start a Live Room</h1>
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
          placeholder="e.g. Nairobi Career Q&A — Ask Me Anything"
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
          placeholder="What will you discuss? Who should join?"
          class="form-input resize-none"
        />
      </div>

      <!-- Category -->
      <div>
        <label class="form-label">Category *</label>
        <select v-model="form.sport_type" class="form-input appearance-none" required>
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
            <template v-if="form.price === 0"><Check class="w-4 h-4 inline" /> Free entry</template>
            <template v-else>Listeners pay KES {{ form.price }}</template>
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
        <Mic v-if="!submitting" class="w-4 h-4" />{{ submitting ? 'Creating…' : 'Start Live Room' }}
      </button>
    </form>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { RouterLink, useRouter } from 'vue-router'
import { ArrowLeft, Check, Mic } from 'lucide-vue-next'
import { commentaryApi } from '@/api/commentary'

const router = useRouter()

const form = ref({
  title: '',
  teaser_hook: '',
  description: '',
  sport_type: '',
  scheduled_at: '',
  price: 0,
  thumbnail_url: '',
})

const submitting = ref(false)
const error = ref('')

async function submit() {
  error.value = ''
  if (!form.value.title.trim()) { error.value = 'Title is required'; return }
  if (!form.value.sport_type) { error.value = 'Please select a category'; return }
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
