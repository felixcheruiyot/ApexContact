<template>
  <RouterLink
    :to="`/events/${event.id}`"
    class="bg-bg-surface rounded-2xl overflow-hidden group cursor-pointer hover:shadow-xl hover:shadow-black/40 transition-all duration-300 hover:-translate-y-1 block"
    @click="emit('click')"
  >
    <!-- Image section -->
    <div class="relative aspect-video overflow-hidden">
      <img
        :src="eventImage(event.thumbnail_url)"
        :alt="event.title"
        class="object-cover w-full h-full transition-transform duration-500 group-hover:scale-105"
        @error="onImageError"
      />

      <!-- Bottom gradient overlay -->
      <div class="absolute inset-0 bg-gradient-to-t from-black/70 to-transparent pointer-events-none" />

      <!-- Top-left: status badge -->
      <div class="absolute top-3 left-3">
        <span
          v-if="event.status === 'live'"
          class="inline-flex items-center gap-1.5 bg-accent-red text-white text-xs font-semibold px-2.5 py-1 rounded-full"
        >
          <span class="w-1.5 h-1.5 bg-white rounded-full animate-pulse" />
          LIVE
        </span>
        <span
          v-else-if="event.status === 'scheduled'"
          class="inline-flex items-center bg-white/15 backdrop-blur-sm text-white text-xs font-medium px-2.5 py-1 rounded-full border border-white/20"
        >
          Upcoming
        </span>
        <span
          v-else-if="event.status === 'completed'"
          class="inline-flex items-center bg-black/50 backdrop-blur-sm text-text-muted text-xs font-medium px-2.5 py-1 rounded-full"
        >
          Ended
        </span>
      </div>

      <!-- Top-right: event_type badge -->
      <div class="absolute top-3 right-3">
        <span class="inline-flex items-center gap-1 bg-black/60 backdrop-blur-sm text-white text-xs font-medium px-2.5 py-1 rounded-full border border-white/10">
          <Video v-if="event.event_type === 'video'" class="w-3 h-3 shrink-0" />
          <Mic v-else class="w-3 h-3 shrink-0" />
          {{ event.event_type === 'video' ? 'Video' : 'Audio' }}
        </span>
      </div>

      <!-- Bottom-left: price overlaid on gradient -->
      <div class="absolute bottom-3 left-3">
        <span class="text-white font-bold text-sm drop-shadow-lg">
          {{ event.price === 0 ? 'FREE' : `KES ${event.price.toLocaleString()}` }}
        </span>
      </div>

      <!-- Centered play button on hover -->
      <div
        class="absolute inset-0 flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity duration-300 pointer-events-none"
      >
        <div class="w-14 h-14 rounded-full bg-accent-red/90 backdrop-blur-sm flex items-center justify-center shadow-xl transform scale-90 group-hover:scale-100 transition-transform duration-300">
          <svg class="w-6 h-6 text-white ml-0.5" fill="currentColor" viewBox="0 0 24 24">
            <path d="M8 5v14l11-7z" />
          </svg>
        </div>
      </div>
    </div>

    <!-- Content section -->
    <div class="p-4">
      <!-- Sport/category badge -->
      <span class="inline-block bg-bg-elevated text-text-muted text-xs px-2 py-0.5 rounded-full capitalize">
        {{ event.sport_type }}
      </span>

      <!-- Title -->
      <h3 class="font-semibold text-white text-base leading-snug mt-2 line-clamp-2 group-hover:text-accent-red transition-colors duration-200">
        {{ event.title }}
      </h3>

      <!-- Teaser hook -->
      <p v-if="event.teaser_hook" class="text-text-muted text-sm mt-1 line-clamp-1">
        {{ event.teaser_hook }}
      </p>

      <!-- Bottom row -->
      <div class="mt-3 flex justify-between items-center">
        <!-- Left: date -->
        <div class="flex items-center gap-1.5 text-text-muted text-xs">
          <svg class="w-3 h-3 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
            <rect x="3" y="4" width="18" height="18" rx="2" ry="2" />
            <line x1="16" y1="2" x2="16" y2="6" />
            <line x1="8" y1="2" x2="8" y2="6" />
            <line x1="3" y1="10" x2="21" y2="10" />
          </svg>
          <span>{{ formattedDate }}</span>
        </div>

        <!-- Right: CTA label -->
        <span
          v-if="event.status === 'live'"
          class="text-accent-red text-sm font-medium"
        >
          Watch Now →
        </span>
        <span
          v-else-if="event.status === 'scheduled'"
          class="text-text-muted text-sm"
        >
          Buy Ticket →
        </span>
      </div>
    </div>
  </RouterLink>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { RouterLink } from 'vue-router'
import { format } from 'date-fns'
import { Video, Mic } from 'lucide-vue-next'
import type { Event } from '@/types'
import { eventImage, onImageError } from '@/utils/eventImage'

const props = defineProps<{ event: Event }>()
const emit = defineEmits<{ click: [] }>()

const formattedDate = computed(() =>
  format(new Date(props.event.scheduled_at), 'MMM d, yyyy · h:mm a')
)
</script>
