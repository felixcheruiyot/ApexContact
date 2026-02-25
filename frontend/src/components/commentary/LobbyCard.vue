<template>
  <RouterLink :to="`/commentary/${lobby.id}`" class="card group block">
    <!-- Thumbnail -->
    <div class="relative aspect-video overflow-hidden bg-bg-elevated">
      <img
        :src="eventImage(lobby.thumbnail_url)"
        :alt="lobby.title"
        class="w-full h-full object-cover transition-transform duration-500 group-hover:scale-105"
        @error="onImageError"
      />
      <!-- Status badge -->
      <div class="absolute top-3 left-3 flex items-center gap-2">
        <span v-if="lobby.status === 'live'" class="flex items-center gap-1.5 px-2 py-0.5 rounded-full bg-accent-orange text-white text-xs font-bold uppercase tracking-wider">
          <span class="w-1.5 h-1.5 bg-white rounded-full animate-pulse" />
          Live Room
        </span>
        <span v-else-if="lobby.status === 'scheduled'" class="badge-upcoming">
          Upcoming
        </span>
        <span v-else-if="lobby.status === 'completed'" class="px-2 py-0.5 rounded-full bg-white/10 text-text-muted text-xs font-medium">
          Ended
        </span>
      </div>
      <!-- Sport badge -->
      <div class="absolute top-3 right-3">
        <span class="bg-black/50 backdrop-blur-sm text-white text-xs px-2 py-1 rounded-full uppercase tracking-wider">
          {{ lobby.sport_type }}
        </span>
      </div>
      <!-- Hover overlay with mic icon -->
      <div class="absolute inset-0 bg-black/40 opacity-0 group-hover:opacity-100 transition-opacity
                  flex items-center justify-center">
        <div class="w-14 h-14 rounded-full bg-accent-orange flex items-center justify-center
                    transform scale-90 group-hover:scale-100 transition-transform">
          <svg class="w-6 h-6 text-white" fill="currentColor" viewBox="0 0 24 24">
            <path d="M12 15c1.66 0 3-1.34 3-3V6c0-1.66-1.34-3-3-3S9 4.34 9 6v6c0 1.66 1.34 3 3 3zm-1-9c0-.55.45-1 1-1s1 .45 1 1v6c0 .55-.45 1-1 1s-1-.45-1-1V6z"/>
            <path d="M17 12c0 2.76-2.24 5-5 5s-5-2.24-5-5H5c0 3.53 2.61 6.43 6 6.92V22h2v-3.08c3.39-.49 6-3.39 6-6.92h-2z"/>
          </svg>
        </div>
      </div>
    </div>

    <!-- Info -->
    <div class="p-4">
      <!-- Commentary type badge -->
      <div class="flex items-center gap-2 mb-2">
        <span class="text-accent-orange text-xs font-semibold uppercase tracking-wider flex items-center gap-1">
          <Mic class="w-3 h-3" /> Live Room
        </span>
      </div>
      <h3 class="text-white font-semibold text-sm leading-tight line-clamp-2 mb-1 group-hover:text-accent-orange
                 transition-colors">
        {{ lobby.title }}
      </h3>
      <p v-if="lobby.teaser_hook" class="text-text-muted text-xs line-clamp-1 mb-2">
        {{ lobby.teaser_hook }}
      </p>
      <div class="flex items-center justify-between">
        <span class="text-text-muted text-xs">{{ formattedDate }}</span>
        <span class="font-bold text-sm" :class="lobby.price === 0 ? 'text-success' : 'text-accent-orange'">
          {{ lobby.price === 0 ? 'Free' : `${lobby.currency} ${lobby.price.toLocaleString()}` }}
        </span>
      </div>
    </div>
  </RouterLink>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { RouterLink } from 'vue-router'
import { Mic } from 'lucide-vue-next'
import { format } from 'date-fns'
import type { Event } from '@/types'
import { eventImage, onImageError } from '@/utils/eventImage'

const props = defineProps<{ lobby: Event }>()

const formattedDate = computed(() =>
  format(new Date(props.lobby.scheduled_at), 'MMM d, yyyy · h:mm a')
)
</script>
