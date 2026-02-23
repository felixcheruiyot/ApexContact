<template>
  <div>
    <!-- Hero -->
    <section class="relative overflow-hidden min-h-[70vh] flex items-center">
      <!-- Background -->
      <div class="absolute inset-0 bg-gradient-to-br from-bg via-bg to-accent-red/10" />
      <div class="absolute inset-0"
        style="background-image: radial-gradient(circle at 70% 50%, rgba(232,0,45,0.08) 0%, transparent 60%);" />

      <div class="relative max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-24 grid grid-cols-1 lg:grid-cols-2 gap-16 items-center">
        <div class="animate-fade-in">
          <p class="text-accent-red font-semibold text-sm uppercase tracking-widest mb-4">For Event Hosts</p>
          <h1 class="font-display text-6xl md:text-7xl uppercase tracking-wide text-white leading-none mb-6">
            Bring Your<br />
            <span class="text-accent-red">Events</span><br />
            Worldwide
          </h1>
          <p class="text-text-muted text-lg leading-relaxed mb-8 max-w-lg">
            Live Streamify gives experts, coaches, promoters, and creators a complete platform
            to schedule, sell tickets, and stream live events — with zero technical setup and
            real-time revenue in your account via M-Pesa.
          </p>
          <div class="flex items-center gap-4">
            <RouterLink v-if="auth.isAuthenticated" to="/dashboard/create" class="btn-primary text-base px-8 py-4">
              Create an Event
            </RouterLink>
            <RouterLink v-else to="/register" class="btn-primary text-base px-8 py-4">
              Get Started Free
            </RouterLink>
            <a href="mailto:promoters@livestreamify.com" class="btn-ghost text-base px-8 py-4">
              Talk to Us
            </a>
          </div>
        </div>

        <!-- Stats -->
        <div class="grid grid-cols-2 gap-4">
          <div v-for="stat in stats" :key="stat.label"
            class="card p-6 flex flex-col justify-between">
            <p class="font-display text-5xl text-accent-red mb-2">{{ stat.value }}</p>
            <p class="text-white font-semibold text-sm">{{ stat.label }}</p>
            <p class="text-text-muted text-xs mt-1">{{ stat.sub }}</p>
          </div>
        </div>
      </div>
    </section>

    <!-- How it works -->
    <section class="bg-bg-surface border-y border-white/5 py-20">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="text-center mb-14">
          <p class="text-accent-red text-sm font-semibold uppercase tracking-widest mb-3">Simple Process</p>
          <h2 class="font-display text-4xl md:text-5xl uppercase tracking-wide text-white">
            Live in 4 Steps
          </h2>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-4 gap-6">
          <div v-for="(step, i) in steps" :key="i" class="relative">
            <!-- Connector line -->
            <div v-if="i < steps.length - 1"
              class="hidden md:block absolute top-8 left-1/2 w-full h-px bg-gradient-to-r from-accent-red/50 to-transparent z-0" />
            <div class="card p-6 relative z-10 text-center">
              <div class="w-14 h-14 rounded-full bg-accent-red/10 border border-accent-red/30
                          flex items-center justify-center mx-auto mb-4">
                <span class="font-display text-2xl text-accent-red">{{ i + 1 }}</span>
              </div>
              <h3 class="text-white font-bold text-sm mb-2">{{ step.title }}</h3>
              <p class="text-text-muted text-xs leading-relaxed">{{ step.desc }}</p>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- Revenue & Pricing -->
    <section class="py-20">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="grid grid-cols-1 lg:grid-cols-2 gap-16 items-center">
          <div>
            <p class="text-accent-red text-sm font-semibold uppercase tracking-widest mb-3">Revenue Split</p>
            <h2 class="font-display text-4xl md:text-5xl uppercase tracking-wide text-white mb-6">
              You Keep<br /><span class="text-accent-red">70%</span> of Every Ticket
            </h2>
            <p class="text-text-muted leading-relaxed mb-6">
              No hidden fees, no monthly subscriptions. Live Streamify takes a
              <strong class="text-white">30% platform fee</strong> per ticket sold to cover streaming
              infrastructure, payment processing (M-Pesa), content protection, and customer support.
            </p>
            <p class="text-text-muted leading-relaxed mb-8">
              Revenue is tracked in real time and visible in your host dashboard. Payouts
              are processed within 48 hours of event completion.
            </p>

            <div class="space-y-3">
              <div v-for="item in revenueItems" :key="item.label"
                class="flex items-center justify-between py-3 border-b border-white/5">
                <span class="text-text-muted text-sm">{{ item.label }}</span>
                <span :class="item.highlight ? 'text-accent-red font-bold' : 'text-white font-semibold'" class="text-sm">
                  {{ item.value }}
                </span>
              </div>
            </div>
          </div>

          <!-- Earnings calculator -->
          <div class="card p-8">
            <h3 class="text-white font-bold text-lg mb-6">Earnings Calculator</h3>
            <div class="space-y-5">
              <div>
                <label class="block text-text-muted text-xs uppercase tracking-wider mb-2">
                  Ticket Price (KES)
                </label>
                <input v-model.number="calc.price" type="range" min="100" max="5000" step="50"
                  class="w-full accent-red-500" />
                <div class="flex justify-between text-text-muted text-xs mt-1">
                  <span>KES 100</span>
                  <span class="text-white font-semibold">KES {{ calc.price.toLocaleString() }}</span>
                  <span>KES 5,000</span>
                </div>
              </div>
              <div>
                <label class="block text-text-muted text-xs uppercase tracking-wider mb-2">
                  Expected Viewers
                </label>
                <input v-model.number="calc.viewers" type="range" min="10" max="10000" step="10"
                  class="w-full accent-red-500" />
                <div class="flex justify-between text-text-muted text-xs mt-1">
                  <span>10</span>
                  <span class="text-white font-semibold">{{ calc.viewers.toLocaleString() }} viewers</span>
                  <span>10,000</span>
                </div>
              </div>
            </div>

            <div class="mt-8 pt-6 border-t border-white/10 space-y-3">
              <div class="flex items-center justify-between">
                <span class="text-text-muted text-sm">Gross revenue</span>
                <span class="text-white font-semibold">KES {{ grossRevenue.toLocaleString() }}</span>
              </div>
              <div class="flex items-center justify-between">
                <span class="text-text-muted text-sm">Platform fee (30%)</span>
                <span class="text-text-muted">− KES {{ platformFee.toLocaleString() }}</span>
              </div>
              <div class="flex items-center justify-between pt-3 border-t border-white/10">
                <span class="text-white font-bold">Your earnings</span>
                <span class="text-accent-red font-bold text-2xl">KES {{ yourEarnings.toLocaleString() }}</span>
              </div>
            </div>

            <RouterLink v-if="auth.isAuthenticated" to="/dashboard/create" class="btn-primary w-full mt-6 text-center block">
              Start Earning →
            </RouterLink>
            <RouterLink v-else to="/register" class="btn-primary w-full mt-6 text-center block">
              Start Earning →
            </RouterLink>
          </div>
        </div>
      </div>
    </section>

    <!-- Features -->
    <section class="bg-bg-surface border-y border-white/5 py-20">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="text-center mb-14">
          <p class="text-accent-red text-sm font-semibold uppercase tracking-widest mb-3">Everything Included</p>
          <h2 class="font-display text-4xl md:text-5xl uppercase tracking-wide text-white">
            Your Complete Toolkit
          </h2>
          <p class="text-text-muted text-base mt-3 max-w-xl mx-auto">
            Everything you need to host sports events, mentorship sessions, workshops, concerts, and more.
          </p>
        </div>

        <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
          <div v-for="feature in features" :key="feature.title"
            class="card p-6 group hover:border hover:border-accent-red/30 transition-all">
            <component :is="feature.icon" class="w-8 h-8 mb-4 text-accent-red" />
            <h3 class="text-white font-bold mb-2">{{ feature.title }}</h3>
            <p class="text-text-muted text-sm leading-relaxed">{{ feature.desc }}</p>
          </div>
        </div>
      </div>
    </section>

    <!-- Anti-piracy trust section -->
    <section class="py-20">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="card p-10 md:p-14 text-center max-w-3xl mx-auto">
          <Shield class="w-16 h-16 mx-auto mb-6 text-accent-red" />
          <h2 class="font-display text-4xl uppercase tracking-wide text-white mb-4">
            Your Content is Protected
          </h2>
          <p class="text-text-muted leading-relaxed mb-8">
            Every ticket generates a unique, single-use stream token. Tokens are bound to the
            viewer's device fingerprint and IP address. Sharing links is technically impossible —
            the moment a second device opens the stream, both sessions are terminated and the
            incident is flagged for review.
          </p>
          <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
            <div v-for="shield in shields" :key="shield" class="bg-bg-elevated rounded-lg px-4 py-3">
              <p class="text-white text-xs font-semibold">{{ shield }}</p>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- CTA -->
    <section class="bg-accent-red py-20">
      <div class="max-w-3xl mx-auto px-4 text-center">
        <h2 class="font-display text-5xl uppercase tracking-wide text-white mb-4">
          Ready to Go Live?
        </h2>
        <p class="text-white/80 text-lg mb-8">
          Create your host account in minutes. No upfront costs — you only pay when tickets sell.
        </p>
        <div class="flex items-center justify-center gap-4 flex-wrap">
          <RouterLink v-if="auth.isAuthenticated" to="/dashboard/create"
            class="bg-white text-accent-red font-bold px-8 py-4 rounded-lg hover:bg-white/90 transition-colors text-base">
            Create an Event
          </RouterLink>
          <RouterLink v-else to="/register"
            class="bg-white text-accent-red font-bold px-8 py-4 rounded-lg hover:bg-white/90 transition-colors text-base">
            Create Host Account
          </RouterLink>
          <a href="mailto:hosts@livestreamify.com"
            class="border-2 border-white/50 text-white font-bold px-8 py-4 rounded-lg
                   hover:border-white transition-colors text-base">
            Contact Sales
          </a>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { RouterLink } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { Calendar, Tv, Smartphone, BarChart2, DollarSign, Shield } from 'lucide-vue-next'

const auth = useAuthStore()

const stats = [
  { value: '70%', label: 'Revenue Share', sub: 'You keep 70% of every ticket sold' },
  { value: 'HLS', label: 'Low-Latency Streaming', sub: '2-second fragments, global delivery' },
  { value: 'M-Pesa', label: 'Payment Method', sub: 'Instant STK Push, no card required' },
  { value: '0 KES', label: 'Setup Cost', sub: 'No monthly fees — pay-per-event model' },
]

const steps = [
  { title: 'Create Your Account', desc: 'Register as a host — free, instant, no credit card required.' },
  { title: 'Schedule Your Event', desc: 'Add event details, set your ticket price (or make it free), and upload a banner image.' },
  { title: 'Go Live', desc: 'Stream from OBS, a phone, or any encoder. Copy your stream key and start broadcasting.' },
  { title: 'Get Paid', desc: 'Attendees pay via M-Pesa. Revenue tracked live, paid out within 48 hours.' },
]

const revenueItems = [
  { label: 'You receive', value: '70% of ticket sales', highlight: true },
  { label: 'Platform fee', value: '30%', highlight: false },
  { label: 'Payment processing', value: 'Included in platform fee', highlight: false },
  { label: 'Streaming infrastructure', value: 'Included in platform fee', highlight: false },
  { label: 'Fraud & piracy protection', value: 'Included in platform fee', highlight: false },
  { label: 'Payout timeline', value: '48 hrs after event ends', highlight: false },
]

const features = [
  {
    icon: Calendar,
    title: 'Event Scheduling',
    desc: 'Create events with custom pricing, descriptions, and banner images. Schedule weeks in advance.',
  },
  {
    icon: Tv,
    title: 'HLS Live Streaming',
    desc: 'Stream from any encoder (OBS, vMix, Wirecast). We handle the HLS conversion and delivery.',
  },
  {
    icon: Smartphone,
    title: 'M-Pesa Payments',
    desc: 'Viewers pay via M-Pesa STK Push — the most popular payment method in East Africa.',
  },
  {
    icon: BarChart2,
    title: 'Real-Time Analytics',
    desc: 'Live viewer counts, peak attendance, total revenue, and ticket sales — all in your dashboard.',
  },
  {
    icon: DollarSign,
    title: 'Revenue Dashboard',
    desc: 'Track earnings per event. See your 70% share in real time as tickets are purchased.',
  },
  {
    icon: Shield,
    title: 'Anti-Piracy Built In',
    desc: 'Unique tokens, IP locking, device fingerprinting, and VPN detection protect every stream.',
  },
]

const shields = [
  'Device Fingerprinting',
  'IP Address Lock',
  'VPN Detection',
  'Single-Session Tokens',
]

// Earnings calculator
const calc = ref({ price: 500, viewers: 500 })
const grossRevenue = computed(() => calc.value.price * calc.value.viewers)
const platformFee = computed(() => Math.round(grossRevenue.value * 0.3))
const yourEarnings = computed(() => grossRevenue.value - platformFee.value)
</script>
