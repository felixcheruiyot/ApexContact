<template>
  <RouterLink :to="`/events/${event.id}`" class="card group block">
    <!-- Thumbnail -->
    <div class="relative aspect-video overflow-hidden bg-bg-elevated">
      <img
        :src="event.thumbnail_url || '/placeholder-event.jpg'"
        :alt="event.title"
        class="w-full h-full object-cover transition-transform duration-500 group-hover:scale-105"
      />
      <!-- Status badge -->
      <div class="absolute top-3 left-3">
        <span v-if="event.status === 'live'" class="badge-live">
          <span class="w-1.5 h-1.5 bg-white rounded-full animate-pulse" />
          Live
        </span>
        <span v-else-if="event.status === 'scheduled'" class="badge-upcoming">
          Upcoming
        </span>
      </div>
      <!-- Sport badge -->
      <div class="absolute top-3 right-3">
        <span class="bg-black/50 backdrop-blur-sm text-white text-xs px-2 py-1 rounded-full uppercase tracking-wider">
          {{ event.sport_type }}
        </span>
      </div>
      <!-- Hover overlay with play icon -->
      <div class="absolute inset-0 bg-black/40 opacity-0 group-hover:opacity-100 transition-opacity
                  flex items-center justify-center">
        <div class="w-14 h-14 rounded-full bg-accent-red flex items-center justify-center
                    transform scale-90 group-hover:scale-100 transition-transform">
          <svg class="w-6 h-6 text-white ml-1" fill="currentColor" viewBox="0 0 24 24">
            <path d="M8 5v14l11-7z" />
          </svg>
        </div>
      </div>
    </div>

    <!-- Info -->
    <div class="p-4">
      <h3 class="text-white font-semibold text-sm leading-tight line-clamp-2 mb-2 group-hover:text-accent-red
                 transition-colors">
        {{ event.title }}
      </h3>
      <div class="flex items-center justify-between">
        <span class="text-text-muted text-xs">{{ formattedDate }}</span>
        <span class="text-accent-red font-bold text-sm">
          {{ event.currency }} {{ event.price.toLocaleString() }}
        </span>
      </div>
    </div>
  </RouterLink>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { RouterLink } from 'vue-router'
import { format } from 'date-fns'
import type { Event } from '@/types'

const props = defineProps<{ event: Event }>()

const formattedDate = computed(() =>
  format(new Date(props.event.scheduled_at), 'MMM d, yyyy · h:mm a')
)
</script>
