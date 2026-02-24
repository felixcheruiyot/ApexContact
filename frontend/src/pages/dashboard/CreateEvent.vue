<template>
  <div class="max-w-3xl">

    <!-- ── Step indicator ─────────────────────────────────────────── -->
    <div class="flex items-center gap-3 mb-10">
      <div
        v-for="(label, i) in ['Broadcast type', 'Event details']"
        :key="label"
        class="flex items-center gap-2"
      >
        <div
          class="w-7 h-7 rounded-full flex items-center justify-center text-xs font-bold transition-colors"
          :class="step > i
            ? 'bg-status-success text-white'
            : step === i
              ? 'bg-accent-red text-white'
              : 'bg-white/10 text-text-muted'"
        >
          <svg v-if="step > i" class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="3" d="M5 13l4 4L19 7" />
          </svg>
          <span v-else>{{ i + 1 }}</span>
        </div>
        <span class="text-sm" :class="step === i ? 'text-white font-semibold' : 'text-text-muted'">
          {{ label }}
        </span>
        <svg v-if="i < 1" class="w-4 h-4 text-white/20" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
        </svg>
      </div>
    </div>

    <!-- ══════════════════════════════════════════════════════════════ -->
    <!-- STEP 0 — Choose broadcast type                                -->
    <!-- ══════════════════════════════════════════════════════════════ -->
    <div v-if="step === 0" class="space-y-6">
      <div>
        <h1 class="text-white font-bold text-2xl mb-1">How do you want to broadcast?</h1>
        <p class="text-text-muted text-sm">Choose the format that best fits your session.</p>
      </div>

      <div class="grid grid-cols-1 gap-4">
        <button
          v-for="mode in broadcastModes"
          :key="mode.type"
          @click="selectType(mode.type)"
          class="relative w-full text-left rounded-xl border p-5 transition-all duration-200 group"
          :class="selectedType === mode.type
            ? 'border-accent-red bg-accent-red/5'
            : 'border-white/10 bg-bg-elevated hover:border-white/25 hover:bg-white/3'"
        >
          <div class="flex items-start gap-4">
            <!-- Icon -->
            <div
              class="w-12 h-12 rounded-lg flex items-center justify-center shrink-0 transition-colors"
              :class="selectedType === mode.type ? mode.iconBg + ' opacity-100' : mode.iconBg + ' opacity-60 group-hover:opacity-80'"
            >
              <component :is="mode.icon" class="w-6 h-6" :class="mode.iconColor" />
            </div>

            <!-- Text -->
            <div class="flex-1 min-w-0">
              <div class="flex items-center gap-2 mb-1">
                <span class="text-white font-semibold text-base">{{ mode.label }}</span>
                <span
                  v-if="mode.badge"
                  class="text-xs font-semibold px-2 py-0.5 rounded-full"
                  :class="mode.badgeClass"
                >{{ mode.badge }}</span>
              </div>
              <p class="text-text-muted text-sm leading-relaxed">{{ mode.description }}</p>
              <div class="flex flex-wrap gap-3 mt-3">
                <span
                  v-for="tag in mode.tags"
                  :key="tag"
                  class="text-xs text-text-muted bg-white/5 border border-white/8 rounded px-2 py-0.5"
                >{{ tag }}</span>
              </div>
            </div>

            <!-- Selection indicator -->
            <div
              class="w-5 h-5 rounded-full border-2 flex items-center justify-center shrink-0 mt-0.5 transition-colors"
              :class="selectedType === mode.type ? 'border-accent-red bg-accent-red' : 'border-white/20'"
            >
              <svg v-if="selectedType === mode.type" class="w-3 h-3 text-white" fill="currentColor" viewBox="0 0 24 24">
                <circle cx="12" cy="12" r="6" />
              </svg>
            </div>
          </div>
        </button>
      </div>

      <div class="flex items-center gap-4 pt-2">
        <button
          @click="step = 1"
          :disabled="!selectedType"
          class="btn-primary disabled:opacity-40 disabled:cursor-not-allowed"
        >
          Continue
          <svg class="w-4 h-4 ml-1 inline-block" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
          </svg>
        </button>
        <RouterLink to="/dashboard" class="btn-ghost">Cancel</RouterLink>
      </div>
    </div>

    <!-- ══════════════════════════════════════════════════════════════ -->
    <!-- STEP 1 — Event details                                        -->
    <!-- ══════════════════════════════════════════════════════════════ -->
    <div v-if="step === 1" class="space-y-6">
      <div class="flex items-center gap-3 mb-2">
        <button
          @click="step = 0"
          class="text-text-muted hover:text-white transition-colors"
          aria-label="Back"
        >
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
          </svg>
        </button>
        <div>
          <h1 class="text-white font-bold text-2xl leading-none">Event details</h1>
          <p class="text-text-muted text-xs mt-1">
            {{ selectedMode?.label }} —
            <button @click="step = 0" class="underline hover:text-white transition-colors">Change</button>
          </p>
        </div>
      </div>

      <!-- Selected type pill -->
      <div class="flex items-center gap-2 p-3 rounded-lg bg-white/3 border border-white/8">
        <component :is="selectedMode?.icon" class="w-4 h-4" :class="selectedMode?.iconColor" />
        <span class="text-white text-sm font-medium">{{ selectedMode?.label }}</span>
        <span class="text-text-muted text-xs ml-1">— {{ selectedMode?.shortDesc }}</span>
      </div>

      <form @submit.prevent="handleSubmit" class="space-y-5">
        <!-- Title -->
        <div>
          <label class="form-label">Event Title *</label>
          <input
            v-model="form.title"
            type="text"
            :placeholder="selectedMode?.titlePlaceholder"
            class="input"
            required
          />
        </div>

        <!-- Teaser hook (always shown for interactive, optional for video) -->
        <div>
          <label class="form-label">
            Teaser Hook
            <span class="text-text-muted font-normal ml-1">(optional)</span>
          </label>
          <input
            v-model="form.teaser_hook"
            type="text"
            maxlength="150"
            placeholder="One-liner to get your audience excited"
            class="input"
          />
        </div>

        <!-- Description -->
        <div>
          <label class="form-label">Description</label>
          <textarea
            v-model="form.description"
            rows="3"
            placeholder="What will happen? Who should attend?"
            class="input resize-none"
          />
        </div>

        <!-- Category + Date -->
        <div class="grid grid-cols-2 gap-4">
          <div>
            <label class="form-label">Category *</label>
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
            <label class="form-label">Date & Time *</label>
            <input v-model="form.scheduled_at" type="datetime-local" class="input" required />
          </div>
        </div>

        <!-- Ticket price -->
        <div>
          <label class="form-label">Ticket Price (KES)</label>
          <div class="flex items-center gap-3">
            <input
              v-model.number="form.price"
              type="number"
              min="0"
              step="1"
              placeholder="0"
              class="input"
            />
            <span class="text-text-muted text-sm whitespace-nowrap shrink-0">0 = Free</span>
          </div>
        </div>

        <!-- Thumbnail URL -->
        <div>
          <label class="form-label">Thumbnail URL <span class="text-text-muted font-normal">(optional)</span></label>
          <input v-model="form.thumbnail_url" type="url" placeholder="https://…" class="input" />
          <div v-if="form.thumbnail_url" class="mt-3 rounded-lg overflow-hidden aspect-video max-w-xs">
            <img :src="form.thumbnail_url" class="w-full h-full object-cover" alt="Thumbnail preview" />
          </div>
        </div>

        <!-- Error -->
        <div v-if="errorMsg" class="bg-status-error/10 border border-status-error/30 text-status-error
                                    text-sm rounded-lg px-4 py-3">
          {{ errorMsg }}
        </div>

        <!-- Actions -->
        <div class="flex items-center gap-4 pt-2">
          <button type="submit" class="btn-primary" :disabled="loading">
            <svg v-if="loading" class="w-4 h-4 animate-spin inline-block mr-2" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"/>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"/>
            </svg>
            {{ loading ? 'Creating…' : 'Create Event' }}
          </button>
          <RouterLink to="/dashboard" class="btn-ghost">Cancel</RouterLink>
        </div>
      </form>
    </div>

  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { RouterLink, useRouter } from 'vue-router'
import { Mic, Video, Monitor } from 'lucide-vue-next'
import { eventsApi } from '@/api/events'
import { commentaryApi } from '@/api/commentary'
import type { SportType } from '@/types'

type BroadcastType = 'audio' | 'audio_video' | 'commercial'

const router = useRouter()
const step = ref(0)
const selectedType = ref<BroadcastType | null>(null)
const loading = ref(false)
const errorMsg = ref('')

const broadcastModes = [
  {
    type: 'audio_video' as BroadcastType,
    label: 'Phone / PC — Audio + Video',
    shortDesc: 'Interactive LiveKit room with camera',
    description: 'Stream directly from your phone or browser camera. Viewers can join as listeners and you can invite speakers on stage. Great for Q&As, masterclasses, and interactive sessions.',
    badge: 'Interactive',
    badgeClass: 'bg-accent-orange/20 text-accent-orange',
    icon: Video,
    iconBg: 'bg-accent-orange/15',
    iconColor: 'text-accent-orange',
    titlePlaceholder: 'e.g. Nairobi Career Q&A — Ask Me Anything',
    tags: ['No OBS needed', 'Browser / Phone', 'Two-way audio', 'LiveKit'],
  },
  {
    type: 'audio' as BroadcastType,
    label: 'Phone / PC — Audio Only',
    shortDesc: 'Interactive LiveKit room, mic only',
    description: 'Mic-only interactive session — like a podcast or Twitter Space. No camera required. Perfect for discussions, panels, and audio-first content.',
    badge: 'Audio room',
    badgeClass: 'bg-accent-red/20 text-accent-red',
    icon: Mic,
    iconBg: 'bg-accent-red/15',
    iconColor: 'text-accent-red',
    titlePlaceholder: 'e.g. Weekly Business Panel — Episode 12',
    tags: ['Mic only', 'No camera needed', 'Browser / Phone', 'LiveKit'],
  },
  {
    type: 'commercial' as BroadcastType,
    label: 'OBS / RTMP Encoder',
    shortDesc: 'Professional RTMP → HLS one-to-many',
    description: 'Use OBS Studio, Streamlabs, or any RTMP encoder. Full-quality one-to-many broadcast with anti-piracy protection. Best for high-production events and large audiences.',
    badge: 'Professional',
    badgeClass: 'bg-white/10 text-text-muted',
    icon: Monitor,
    iconBg: 'bg-white/10',
    iconColor: 'text-text-muted',
    titlePlaceholder: 'e.g. Nairobi Boxing Championship 2026',
    tags: ['OBS / Streamlabs', 'RTMP → HLS', 'Anti-piracy', 'Unlimited viewers'],
  },
]

const selectedMode = computed(() =>
  broadcastModes.find(m => m.type === selectedType.value) ?? null
)

function selectType(type: BroadcastType) {
  selectedType.value = type
}

const form = ref({
  title: '',
  teaser_hook: '',
  description: '',
  sport_type: '' as SportType | '',
  scheduled_at: '',
  price: 0,
  currency: 'KES',
  thumbnail_url: '',
})

async function handleSubmit() {
  if (!form.value.sport_type || !selectedType.value) return
  errorMsg.value = ''
  loading.value = true

  try {
    const scheduledAt = new Date(form.value.scheduled_at).toISOString()
    const type = selectedType.value

    if (type === 'commercial') {
      // OBS/RTMP → goes through standard events API (needs review)
      await eventsApi.create({
        title: form.value.title,
        description: form.value.description,
        teaser_hook: form.value.teaser_hook,
        sport_type: form.value.sport_type as SportType,
        scheduled_at: scheduledAt,
        price: form.value.price,
        currency: form.value.currency,
        thumbnail_url: form.value.thumbnail_url,
        event_type: 'video',
      })
      router.push('/dashboard')
    } else {
      // audio or audio_video → commentary API (no review, goes live immediately)
      const res = await commentaryApi.create({
        title: form.value.title,
        teaser_hook: form.value.teaser_hook,
        description: form.value.description,
        sport_type: form.value.sport_type as string,
        scheduled_at: scheduledAt,
        price: form.value.price,
        thumbnail_url: form.value.thumbnail_url,
        event_type: type === 'audio' ? 'audio' : 'audio_video',
      })
      const created = res.data.data
      if (created) {
        router.push(`/commentary/${created.id}`)
      } else {
        router.push('/dashboard')
      }
    }
  } catch (e: any) {
    errorMsg.value = e.response?.data?.error ?? 'Failed to create event'
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.form-label {
  @apply block text-text-muted text-xs uppercase tracking-wider mb-2 font-semibold;
}
</style>
