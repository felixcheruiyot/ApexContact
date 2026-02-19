<template>
  <div class="fixed inset-0 z-50 flex items-center justify-center p-4">
    <div class="absolute inset-0 bg-black/70 backdrop-blur-sm" @click="$emit('close')" />

    <div class="relative z-10 bg-bg-surface rounded-2xl p-6 w-full max-w-md border border-white/10">
      <h2 class="text-xl font-bold text-white mb-1">Choose your nickname</h2>
      <p class="text-text-muted text-sm mb-6">Your real name won't be visible to other listeners.</p>

      <!-- Suggested nicknames -->
      <div v-if="suggestions.length" class="mb-4">
        <p class="text-text-muted text-xs uppercase tracking-wider mb-3">Suggestions</p>
        <div class="flex flex-wrap gap-2">
          <button
            v-for="s in suggestions"
            :key="s"
            @click="selectSuggestion(s)"
            class="px-3 py-1.5 rounded-full text-sm font-medium border transition-all duration-150"
            :class="nickname === s
              ? 'bg-accent-orange border-accent-orange text-white shadow-[0_0_12px_rgba(255,107,0,0.4)]'
              : 'bg-bg-elevated border-white/10 text-text-muted hover:border-accent-orange hover:text-white'"
          >
            {{ s }}
          </button>
        </div>
      </div>

      <!-- Skeleton suggestions loading -->
      <div v-else-if="loadingSuggestions" class="flex gap-2 mb-4">
        <div v-for="i in 4" :key="i" class="h-8 w-28 bg-bg-elevated rounded-full animate-pulse" />
      </div>

      <!-- Custom input -->
      <div class="mb-6">
        <label class="block text-text-muted text-xs uppercase tracking-wider mb-2">Or type your own</label>
        <input
          v-model="nickname"
          type="text"
          maxlength="32"
          placeholder="e.g. RingFan42"
          class="w-full bg-bg-elevated border border-white/10 rounded-lg px-4 py-2.5 text-white
                 placeholder-text-muted focus:outline-none focus:border-accent-orange transition-colors"
          @keyup.enter="confirm"
        />
      </div>

      <!-- Actions -->
      <div class="flex gap-3">
        <button
          @click="$emit('close')"
          class="flex-1 py-2.5 rounded-lg border border-white/10 text-text-muted hover:text-white
                 hover:border-white/30 transition-colors text-sm font-medium"
        >
          Cancel
        </button>
        <button
          @click="confirm"
          :disabled="!nickname.trim() || confirming"
          class="flex-1 py-2.5 rounded-lg bg-accent-orange text-white font-semibold text-sm
                 disabled:opacity-40 disabled:cursor-not-allowed hover:bg-orange-500 transition-colors"
        >
          {{ confirming ? 'Joining…' : (isFree ? 'Join Free' : 'Continue to Pay') }}
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { commentaryApi } from '@/api/commentary'

const props = defineProps<{
  eventId: string
  isFree: boolean
}>()

const emit = defineEmits<{
  (e: 'close'): void
  (e: 'confirm', nickname: string): void
}>()

const nickname = ref('')
const suggestions = ref<string[]>([])
const loadingSuggestions = ref(true)
const confirming = ref(false)

onMounted(async () => {
  try {
    const res = await commentaryApi.suggestNicknames(props.eventId)
    suggestions.value = res.data.data ?? []
  } catch {
    // Suggestions are optional
  } finally {
    loadingSuggestions.value = false
  }
})

function selectSuggestion(s: string) {
  nickname.value = s
}

async function confirm() {
  const name = nickname.value.trim()
  if (!name) return
  confirming.value = true
  try {
    emit('confirm', name)
  } finally {
    confirming.value = false
  }
}
</script>
