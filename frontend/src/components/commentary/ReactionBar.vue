<template>
  <div class="flex items-center justify-center gap-2 py-3 px-4 bg-bg-surface/80 backdrop-blur-sm rounded-xl border border-white/5">
    <button
      v-for="emoji in emojis"
      :key="emoji"
      @click="sendReaction(emoji)"
      class="text-xl hover:scale-125 active:scale-95 transition-transform duration-150 select-none
             cursor-pointer p-1.5 rounded-lg hover:bg-white/5"
      :title="`React with ${emoji}`"
    >
      {{ emoji }}
    </button>
  </div>

  <!-- Floating reactions overlay -->
  <Teleport to="body">
    <div class="pointer-events-none fixed inset-0 z-[100] overflow-hidden">
      <TransitionGroup name="float">
        <div
          v-for="r in activeReactions"
          :key="r.id"
          class="absolute text-3xl select-none"
          :style="r.style"
        >
          {{ r.emoji }}
        </div>
      </TransitionGroup>
    </div>
  </Teleport>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const emit = defineEmits<{
  (e: 'reaction', emoji: string): void
}>()

const emojis = ['🔥', '💯', '👊', '🎯', '💡', '😂', '🤯', '👀']

interface FloatingReaction {
  id: number
  emoji: string
  style: Record<string, string>
}

const activeReactions = ref<FloatingReaction[]>([])
let nextId = 0

function sendReaction(emoji: string) {
  emit('reaction', emoji)
  spawnFloat(emoji)
}

function spawnFloat(emoji: string) {
  const id = nextId++
  const left = 10 + Math.random() * 80  // 10%–90% horizontal
  const reaction: FloatingReaction = {
    id,
    emoji,
    style: {
      left: `${left}%`,
      bottom: '80px',
      animation: 'float-up 2s ease-out forwards',
    },
  }
  activeReactions.value.push(reaction)
  setTimeout(() => {
    activeReactions.value = activeReactions.value.filter((r) => r.id !== id)
  }, 2200)
}

// Expose for external reactions (e.g. received via WebSocket)
defineExpose({ spawnFloat })
</script>

<style scoped>
.float-enter-active {
  animation: float-up 2s ease-out forwards;
}
.float-leave-active {
  opacity: 0;
}

@keyframes float-up {
  0%   { transform: translateY(0) scale(1);   opacity: 1; }
  80%  { transform: translateY(-200px) scale(1.3); opacity: 0.8; }
  100% { transform: translateY(-280px) scale(0.8); opacity: 0; }
}
</style>
