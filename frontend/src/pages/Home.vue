<template>
  <div class="overflow-x-hidden">

    <!-- ============================================================
         SECTION 1 — HERO
    ============================================================ -->
    <section class="min-h-screen flex items-center relative overflow-hidden">
      <!-- Background radial glow -->
      <div class="absolute inset-0 pointer-events-none">
        <div
          class="absolute top-1/3 left-1/4 w-[600px] h-[600px] rounded-full opacity-10"
          style="background: radial-gradient(circle, #E8002D 0%, transparent 70%); filter: blur(80px);"
        />
        <div
          class="absolute bottom-0 right-0 w-[400px] h-[400px] rounded-full opacity-8"
          style="background: radial-gradient(circle, #FF6B00 0%, transparent 70%); filter: blur(100px);"
        />
      </div>

      <div class="relative max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-20 lg:py-0 w-full">
        <div class="grid lg:grid-cols-[55%_45%] gap-12 lg:gap-16 items-center min-h-screen lg:min-h-0 lg:py-24">

          <!-- Left: copy -->
          <div class="animate-fade-in space-y-8">
            <!-- Live pill badge -->
            <div class="inline-flex items-center gap-2.5 bg-accent-red/15 border border-accent-red/25
                        text-white text-xs font-semibold px-4 py-2 rounded-full">
              <span class="relative flex h-2 w-2">
                <span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-accent-red opacity-75" />
                <span class="relative inline-flex rounded-full h-2 w-2 bg-accent-red" />
              </span>
              LIVE NOW &nbsp;·&nbsp; 47 Creators Streaming
            </div>

            <!-- H1 -->
            <h1 class="font-display text-6xl sm:text-7xl lg:text-8xl text-white leading-none uppercase tracking-wide">
              Where Africa's<br />
              <span class="text-accent-red">Experts</span><br />
              Go Live
            </h1>

            <!-- Subtitle -->
            <p class="text-text-muted text-lg leading-relaxed max-w-lg">
              Join thousands learning sales, law, business, fitness and more —
              live and paid direct to the creator.
            </p>

            <!-- CTA row -->
            <div class="flex items-center flex-wrap gap-4">
              <RouterLink
                to="/"
                class="inline-flex items-center px-7 py-3.5 rounded-full bg-accent-red
                       hover:bg-accent-red-hover text-white font-semibold text-sm
                       transition-all duration-200 active:scale-95 shadow-lg shadow-accent-red/25"
              >
                Browse Live Events
              </RouterLink>
              <RouterLink
                to="/promoters"
                class="inline-flex items-center px-7 py-3.5 rounded-full border border-white/25
                       hover:border-white/50 text-white font-semibold text-sm
                       transition-all duration-200 hover:bg-white/5"
              >
                Start Hosting
              </RouterLink>
            </div>

            <!-- Stats row -->
            <div class="flex items-center flex-wrap gap-x-6 gap-y-2 text-text-muted text-sm pt-2">
              <span class="flex items-center gap-1.5">
                <span class="font-semibold text-white">12,400+</span> Creators
              </span>
              <span class="text-white/20">|</span>
              <span class="flex items-center gap-1.5">
                <span class="font-semibold text-white">58,000+</span> Viewers
              </span>
              <span class="text-white/20">|</span>
              <span class="flex items-center gap-1.5">
                <span class="font-semibold text-white">KES 4.2M+</span> Paid Out
              </span>
            </div>
          </div>

          <!-- Right: 2x2 mosaic creator grid -->
          <div class="hidden lg:grid grid-cols-2 gap-3 animate-fade-in">
            <!-- Top-left tall -->
            <div class="relative rounded-2xl overflow-hidden h-72 group">
              <img
                src="https://images.pexels.com/photos/3769021/pexels-photo-3769021.jpeg?auto=compress&cs=tinysrgb&w=600&q=80"
                alt="Creator Amara K."
                class="w-full h-full object-cover transition-transform duration-500 group-hover:scale-105"
              />
              <div class="absolute inset-0 bg-gradient-to-t from-black/70 via-transparent to-transparent" />
              <div class="absolute bottom-3 left-3">
                <p class="text-white text-sm font-semibold leading-tight">Amara K.</p>
                <p class="text-text-muted text-xs">Business</p>
              </div>
            </div>

            <!-- Top-right short -->
            <div class="relative rounded-2xl overflow-hidden h-48 group self-end">
              <img
                src="https://images.pexels.com/photos/1139743/pexels-photo-1139743.jpeg?auto=compress&cs=tinysrgb&w=600&q=80"
                alt="Creator James O."
                class="w-full h-full object-cover transition-transform duration-500 group-hover:scale-105"
              />
              <div class="absolute inset-0 bg-gradient-to-t from-black/70 via-transparent to-transparent" />
              <div class="absolute bottom-3 left-3">
                <p class="text-white text-sm font-semibold leading-tight">James O.</p>
                <p class="text-text-muted text-xs">Sales</p>
              </div>
            </div>

            <!-- Bottom-left short -->
            <div class="relative rounded-2xl overflow-hidden h-48 group self-start">
              <img
                src="https://images.pexels.com/photos/3756879/pexels-photo-3756879.jpeg?auto=compress&cs=tinysrgb&w=600&q=80"
                alt="Creator Fatima A."
                class="w-full h-full object-cover transition-transform duration-500 group-hover:scale-105"
              />
              <div class="absolute inset-0 bg-gradient-to-t from-black/70 via-transparent to-transparent" />
              <div class="absolute bottom-3 left-3">
                <p class="text-white text-sm font-semibold leading-tight">Fatima A.</p>
                <p class="text-text-muted text-xs">Legal</p>
              </div>
            </div>

            <!-- Bottom-right tall -->
            <div class="relative rounded-2xl overflow-hidden h-72 group">
              <img
                src="https://images.pexels.com/photos/5717546/pexels-photo-5717546.jpeg?auto=compress&cs=tinysrgb&w=600&q=80"
                alt="Creator Zara M."
                class="w-full h-full object-cover transition-transform duration-500 group-hover:scale-105"
              />
              <div class="absolute inset-0 bg-gradient-to-t from-black/70 via-transparent to-transparent" />
              <div class="absolute bottom-3 left-3">
                <p class="text-white text-sm font-semibold leading-tight">Zara M.</p>
                <p class="text-text-muted text-xs">Lifestyle</p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- ============================================================
         SECTION 2 — LIVE NOW (conditional)
    ============================================================ -->
    <section v-if="liveEvents.length" class="py-16 animate-fade-in">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <!-- Header row -->
        <div class="flex items-center justify-between mb-8">
          <div class="flex items-center gap-3">
            <span class="inline-flex items-center gap-1.5 bg-accent-red text-white text-xs font-bold
                         uppercase tracking-wider px-3 py-1.5 rounded-full">
              <span class="w-1.5 h-1.5 bg-white rounded-full animate-pulse" />
              Live
            </span>
            <h2 class="font-display text-3xl md:text-4xl uppercase tracking-wide text-white">
              Happening Right Now
            </h2>
          </div>
          <RouterLink
            to="/"
            class="text-text-muted hover:text-white text-sm font-medium transition-colors hidden sm:block"
          >
            View all &rarr;
          </RouterLink>
        </div>

        <!-- Scrollable event row -->
        <div class="flex gap-4 overflow-x-auto scrollbar-none pb-2 -mx-4 px-4">
          <div
            v-for="event in liveEvents"
            :key="event.id"
            class="flex-none w-64 sm:w-72"
          >
            <EventCard :event="event" />
          </div>
        </div>

        <RouterLink
          to="/"
          class="sm:hidden mt-4 inline-flex text-text-muted hover:text-white text-sm
                 font-medium transition-colors"
        >
          View all &rarr;
        </RouterLink>
      </div>
    </section>

    <!-- ============================================================
         SECTION 3 — BROWSE CATEGORIES
    ============================================================ -->
    <section class="py-16 animate-fade-in">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="text-center mb-12">
          <p class="text-accent-red font-semibold text-xs uppercase tracking-widest mb-3">
            Topics
          </p>
          <h2 class="font-display text-4xl md:text-5xl uppercase tracking-wide text-white">
            Find Your Topic
          </h2>
        </div>

        <div class="grid grid-cols-2 sm:grid-cols-4 gap-4">
          <div
            v-for="cat in categories"
            :key="cat.name"
            class="group bg-bg-elevated border border-transparent hover:border-accent-red/50
                   rounded-2xl p-6 cursor-pointer transition-all duration-200 hover:bg-bg-elevated/80
                   hover:-translate-y-1 text-center"
          >
            <component :is="cat.icon" class="w-8 h-8 mx-auto mb-3 text-text-muted group-hover:text-accent-red transition-colors duration-200" />
            <p class="text-white font-semibold text-sm mb-1">{{ cat.name }}</p>
            <p class="text-text-muted text-xs">{{ cat.sessions }}+ sessions</p>
          </div>
        </div>
      </div>
    </section>

    <!-- ============================================================
         SECTION 4 — CREATOR SPOTLIGHT
    ============================================================ -->
    <section class="py-16 bg-bg-surface animate-fade-in">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="grid lg:grid-cols-2 gap-12 lg:gap-20 items-center">

          <!-- Image -->
          <div class="relative">
            <div class="aspect-square rounded-3xl overflow-hidden">
              <img
                src="https://images.pexels.com/photos/2379004/pexels-photo-2379004.jpeg?auto=compress&cs=tinysrgb&w=1260&q=80"
                alt="David Mwangi — Business Strategy Coach"
                class="w-full h-full object-cover"
              />
            </div>
            <!-- Floating stat card -->
            <div class="absolute -bottom-4 -right-4 sm:bottom-6 sm:right-6
                        bg-bg-elevated border border-white/10 rounded-2xl px-5 py-4 shadow-2xl">
              <p class="text-accent-red font-display text-3xl">KES 340K</p>
              <p class="text-text-muted text-xs mt-0.5">Earned in 3 months</p>
            </div>
          </div>

          <!-- Copy -->
          <div class="space-y-6 lg:space-y-8">
            <p class="text-accent-orange text-xs font-bold uppercase tracking-widest">
              Creator Spotlight
            </p>
            <div>
              <h3 class="font-display text-5xl md:text-6xl uppercase tracking-wide text-white leading-none">
                David Mwangi
              </h3>
              <p class="text-text-muted mt-2 text-base">
                Business Strategy Coach — Nairobi
              </p>
            </div>
            <blockquote class="text-white text-lg leading-relaxed border-l-2 border-accent-red pl-5">
              "I earned KES 340,000 in my first 3 months hosting live workshops on
              negotiation and sales. This platform gave me the tools and the audience."
            </blockquote>

            <!-- Stats row -->
            <div class="flex items-center gap-8 pt-2">
              <div>
                <p class="text-white font-bold text-2xl">2,400</p>
                <p class="text-text-muted text-xs mt-0.5">Total Viewers</p>
              </div>
              <div class="w-px h-10 bg-white/10" />
              <div>
                <p class="text-white font-bold text-2xl">18</p>
                <p class="text-text-muted text-xs mt-0.5">Sessions Hosted</p>
              </div>
              <div class="w-px h-10 bg-white/10" />
              <div>
                <p class="text-white font-bold text-2xl">KES 340K</p>
                <p class="text-text-muted text-xs mt-0.5">Total Earned</p>
              </div>
            </div>

            <RouterLink
              to="/promoters"
              class="inline-flex items-center px-7 py-3.5 rounded-full bg-accent-red
                     hover:bg-accent-red-hover text-white font-semibold text-sm
                     transition-all duration-200 active:scale-95"
            >
              Become a Creator &rarr;
            </RouterLink>
          </div>
        </div>
      </div>
    </section>

    <!-- ============================================================
         SECTION 5 — HOW IT WORKS
    ============================================================ -->
    <section class="py-20 animate-fade-in">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="text-center mb-14">
          <p class="text-accent-red font-semibold text-xs uppercase tracking-widest mb-3">Simple</p>
          <h2 class="font-display text-4xl md:text-5xl uppercase tracking-wide text-white">
            Go Live in 3 Steps
          </h2>
        </div>

        <div class="grid sm:grid-cols-3 gap-6">
          <div
            v-for="step in howItWorks"
            :key="step.number"
            class="bg-bg-elevated rounded-2xl p-8 border border-white/5
                   hover:border-accent-red/20 transition-all duration-300 hover:-translate-y-1 group"
          >
            <!-- Step number -->
            <div class="font-display text-6xl text-accent-red/20 group-hover:text-accent-red/40
                        transition-colors leading-none mb-4 select-none">
              {{ step.number }}
            </div>
            <!-- Icon -->
            <div class="mb-4">
              <component :is="step.icon" class="text-text-muted" />
            </div>
            <h3 class="text-white font-bold text-lg mb-2">{{ step.title }}</h3>
            <p class="text-text-muted text-sm leading-relaxed">{{ step.body }}</p>
          </div>
        </div>
      </div>
    </section>

    <!-- ============================================================
         SECTION 6 — TESTIMONIALS
    ============================================================ -->
    <section class="py-20 bg-bg-surface animate-fade-in">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="text-center mb-14">
          <p class="text-accent-orange font-semibold text-xs uppercase tracking-widest mb-3">
            Social Proof
          </p>
          <h2 class="font-display text-4xl md:text-5xl uppercase tracking-wide text-white">
            Loved by Creators Across Africa
          </h2>
        </div>

        <div class="grid sm:grid-cols-2 lg:grid-cols-3 gap-6">
          <div
            v-for="t in testimonials"
            :key="t.name"
            class="bg-bg-elevated rounded-2xl p-8 border border-white/5
                   hover:border-white/10 transition-all duration-200 flex flex-col gap-5"
          >
            <!-- Stars -->
            <div class="flex gap-1">
              <svg
                v-for="n in 5"
                :key="n"
                class="w-4 h-4 text-amber-400"
                fill="currentColor"
                viewBox="0 0 20 20"
              >
                <path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z"/>
              </svg>
            </div>

            <!-- Quote -->
            <p class="text-text-muted text-sm leading-relaxed flex-1">
              "{{ t.quote }}"
            </p>

            <!-- Author -->
            <div class="flex items-center gap-3 pt-4 border-t border-white/5">
              <img
                :src="t.avatar"
                :alt="t.name"
                class="w-10 h-10 rounded-full object-cover ring-2 ring-white/10"
              />
              <div>
                <p class="text-white text-sm font-bold leading-tight">{{ t.name }}</p>
                <p class="text-text-muted text-xs">{{ t.role }}</p>
                <p class="text-text-muted text-xs">{{ t.city }}</p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- ============================================================
         SECTION 7 — BOTTOM CTA
    ============================================================ -->
    <section class="py-20 animate-fade-in">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div
          class="rounded-3xl border border-accent-red/20 p-12 md:p-20 text-center
                 bg-gradient-to-br from-accent-red/15 via-bg-elevated/50 to-transparent
                 relative overflow-hidden"
        >
          <!-- Glow -->
          <div
            class="absolute top-0 left-1/2 -translate-x-1/2 w-[500px] h-[300px] pointer-events-none"
            style="background: radial-gradient(ellipse, rgba(232,0,45,0.15) 0%, transparent 70%); filter: blur(40px);"
          />

          <div class="relative space-y-6">
            <h2 class="font-display text-5xl md:text-7xl uppercase tracking-wide text-white leading-none">
              Ready to Go Live?
            </h2>
            <p class="text-text-muted text-lg max-w-xl mx-auto leading-relaxed">
              Join 12,000+ creators earning from their knowledge on Africa's #1 live platform.
            </p>
            <div class="flex items-center justify-center flex-wrap gap-4 pt-2">
              <RouterLink
                to="/register"
                class="inline-flex items-center px-8 py-4 rounded-full bg-white text-black
                       font-bold text-sm hover:bg-gray-100 transition-all duration-200
                       active:scale-95 shadow-xl"
              >
                Create Free Account
              </RouterLink>
              <RouterLink
                to="/"
                class="inline-flex items-center px-8 py-4 rounded-full border border-white/25
                       hover:border-white/50 text-white font-semibold text-sm
                       transition-all duration-200 hover:bg-white/5"
              >
                Browse Events
              </RouterLink>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- ============================================================
         PAYMENT MODAL
    ============================================================ -->
    <MpesaModal
      v-if="selectedEvent && showPaymentModal"
      :event="selectedEvent"
      @close="showPaymentModal = false"
    />

  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch, h } from 'vue'
import { RouterLink, useRouter, useRoute } from 'vue-router'
import { TrendingUp, Scale, Dumbbell, Music, BookOpen, ChefHat, Briefcase, Globe } from 'lucide-vue-next'
import { useEventsStore } from '@/stores/events'
import { useAuthStore } from '@/stores/auth'
import { useCommentaryStore } from '@/stores/commentary'
import EventCard from '@/components/events/EventCard.vue'
import MpesaModal from '@/components/payment/MpesaModal.vue'
import type { Event } from '@/types'

// Stores
const eventsStore = useEventsStore()
const commentaryStore = useCommentaryStore()
const auth = useAuthStore()
const route = useRoute()
const router = useRouter()

// Payment modal state
const showPaymentModal = ref(false)
const selectedEvent = ref<Event | null>(null)

// Data loading
function load() {
  const category = route.query.category as string | undefined
  eventsStore.fetchEvents(category ? { sport: category } : undefined)
  commentaryStore.fetchLobbies(category ? { sport: category } : undefined)
}

onMounted(load)
watch(() => route.query.category, load)

// Derived event lists
const liveEvents = computed(() => eventsStore.events.filter((e) => e.status === 'live'))

function openPayment(event: Event) {
  if (!auth.isAuthenticated) {
    router.push({ name: 'login', query: { redirect: `/events/${event.id}` } })
    return
  }
  selectedEvent.value = event
  showPaymentModal.value = true
}
// kept for future use — suppress unused warning
void openPayment

// ---- STATIC DATA ----

const categories = [
  { icon: TrendingUp, name: 'Business & Finance', sessions: 340 },
  { icon: Scale, name: 'Legal & Visa', sessions: 128 },
  { icon: Dumbbell, name: 'Fitness & Wellness', sessions: 215 },
  { icon: Music, name: 'Music & Arts', sessions: 97 },
  { icon: BookOpen, name: 'Education', sessions: 402 },
  { icon: ChefHat, name: 'Cooking', sessions: 81 },
  { icon: Briefcase, name: 'Sales', sessions: 189 },
  { icon: Globe, name: 'Community', sessions: 143 },
]

// Inline SVG icon components for How It Works
const PencilIcon = {
  render() {
    return h('svg', { class: 'w-7 h-7', fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24' }, [
      h('path', {
        'stroke-linecap': 'round',
        'stroke-linejoin': 'round',
        'stroke-width': '2',
        d: 'M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z',
      }),
    ])
  },
}

const VideoIcon = {
  render() {
    return h('svg', { class: 'w-7 h-7', fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24' }, [
      h('path', {
        'stroke-linecap': 'round',
        'stroke-linejoin': 'round',
        'stroke-width': '2',
        d: 'M15 10l4.553-2.276A1 1 0 0121 8.677v6.646a1 1 0 01-1.447.894L15 14M5 18h8a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v8a2 2 0 002 2z',
      }),
    ])
  },
}

const BanknoteIcon = {
  render() {
    return h('svg', { class: 'w-7 h-7', fill: 'none', stroke: 'currentColor', viewBox: '0 0 24 24' }, [
      h('path', {
        'stroke-linecap': 'round',
        'stroke-linejoin': 'round',
        'stroke-width': '2',
        d: 'M17 9V7a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2m2 4h10a2 2 0 002-2v-6a2 2 0 00-2-2H9a2 2 0 00-2 2v6a2 2 0 002 2zm7-5a2 2 0 11-4 0 2 2 0 014 0z',
      }),
    ])
  },
}

const howItWorks = [
  {
    number: '01',
    icon: PencilIcon,
    title: 'Create',
    body: 'Set your price, pick a date, write your topic. Your event page is live in minutes.',
  },
  {
    number: '02',
    icon: VideoIcon,
    title: 'Go Live',
    body: 'Stream directly from your browser or OBS. We handle the infrastructure.',
  },
  {
    number: '03',
    icon: BanknoteIcon,
    title: 'Get Paid',
    body: '70% goes to you via M-Pesa, instantly after every session. No delays.',
  },
]

const testimonials = [
  {
    name: 'Fatima A.',
    role: 'Immigration Consultant',
    city: 'Nairobi',
    avatar:
      'https://images.pexels.com/photos/1181690/pexels-photo-1181690.jpeg?auto=compress&cs=tinysrgb&w=400&q=80',
    quote:
      'This platform changed how I monetize my expertise. I ran a live visa consultation and made KES 85,000 in one evening.',
  },
  {
    name: 'James O.',
    role: 'Sales Trainer',
    city: 'Kampala',
    avatar:
      'https://images.pexels.com/photos/1139743/pexels-photo-1139743.jpeg?auto=compress&cs=tinysrgb&w=400&q=80',
    quote:
      'As a sales trainer, I needed a platform that Kenyans trust. Live Streamify pays via M-Pesa instantly. Game changer.',
  },
  {
    name: 'Zara M.',
    role: 'Chef & Lifestyle Creator',
    city: 'Lagos',
    avatar:
      'https://images.pexels.com/photos/3756879/pexels-photo-3756879.jpeg?auto=compress&cs=tinysrgb&w=400&q=80',
    quote:
      'I ran a live cooking masterclass for 800 people. The quality was flawless. My followers finally took me seriously.',
  },
]
</script>
