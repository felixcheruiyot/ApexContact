<template>
  <div v-if="event" class="relative min-h-[70vh] flex items-end overflow-hidden">
    <!-- Background image -->
    <div class="absolute inset-0">
      <img
        :src="event.thumbnail_url || '/placeholder-event.jpg'"
        :alt="event.title"
        class="w-full h-full object-cover"
      />
      <div class="absolute inset-0 bg-gradient-to-t from-bg via-bg/60 to-transparent" />
      <div class="absolute inset-0 bg-gradient-to-r from-bg/80 to-transparent" />
    </div>

    <!-- Content -->
    <div class="relative z-10 max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 pb-16 w-full">
      <div class="max-w-2xl animate-slide-up">
        <!-- Badges -->
        <div class="flex items-center gap-3 mb-4">
          <span v-if="event.status === 'live'" class="badge-live">
            <span class="w-1.5 h-1.5 bg-white rounded-full animate-pulse" />
            Live Now
          </span>
          <span class="bg-bg-elevated/80 backdrop-blur-sm border border-white/10 text-text-muted
                       text-xs font-semibold uppercase tracking-wider px-3 py-1 rounded-full">
            {{ event.sport_type }}
          </span>
        </div>

        <!-- Title -->
        <h1 class="font-display text-5xl md:text-7xl uppercase tracking-wide text-white mb-4 leading-none">
          {{ event.title }}
        </h1>

        <!-- Description -->
        <p class="text-text-muted text-base leading-relaxed mb-6 max-w-xl line-clamp-3">
          {{ event.description }}
        </p>

        <!-- Date + Price -->
        <div class="flex items-center gap-6 mb-8">
          <div>
            <p class="text-text-muted text-xs uppercase tracking-wider mb-1">Date</p>
            <p class="text-white font-semibold">{{ formattedDate }}</p>
          </div>
          <div class="w-px h-10 bg-white/10" />
          <div>
            <p class="text-text-muted text-xs uppercase tracking-wider mb-1">Price</p>
            <p class="text-accent-red font-bold text-xl">
              {{ event.currency }} {{ event.price.toLocaleString() }}
            </p>
          </div>
        </div>

        <!-- CTA -->
        <div class="flex items-center gap-4">
          <button v-if="event.status === 'live'" @click="$emit('buy-ticket')" class="btn-primary text-base px-8 py-4">
            Watch Live
          </button>
          <button v-else-if="event.status === 'scheduled'" @click="$emit('buy-ticket')" class="btn-primary text-base px-8 py-4">
            Buy Ticket
          </button>
          <CountdownTimer v-if="event.status === 'scheduled'" :target="event.scheduled_at" />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { format } from 'date-fns'
import type { Event } from '@/types'
import CountdownTimer from './CountdownTimer.vue'

const props = defineProps<{ event: Event }>()
defineEmits<{ 'buy-ticket': [] }>()

const formattedDate = computed(() =>
  format(new Date(props.event.scheduled_at), 'EEEE, MMMM d, yyyy · h:mm a')
)
</script>
