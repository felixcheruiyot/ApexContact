<template>
  <div v-if="event" class="relative min-h-[70vh] flex items-end overflow-hidden">
    <!-- Background image -->
    <img
      :src="eventImage(event.thumbnail_url)"
      :alt="event.title"
      class="absolute inset-0 w-full h-full object-cover"
      @error="onImageError"
    />

    <!-- Bottom fade overlay -->
    <div class="absolute inset-0 bg-gradient-to-t from-black via-black/60 to-transparent pointer-events-none" />
    <!-- Left fade overlay -->
    <div class="absolute inset-0 bg-gradient-to-r from-black/80 to-transparent pointer-events-none" />

    <!-- Content -->
    <div class="relative z-10 p-8 md:p-16 max-w-3xl w-full">
      <!-- Badge row -->
      <div class="flex flex-wrap items-center gap-2 mb-5">
        <span class="inline-flex items-center bg-white/10 backdrop-blur-sm border border-white/20 text-white text-xs font-semibold uppercase tracking-wider px-3 py-1 rounded-full capitalize">
          {{ event.sport_type }}
        </span>
        <span
          v-if="event.status === 'live'"
          class="inline-flex items-center gap-1.5 bg-accent-red text-white text-xs font-semibold px-3 py-1 rounded-full"
        >
          <span class="w-1.5 h-1.5 bg-white rounded-full animate-pulse" />
          Live Now
        </span>
        <span
          v-else-if="event.status === 'scheduled'"
          class="inline-flex items-center bg-white/10 backdrop-blur-sm border border-white/20 text-white text-xs font-medium px-3 py-1 rounded-full"
        >
          Upcoming
        </span>
        <span
          v-else-if="event.status === 'completed'"
          class="inline-flex items-center bg-black/40 backdrop-blur-sm border border-white/10 text-text-muted text-xs font-medium px-3 py-1 rounded-full"
        >
          Ended
        </span>

        <!-- Event type badge -->
        <span class="inline-flex items-center gap-1 bg-white/10 backdrop-blur-sm border border-white/10 text-white text-xs px-3 py-1 rounded-full">
          <Video v-if="event.event_type === 'video'" class="w-3.5 h-3.5" />
          <Mic v-else class="w-3.5 h-3.5" />
          {{ event.event_type === 'video' ? 'Video' : 'Audio' }}
        </span>
      </div>

      <!-- Title -->
      <h1 class="font-display text-5xl md:text-7xl uppercase text-white leading-none tracking-wide">
        {{ event.title }}
      </h1>

      <!-- Teaser hook -->
      <p v-if="event.teaser_hook" class="text-text-muted text-xl mt-4 max-w-xl leading-relaxed">
        {{ event.teaser_hook }}
      </p>

      <!-- Meta row -->
      <div class="text-text-muted text-sm flex flex-wrap gap-x-6 gap-y-2 mt-6 items-center">
        <span class="flex items-center gap-1.5">
          <Calendar class="w-4 h-4 shrink-0" />
          {{ formattedDate }}
        </span>
        <span class="flex items-center gap-1.5">
          <DollarSign class="w-4 h-4 shrink-0" />
          {{ event.price === 0 ? 'Free' : `${event.currency} ${event.price.toLocaleString()}` }}
        </span>
        <span class="flex items-center gap-1.5">
          <Tag class="w-4 h-4 shrink-0" />
          <span class="capitalize">{{ event.sport_type }}</span>
        </span>
      </div>

      <!-- CTA buttons row -->
      <div class="mt-8 flex flex-wrap items-center gap-4">
        <!-- Live + has token: Watch Live -->
        <button
          v-if="event.status === 'live' && hasToken"
          class="inline-flex items-center gap-2 bg-accent-red hover:bg-accent-red-hover text-white font-bold text-lg px-8 py-4 rounded-full transition-colors duration-200 shadow-lg shadow-accent-red/30"
        >
          <Play class="w-5 h-5 fill-current" />
          Watch Live
        </button>

        <!-- Live + no token: Buy Ticket -->
        <button
          v-else-if="event.status === 'live' && !hasToken"
          @click="onBuyTicket"
          class="inline-flex items-center gap-2 bg-accent-red hover:bg-accent-red-hover text-white font-bold text-lg px-8 py-4 rounded-full transition-colors duration-200 shadow-lg shadow-accent-red/30"
        >
          Buy Ticket · {{ event.currency }} {{ event.price.toLocaleString() }}
        </button>

        <!-- Scheduled: Reserve + Notify Me -->
        <template v-else-if="event.status === 'scheduled'">
          <button
            @click="onBuyTicket"
            class="inline-flex items-center gap-2 bg-accent-red hover:bg-accent-red-hover text-white font-bold text-base px-8 py-4 rounded-full transition-colors duration-200 shadow-lg shadow-accent-red/30"
          >
            <Ticket class="w-5 h-5 shrink-0" /> Reserve Your Spot · {{ event.currency }} {{ event.price.toLocaleString() }}
          </button>
          <button
            class="inline-flex items-center gap-2 border-2 border-white/30 hover:border-white/60 text-white font-medium text-base px-6 py-4 rounded-full transition-colors duration-200 backdrop-blur-sm"
          >
            Notify Me
          </button>
        </template>

        <!-- Completed -->
        <button
          v-else-if="event.status === 'completed'"
          disabled
          class="inline-flex items-center gap-2 border-2 border-white/20 text-text-muted font-medium text-base px-8 py-4 rounded-full cursor-not-allowed opacity-60"
        >
          Event Ended
        </button>
      </div>

      <!-- Countdown timer (only for scheduled events) -->
      <div v-if="event.status === 'scheduled' && !countdownExpired" class="mt-6">
        <p class="text-text-muted text-xs uppercase tracking-widest mb-2 font-semibold">Starts in</p>
        <div class="flex items-center gap-3">
          <div
            v-for="unit in countdownUnits"
            :key="unit.label"
            class="text-center"
          >
            <div class="bg-white/10 backdrop-blur-sm border border-white/20 rounded-xl px-4 py-2 min-w-[56px]">
              <span class="font-display text-3xl text-white leading-none">{{ unit.value }}</span>
            </div>
            <p class="text-text-muted text-xs uppercase tracking-wider mt-1.5">{{ unit.label }}</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, onMounted, onUnmounted } from 'vue'
import { format } from 'date-fns'
import { Video, Mic, Calendar, DollarSign, Tag, Play, Ticket } from 'lucide-vue-next'
import type { Event } from '@/types'
import { eventImage, onImageError } from '@/utils/eventImage'

const props = defineProps<{
  event: Event
  onBuyTicket: () => void
  hasToken?: boolean
}>()

const formattedDate = computed(() =>
  format(new Date(props.event.scheduled_at), 'EEEE, MMMM d, yyyy · h:mm a')
)

// Inline countdown timer
const now = ref(Date.now())
let timer: ReturnType<typeof setInterval>

onMounted(() => {
  timer = setInterval(() => {
    now.value = Date.now()
  }, 1000)
})

onUnmounted(() => {
  clearInterval(timer)
})

const countdownDiff = computed(() =>
  Math.max(0, new Date(props.event.scheduled_at).getTime() - now.value)
)

const countdownExpired = computed(() => countdownDiff.value === 0)

const countdownUnits = computed(() => {
  const s = Math.floor(countdownDiff.value / 1000)
  return [
    { label: 'Days',    value: String(Math.floor(s / 86400)).padStart(2, '0') },
    { label: 'Hours',   value: String(Math.floor((s % 86400) / 3600)).padStart(2, '0') },
    { label: 'Mins',    value: String(Math.floor((s % 3600) / 60)).padStart(2, '0') },
    { label: 'Secs',    value: String(s % 60).padStart(2, '0') },
  ]
})
</script>
