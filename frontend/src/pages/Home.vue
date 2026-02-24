<template>
  <div class="overflow-x-hidden">

    <!-- ============================================================
         SECTION 1 — HERO
    ============================================================ -->
    <section class="min-h-[80vh] flex items-center relative overflow-hidden">
      <!-- Subtle background glow -->
      <div class="absolute inset-0 pointer-events-none">
        <div
          class="absolute top-1/2 left-1/3 w-[500px] h-[500px] rounded-full opacity-8"
          style="background: radial-gradient(circle, #E8002D 0%, transparent 70%); filter: blur(120px);"
        />
      </div>

      <div class="relative max-w-5xl mx-auto px-4 sm:px-6 lg:px-8 py-24 w-full">
        <div class="animate-fade-in space-y-8 max-w-3xl">

          <!-- Eyebrow -->
          <div class="inline-flex items-center gap-2 bg-accent-red/10 border border-accent-red/20
                      text-accent-red text-xs font-semibold px-3 py-1.5 rounded">
            <span class="w-1.5 h-1.5 rounded-full bg-accent-red animate-pulse" />
            No account required to start
          </div>

          <!-- H1 -->
          <h1 class="font-display text-6xl sm:text-7xl lg:text-8xl text-white leading-none uppercase tracking-wide">
            Go Live.<br />
            <span class="text-accent-red">Get Paid.</span>
          </h1>

          <!-- Subtitle -->
          <p class="text-text-muted text-xl leading-relaxed max-w-xl">
            Test your stream in under a minute — no account, no card.
            Share a link with your audience. Create an account when you're
            ready to charge them.
          </p>

          <!-- Primary CTA -->
          <div class="flex items-center flex-wrap gap-4 pt-2">
            <RouterLink
              to="/try"
              class="inline-flex items-center gap-2 px-8 py-4 rounded-lg bg-accent-red
                     hover:bg-accent-red-hover text-white font-bold text-base
                     transition-all duration-200 active:scale-95 shadow-lg shadow-accent-red/20"
            >
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                  d="M15 10l4.553-2.276A1 1 0 0121 8.677v6.646a1 1 0 01-1.447.894L15 14M5 18h8a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v8a2 2 0 002 2z" />
              </svg>
              Start a Free Test Stream
            </RouterLink>
            <RouterLink
              v-if="liveEvents.length"
              to="/#live"
              class="inline-flex items-center gap-2 px-6 py-4 rounded-lg border border-white/20
                     hover:border-white/40 text-white text-base font-medium
                     transition-all duration-200 hover:bg-white/5"
            >
              Browse Live Events
            </RouterLink>
          </div>

          <!-- Simple trust line — no fake numbers -->
          <p class="text-text-muted text-sm pt-2">
            5 minutes free &nbsp;·&nbsp; M-Pesa payouts &nbsp;·&nbsp; Works with OBS or browser
          </p>
        </div>
      </div>
    </section>

    <!-- ============================================================
         SECTION 2 — LIVE NOW (only shown when real events exist)
    ============================================================ -->
    <section v-if="liveEvents.length" id="live" class="py-16 animate-fade-in">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex items-center justify-between mb-8">
          <div class="flex items-center gap-3">
            <span class="inline-flex items-center gap-1.5 bg-accent-red text-white text-xs font-bold
                         uppercase tracking-wider px-3 py-1.5 rounded">
              <span class="w-1.5 h-1.5 bg-white rounded-full animate-pulse" />
              Live
            </span>
            <h2 class="font-display text-3xl md:text-4xl uppercase tracking-wide text-white">
              Happening Right Now
            </h2>
          </div>
        </div>

        <div class="flex gap-4 overflow-x-auto scrollbar-none pb-2 -mx-4 px-4">
          <div
            v-for="event in liveEvents"
            :key="event.id"
            class="flex-none w-64 sm:w-72"
          >
            <EventCard :event="event" />
          </div>
        </div>
      </div>
    </section>

    <!-- ============================================================
         SECTION 3 — HOW IT WORKS (updated for try-first flow)
    ============================================================ -->
    <section class="py-20 animate-fade-in">
      <div class="max-w-5xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="mb-14">
          <p class="text-accent-red font-semibold text-xs uppercase tracking-widest mb-3">Simple</p>
          <h2 class="font-display text-4xl md:text-5xl uppercase tracking-wide text-white">
            How It Works
          </h2>
        </div>

        <div class="grid sm:grid-cols-3 gap-6">
          <div
            v-for="step in howItWorks"
            :key="step.number"
            class="bg-bg-elevated rounded-xl p-8 border border-white/5
                   hover:border-accent-red/20 transition-all duration-300 group"
          >
            <div class="font-display text-5xl text-accent-red/20 group-hover:text-accent-red/40
                        transition-colors leading-none mb-5 select-none">
              {{ step.number }}
            </div>
            <h3 class="text-white font-bold text-lg mb-2">{{ step.title }}</h3>
            <p class="text-text-muted text-sm leading-relaxed">{{ step.body }}</p>
          </div>
        </div>

        <div class="mt-10">
          <RouterLink
            to="/try"
            class="inline-flex items-center gap-2 px-7 py-3.5 rounded-lg bg-accent-red
                   hover:bg-accent-red-hover text-white font-semibold text-sm
                   transition-all duration-200 active:scale-95"
          >
            Try it now — no account needed
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
            </svg>
          </RouterLink>
        </div>
      </div>
    </section>

    <!-- ============================================================
         SECTION 4 — UPCOMING EVENTS (if any)
    ============================================================ -->
    <section v-if="upcomingEvents.length" class="py-16 bg-bg-surface animate-fade-in">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex items-center justify-between mb-8">
          <h2 class="font-display text-3xl md:text-4xl uppercase tracking-wide text-white">
            Scheduled
          </h2>
        </div>
        <div class="grid sm:grid-cols-2 lg:grid-cols-3 gap-4">
          <EventCard v-for="event in upcomingEvents" :key="event.id" :event="event" />
        </div>
      </div>
    </section>

    <!-- ============================================================
         SECTION 5 — BOTTOM CTA
    ============================================================ -->
    <section class="py-20 animate-fade-in">
      <div class="max-w-5xl mx-auto px-4 sm:px-6 lg:px-8">
        <div
          class="rounded-xl border border-white/10 p-12 md:p-16 bg-bg-elevated"
        >
          <h2 class="font-display text-5xl md:text-6xl uppercase tracking-wide text-white leading-none mb-4">
            Ready to Charge<br />Your Audience?
          </h2>
          <p class="text-text-muted text-base max-w-lg leading-relaxed mb-8">
            Once you've tested the platform and seen it works, create an account.
            Set a ticket price. We handle payment collection via M-Pesa and
            pay you 70% after every session.
          </p>
          <div class="flex items-center flex-wrap gap-4">
            <RouterLink
              to="/register"
              class="inline-flex items-center gap-2 px-7 py-3.5 rounded-lg bg-accent-red
                     hover:bg-accent-red-hover text-white font-semibold text-sm
                     transition-all duration-200 active:scale-95"
            >
              Create Free Account
            </RouterLink>
            <RouterLink
              to="/try"
              class="inline-flex items-center gap-2 px-7 py-3.5 rounded-lg border border-white/20
                     hover:border-white/40 text-white text-sm font-medium
                     transition-all duration-200 hover:bg-white/5"
            >
              Test the stream first
            </RouterLink>
          </div>
        </div>
      </div>
    </section>

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
import { RouterLink, useRouter, useRoute } from 'vue-router'
import { useEventsStore } from '@/stores/events'
import { useAuthStore } from '@/stores/auth'
import EventCard from '@/components/events/EventCard.vue'
import MpesaModal from '@/components/payment/MpesaModal.vue'
import type { Event } from '@/types'

const eventsStore = useEventsStore()
const auth = useAuthStore()
const route = useRoute()
const router = useRouter()

const showPaymentModal = ref(false)
const selectedEvent = ref<Event | null>(null)

function load() {
  const category = route.query.category as string | undefined
  eventsStore.fetchEvents(category ? { sport: category } : undefined)
}

onMounted(load)
watch(() => route.query.category, load)

const liveEvents = computed(() => eventsStore.events.filter((e) => e.status === 'live'))
const upcomingEvents = computed(() => eventsStore.events.filter((e) => e.status === 'scheduled'))

function openPayment(event: Event) {
  if (!auth.isAuthenticated) {
    router.push({ name: 'login', query: { redirect: `/events/${event.id}` } })
    return
  }
  selectedEvent.value = event
  showPaymentModal.value = true
}
void openPayment

const howItWorks = [
  {
    number: '01',
    title: 'Test for free',
    body: 'Start a 5-minute stream right now — no account, no credit card. Get a real RTMP stream key and a shareable viewer link.',
  },
  {
    number: '02',
    title: 'Share the link',
    body: 'Send your viewer URL to your audience. They watch without paying. You see if the platform delivers.',
  },
  {
    number: '03',
    title: 'Go paid when ready',
    body: 'Create an account, set a ticket price, and charge your audience via M-Pesa. 70% goes to you after every session.',
  },
]
</script>
