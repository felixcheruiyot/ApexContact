<template>
  <RouterLink :to="to"
    class="flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm font-medium transition-colors"
    :class="isActive
      ? 'bg-accent-red/10 text-accent-red'
      : 'text-text-muted hover:text-white hover:bg-bg-elevated'"
  >
    <!-- SVG icon set -->
    <svg class="w-4 h-4 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
      <!-- grid / overview -->
      <path v-if="icon === 'grid'" stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
        d="M4 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2V6zM14 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2V6zM4 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2v-2zM14 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2v-2z" />
      <!-- users -->
      <path v-else-if="icon === 'users'" stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
        d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" />
      <!-- film / events -->
      <path v-else-if="icon === 'film'" stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
        d="M7 4v16M17 4v16M3 8h4m10 0h4M3 12h18M3 16h4m10 0h4M4 20h16a1 1 0 001-1V5a1 1 0 00-1-1H4a1 1 0 00-1 1v14a1 1 0 001 1z" />
      <!-- shield / fraud -->
      <path v-else-if="icon === 'shield'" stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
        d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z" />
      <!-- plus / create -->
      <path v-else-if="icon === 'plus'" stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
        d="M12 4v16m8-8H4" />
      <!-- bar-chart / analytics -->
      <path v-else-if="icon === 'bar-chart'" stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
        d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
      <!-- radio / stream -->
      <path v-else-if="icon === 'radio'" stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
        d="M8.111 16.404a5.5 5.5 0 017.778 0M12 20h.01m-7.08-7.071c3.904-3.905 10.236-3.905 14.141 0M1.394 9.393c5.857-5.857 15.355-5.857 21.213 0" />
      <!-- user / profile -->
      <path v-else-if="icon === 'user'" stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
        d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
      <!-- fallback circle -->
      <circle v-else cx="12" cy="12" r="4" stroke-width="2" />
    </svg>
    <span>{{ label }}</span>
  </RouterLink>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { RouterLink, useRoute } from 'vue-router'

const props = defineProps<{ to: string; icon: string; label: string; exact?: boolean }>()
const route = useRoute()
const isActive = computed(() =>
  props.exact ? route.path === props.to : route.path.startsWith(props.to),
)
</script>
