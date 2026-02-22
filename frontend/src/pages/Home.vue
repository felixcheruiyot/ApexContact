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

    <!-- Platform intro (shown when no events are featured) -->
    <div v-else class="relative overflow-hidden min-h-[60vh] flex items-center">
      <div class="absolute inset-0 bg-gradient-to-br from-bg via-bg to-accent-red/10" />
      <div class="absolute inset-0"
        style="background-image: radial-gradient(circle at 70% 50%, rgba(232,0,45,0.06) 0%, transparent 60%);" />
      <div class="relative max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-24 text-center">
        <p class="text-accent-red font-semibold text-sm uppercase tracking-widest mb-4">Live · Paid · Real-time</p>
        <h1 class="font-display text-6xl md:text-8xl uppercase tracking-wide text-white leading-none mb-6">
          Knowledge<br /><span class="text-accent-red">Goes Live</span>
        </h1>
        <p class="text-text-muted text-lg md:text-xl leading-relaxed mb-10 max-w-2xl mx-auto">
          Host paid or free events online — sports, mentorship sessions, legal consultations,
          music performances, and more. Get paid via M-Pesa, reach anyone, anywhere.
        </p>
        <div class="flex items-center justify-center gap-4 flex-wrap">
          <RouterLink to="/register" class="btn-primary text-base px-8 py-4">Start for Free</RouterLink>
          <RouterLink to="/use-cases" class="btn-ghost text-base px-8 py-4">See Use Cases</RouterLink>
        </div>
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
        <h2 class="section-heading mb-6">Upcoming Events</h2>
        <div class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
          <EventCard v-for="event in upcomingEvents" :key="event.id" :event="event" />
        </div>
      </section>

      <!-- Live Rooms (Commentary) -->
      <section v-if="allLobbies.length">
        <h2 class="section-heading mb-2 flex items-center gap-3">
          <span class="text-accent-orange">🎙</span> Live Rooms
        </h2>
        <p class="text-text-muted text-sm mb-6">Join live audio discussions and expert sessions</p>
        <div class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
          <LobbyCard v-for="lobby in allLobbies" :key="lobby.id" :lobby="lobby" />
        </div>
      </section>

      <!-- Empty state -->
      <div v-if="!eventsStore.loading && !eventsStore.events.length && !allLobbies.length" class="text-center py-16">
        <p class="text-5xl mb-4">🎬</p>
        <h3 class="text-white font-semibold text-xl mb-2">No events scheduled yet</h3>
        <p class="text-text-muted mb-6">Be the first to host an event on Live Streamify.</p>
        <RouterLink to="/register" class="btn-primary text-sm px-6 py-3">Host an Event</RouterLink>
      </div>

      <!-- Use cases teaser -->
      <section class="bg-bg-surface rounded-2xl border border-white/5 p-8 md:p-12">
        <div class="text-center mb-10">
          <p class="text-accent-red font-semibold text-sm uppercase tracking-widest mb-3">Built for Everyone</p>
          <h2 class="font-display text-4xl md:text-5xl uppercase tracking-wide text-white mb-4">
            What Will You Stream?
          </h2>
          <p class="text-text-muted text-base max-w-xl mx-auto">
            From boxing rings to boardrooms — if you have an audience, we have a platform.
          </p>
        </div>
        <div class="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-6 gap-3 mb-8">
          <div v-for="uc in useCaseTeaser" :key="uc.label"
            class="bg-bg-elevated rounded-xl p-4 text-center hover:bg-white/5 transition-colors cursor-default">
            <div class="text-2xl mb-2">{{ uc.icon }}</div>
            <p class="text-white text-xs font-semibold">{{ uc.label }}</p>
          </div>
        </div>
        <div class="text-center">
          <RouterLink to="/use-cases" class="btn-ghost text-sm px-6 py-3">
            Explore All Use Cases →
          </RouterLink>
        </div>
      </section>

      <!-- Testimonials -->
      <section>
        <div class="text-center mb-10">
          <p class="text-accent-orange font-semibold text-sm uppercase tracking-widest mb-3">What Hosts Say</p>
          <h2 class="font-display text-4xl md:text-5xl uppercase tracking-wide text-white">
            Trusted by Experts
          </h2>
        </div>
        <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
          <div v-for="t in testimonials" :key="t.name"
            class="card p-6 flex flex-col gap-4">
            <div class="flex gap-1">
              <span v-for="n in 5" :key="n" class="text-accent-orange text-sm">★</span>
            </div>
            <p class="text-text-muted text-sm leading-relaxed flex-1">"{{ t.quote }}"</p>
            <div class="flex items-center gap-3 border-t border-white/5 pt-4">
              <div class="w-9 h-9 rounded-full bg-accent-red/20 flex items-center justify-center text-accent-red font-bold text-sm">
                {{ t.initials }}
              </div>
              <div>
                <p class="text-white text-sm font-semibold">{{ t.name }}</p>
                <p class="text-text-muted text-xs">{{ t.role }}</p>
              </div>
            </div>
          </div>
        </div>
      </section>
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
import { useRoute, useRouter, RouterLink } from 'vue-router'
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
  const category = route.query.category as string | undefined
  eventsStore.fetchEvents(category ? { sport: category } : undefined)
  commentaryStore.fetchLobbies(category ? { sport: category } : undefined)
}

onMounted(load)
watch(() => route.query.category, load)

const liveEvents = computed(() => eventsStore.events.filter((e) => e.status === 'live'))
const upcomingEvents = computed(() => eventsStore.events.filter((e) => e.status === 'scheduled'))
const liveLobbies = computed(() => commentaryStore.lobbies.filter((l) => l.status === 'live'))
const upcomingLobbies = computed(() => commentaryStore.lobbies.filter((l) => l.status === 'scheduled'))
const allLobbies = computed(() => [...liveLobbies.value, ...upcomingLobbies.value])

const featuredEvent = computed<Event | undefined>(() =>
  liveEvents.value[0] ?? upcomingEvents.value[0]
)

const useCaseTeaser = [
  { icon: '🥊', label: 'Sports' },
  { icon: '🏫', label: 'Mentoring' },
  { icon: '🌍', label: 'Visa Guides' },
  { icon: '🎵', label: 'Music' },
  { icon: '⚖️', label: 'Legal' },
  { icon: '📚', label: 'Education' },
]

const testimonials = [
  {
    name: 'Sarah Mutua',
    role: 'Boxing Promoter, Nairobi',
    initials: 'SM',
    quote: 'We sold 2,400 tickets for our championship fight. Revenue hit our account within 48 hours. Nothing else comes close for African promoters.',
  },
  {
    name: 'Dr. James Odhiambo',
    role: 'Career Mentor & Coach',
    initials: 'JO',
    quote: 'My mentorship sessions now reach students across East Africa. The M-Pesa integration means no one is left out — it just works.',
  },
  {
    name: 'Aisha Kamau',
    role: 'Immigration Consultant',
    initials: 'AK',
    quote: 'I turned my visa application workshops into a recurring income stream. Clients attend live, ask questions, and pay with M-Pesa in seconds.',
  },
]

function openPayment(event: Event) {
  if (!auth.isAuthenticated) {
    router.push({ name: 'login', query: { redirect: `/events/${event.id}` } })
    return
  }
  selectedEvent.value = event
  showPaymentModal.value = true
}
</script>
