<template>
  <div>
    <!-- Hero: featured live or next scheduled event -->
    <EventHero v-if="featuredEvent" :event="featuredEvent" @buy-ticket="openPayment(featuredEvent!)" />

    <!-- Skeleton hero if loading -->
    <div v-else-if="eventsStore.loading"
      class="min-h-[60vh] bg-bg-surface animate-pulse flex items-end p-12">
      <div class="max-w-lg space-y-4">
        <div class="h-4 bg-bg-elevated rounded w-24" />
        <div class="h-16 bg-bg-elevated rounded w-96" />
        <div class="h-4 bg-bg-elevated rounded w-64" />
      </div>
    </div>

    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12 space-y-14">
      <!-- Live Now -->
      <section v-if="liveEvents.length">
        <h2 class="section-heading mb-6 flex items-center gap-3">
          Live Now
          <span class="badge-live text-xs">
            <span class="w-1.5 h-1.5 bg-white rounded-full animate-pulse" />
            {{ liveEvents.length }}
          </span>
        </h2>
        <div class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
          <EventCard v-for="event in liveEvents" :key="event.id" :event="event" />
        </div>
      </section>

      <!-- Upcoming -->
      <section v-if="upcomingEvents.length">
        <h2 class="section-heading mb-6">
          {{ currentSport ? `Upcoming ${currentSport.charAt(0).toUpperCase() + currentSport.slice(1)}` : 'Upcoming Events' }}
        </h2>
        <div class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
          <EventCard v-for="event in upcomingEvents" :key="event.id" :event="event" />
        </div>
      </section>

      <!-- Commentary Lobbies -->
      <section v-if="allLobbies.length">
        <h2 class="section-heading mb-2 flex items-center gap-3">
          <span class="text-accent-orange">🎙</span> Commentary
        </h2>
        <p class="text-text-muted text-sm mb-6">Join live audio discussions from fans worldwide</p>
        <div class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
          <LobbyCard v-for="lobby in allLobbies" :key="lobby.id" :lobby="lobby" />
        </div>
      </section>

      <!-- Empty state -->
      <div v-if="!eventsStore.loading && !eventsStore.events.length && !allLobbies.length" class="text-center py-24">
        <p class="text-5xl mb-4">🎬</p>
        <h3 class="text-white font-semibold text-xl mb-2">No events scheduled yet</h3>
        <p class="text-text-muted">Check back soon for upcoming boxing and racing events.</p>
      </div>
    </div>

    <!-- Payment modal -->
    <MpesaModal
      v-if="selectedEvent && showPaymentModal"
      :event="selectedEvent"
      @close="showPaymentModal = false"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useEventsStore } from '@/stores/events'
import { useAuthStore } from '@/stores/auth'
import { useCommentaryStore } from '@/stores/commentary'
import EventHero from '@/components/events/EventHero.vue'
import EventCard from '@/components/events/EventCard.vue'
import LobbyCard from '@/components/commentary/LobbyCard.vue'
import MpesaModal from '@/components/payment/MpesaModal.vue'
import type { Event } from '@/types'

const eventsStore = useEventsStore()
const commentaryStore = useCommentaryStore()
const auth = useAuthStore()
const route = useRoute()
const router = useRouter()

const showPaymentModal = ref(false)
const selectedEvent = ref<Event | null>(null)

function load() {
  const sport = route.query.sport as string | undefined
  eventsStore.fetchEvents(sport ? { sport } : undefined)
  commentaryStore.fetchLobbies(sport ? { sport } : undefined)
}

onMounted(load)

// Re-fetch when the ?sport= query param changes (e.g. clicking Boxing / Racing nav)
watch(() => route.query.sport, load)

const liveEvents = computed(() => eventsStore.events.filter((e) => e.status === 'live'))
const upcomingEvents = computed(() => eventsStore.events.filter((e) => e.status === 'scheduled'))
const liveLobbies = computed(() => commentaryStore.lobbies.filter((l) => l.status === 'live'))
const upcomingLobbies = computed(() => commentaryStore.lobbies.filter((l) => l.status === 'scheduled'))
const allLobbies = computed(() => [...liveLobbies.value, ...upcomingLobbies.value])

const featuredEvent = computed<Event | undefined>(() =>
  liveEvents.value[0] ?? upcomingEvents.value[0]
)

const currentSport = computed(() => route.query.sport as string | undefined)

function openPayment(event: Event) {
  if (!auth.isAuthenticated) {
    router.push({ name: 'login', query: { redirect: `/events/${event.id}` } })
    return
  }
  selectedEvent.value = event
  showPaymentModal.value = true
}
</script>
