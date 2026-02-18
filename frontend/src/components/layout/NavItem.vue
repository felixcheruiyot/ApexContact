<template>
  <RouterLink :to="to"
    class="flex items-center gap-3 px-3 py-2.5 rounded-lg text-sm font-medium transition-colors"
    :class="isActive
      ? 'bg-accent-red/10 text-accent-red'
      : 'text-text-muted hover:text-white hover:bg-bg-elevated'"
  >
    <span>{{ icon }}</span>
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
