<template>
  <div v-if="!expired" class="flex items-center gap-4">
    <div v-for="unit in units" :key="unit.label" class="text-center">
      <div class="bg-bg-elevated border border-white/10 rounded-lg px-3 py-2 min-w-[52px]">
        <span class="font-display text-2xl text-white">{{ unit.value }}</span>
      </div>
      <p class="text-text-muted text-xs uppercase tracking-wider mt-1">{{ unit.label }}</p>
    </div>
  </div>
  <span v-else class="badge-live">
    <span class="w-1.5 h-1.5 bg-white rounded-full animate-pulse" />
    Happening Now
  </span>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'

const props = defineProps<{ target: string }>()

const now = ref(Date.now())
let timer: ReturnType<typeof setInterval>

onMounted(() => { timer = setInterval(() => { now.value = Date.now() }, 1000) })
onUnmounted(() => clearInterval(timer))

const diff = computed(() => Math.max(0, new Date(props.target).getTime() - now.value))
const expired = computed(() => diff.value === 0)

const units = computed(() => {
  const s = Math.floor(diff.value / 1000)
  return [
    { label: 'Days',    value: String(Math.floor(s / 86400)).padStart(2, '0') },
    { label: 'Hours',   value: String(Math.floor((s % 86400) / 3600)).padStart(2, '0') },
    { label: 'Minutes', value: String(Math.floor((s % 3600) / 60)).padStart(2, '0') },
    { label: 'Seconds', value: String(s % 60).padStart(2, '0') },
  ]
})
</script>
