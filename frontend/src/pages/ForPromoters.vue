<template>
  <div class="bg-bg text-white">

    <!-- ============================================================ -->
    <!-- HERO                                                          -->
    <!-- ============================================================ -->
    <section class="relative py-24 lg:py-36 overflow-hidden">
      <div class="absolute inset-0 pointer-events-none">
        <div
          class="absolute top-0 right-0 w-[700px] h-[700px] rounded-full opacity-5"
          style="background: radial-gradient(circle, #E8002D 0%, transparent 65%); filter: blur(100px); transform: translate(30%, -30%);"
        />
        <div
          class="absolute bottom-0 left-0 w-[500px] h-[500px] rounded-full opacity-4"
          style="background: radial-gradient(circle, #FF6B00 0%, transparent 65%); filter: blur(120px); transform: translate(-20%, 20%);"
        />
      </div>

      <div class="relative max-w-6xl mx-auto px-6 sm:px-8 lg:px-16">
        <div class="grid lg:grid-cols-2 gap-16 items-center">
          <!-- Left: Copy -->
          <div class="space-y-8">
            <div class="inline-flex items-center gap-2 bg-accent-red/10 border border-accent-red/25
                        text-accent-red text-xs font-bold uppercase tracking-widest px-3 py-1.5 rounded">
              For Promoters
            </div>

            <h1 class="font-display text-6xl md:text-7xl uppercase text-white leading-none tracking-wide">
              Host Live.<br />
              Get Paid.<br />
              <span class="text-accent-red">Keep 70%.</span>
            </h1>

            <p class="text-text-muted text-lg leading-relaxed max-w-md">
              Schedule live events, set your ticket price, and collect M-Pesa payments directly.
              We handle the stream infrastructure — you keep the majority.
            </p>

            <div class="flex flex-wrap items-center gap-4 pt-2">
              <RouterLink
                to="/try"
                class="inline-flex items-center gap-2 px-8 py-4 rounded-lg bg-accent-red
                       hover:bg-accent-red-hover text-white font-bold text-base
                       transition-all duration-200 active:scale-95"
              >
                <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                    d="M15 10l4.553-2.276A1 1 0 0121 8.677v6.646a1 1 0 01-1.447.894L15 14M5 18h8a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v8a2 2 0 002 2z" />
                </svg>
                Try a Free Test Stream
              </RouterLink>
              <RouterLink
                to="/register"
                class="inline-flex items-center gap-2 px-6 py-4 rounded-lg border border-white/20
                       hover:border-white/40 text-white text-base font-medium
                       transition-all duration-200 hover:bg-white/5"
              >
                Create Account
              </RouterLink>
            </div>

            <div class="flex flex-wrap items-center gap-x-6 gap-y-2 text-text-muted text-sm pt-1">
              <span class="flex items-center gap-1.5">
                <Check class="w-4 h-4 text-status-success shrink-0" />
                No monthly fees
              </span>
              <span class="flex items-center gap-1.5">
                <Check class="w-4 h-4 text-status-success shrink-0" />
                No account needed to test
              </span>
              <span class="flex items-center gap-1.5">
                <Check class="w-4 h-4 text-status-success shrink-0" />
                M-Pesa payouts
              </span>
            </div>
          </div>

          <!-- Right: Key numbers -->
          <div class="grid grid-cols-2 gap-4">
            <div class="bg-bg-elevated border border-white/5 rounded-xl p-6 flex flex-col gap-2">
              <span class="font-display text-5xl text-accent-red leading-none">70%</span>
              <span class="text-white font-semibold text-sm">Revenue to you</span>
              <span class="text-text-muted text-xs leading-relaxed">Every ticket sold. No negotiations.</span>
            </div>
            <div class="bg-bg-elevated border border-white/5 rounded-xl p-6 flex flex-col gap-2">
              <span class="font-display text-5xl text-white leading-none">&lt;24h</span>
              <span class="text-white font-semibold text-sm">M-Pesa payout</span>
              <span class="text-text-muted text-xs leading-relaxed">Direct to your mobile money.</span>
            </div>
            <div class="bg-bg-elevated border border-white/5 rounded-xl p-6 flex flex-col gap-2 col-span-2">
              <span class="font-display text-5xl text-accent-orange leading-none">KES 0</span>
              <span class="text-white font-semibold text-sm">Upfront cost</span>
              <span class="text-text-muted text-xs leading-relaxed">No subscription. Pay nothing until you earn.</span>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- ============================================================ -->
    <!-- REVENUE CALCULATOR                                            -->
    <!-- ============================================================ -->
    <section class="py-24 bg-bg-surface border-y border-white/5">
      <div class="max-w-6xl mx-auto px-6 sm:px-8 lg:px-16">
        <div class="grid lg:grid-cols-2 gap-16 items-start">
          <!-- Copy -->
          <div class="space-y-5 lg:pt-4">
            <p class="text-accent-red text-sm font-bold uppercase tracking-widest">Earnings</p>
            <h2 class="font-display text-4xl md:text-5xl uppercase text-white leading-none">
              Calculate your<br />potential earnings
            </h2>
            <p class="text-text-muted text-base leading-relaxed">
              Drag the sliders to estimate what you'd earn from your event.
              70% of every ticket goes straight to you — always.
            </p>

            <div class="space-y-3 pt-4">
              <div class="flex items-center gap-3 text-text-muted text-sm">
                <div class="w-8 h-8 rounded-lg bg-accent-red/10 border border-accent-red/20 flex items-center justify-center shrink-0">
                  <Ticket class="w-4 h-4 text-accent-red" />
                </div>
                Set any price from KES 0 (free) to KES 100,000
              </div>
              <div class="flex items-center gap-3 text-text-muted text-sm">
                <div class="w-8 h-8 rounded-lg bg-accent-red/10 border border-accent-red/20 flex items-center justify-center shrink-0">
                  <Smartphone class="w-4 h-4 text-accent-red" />
                </div>
                Paid to M-Pesa within 24 hours of the event
              </div>
              <div class="flex items-center gap-3 text-text-muted text-sm">
                <div class="w-8 h-8 rounded-lg bg-accent-red/10 border border-accent-red/20 flex items-center justify-center shrink-0">
                  <BarChart2 class="w-4 h-4 text-accent-red" />
                </div>
                Full breakdown visible in your host dashboard
              </div>
            </div>
          </div>

          <!-- Calculator -->
          <div class="bg-bg-elevated rounded-xl border border-white/5 overflow-hidden">
            <div class="p-8 space-y-8">
              <!-- Slider 1 -->
              <div>
                <div class="flex items-center justify-between mb-3">
                  <label class="text-text-muted text-xs uppercase tracking-widest font-semibold">
                    Ticket Price
                  </label>
                  <span class="font-display text-3xl text-white leading-none">
                    KES {{ calc.price.toLocaleString() }}
                  </span>
                </div>
                <input
                  v-model.number="calc.price"
                  type="range"
                  min="100"
                  max="5000"
                  step="50"
                  class="w-full h-2 appearance-none bg-white/10 cursor-pointer rounded-full"
                  style="accent-color: #E8002D;"
                />
                <div class="flex justify-between text-text-muted text-xs mt-2">
                  <span>KES 100</span>
                  <span>KES 5,000</span>
                </div>
              </div>

              <!-- Slider 2 -->
              <div>
                <div class="flex items-center justify-between mb-3">
                  <label class="text-text-muted text-xs uppercase tracking-widest font-semibold">
                    Viewers
                  </label>
                  <span class="font-display text-3xl text-white leading-none">
                    {{ calc.viewers.toLocaleString() }}
                  </span>
                </div>
                <input
                  v-model.number="calc.viewers"
                  type="range"
                  min="10"
                  max="10000"
                  step="10"
                  class="w-full h-2 appearance-none bg-white/10 cursor-pointer rounded-full"
                  style="accent-color: #E8002D;"
                />
                <div class="flex justify-between text-text-muted text-xs mt-2">
                  <span>10</span>
                  <span>10,000</span>
                </div>
              </div>
            </div>

            <!-- Results -->
            <div class="border-t border-white/8 bg-black/20 p-8 space-y-4">
              <div class="flex items-center justify-between text-sm">
                <span class="text-text-muted">Gross revenue</span>
                <span class="text-white font-medium">KES {{ grossRevenue.toLocaleString() }}</span>
              </div>
              <div class="flex items-center justify-between text-sm">
                <span class="text-text-muted">Platform fee (30%)</span>
                <span class="text-text-muted">− KES {{ platformFee.toLocaleString() }}</span>
              </div>
              <div class="flex items-center justify-between pt-4 border-t border-white/8">
                <span class="text-white font-bold">Your payout</span>
                <span class="font-display text-4xl text-status-success leading-none">
                  KES {{ yourEarnings.toLocaleString() }}
                </span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- ============================================================ -->
    <!-- HOW IT WORKS                                                  -->
    <!-- ============================================================ -->
    <section id="how-it-works" class="py-24">
      <div class="max-w-6xl mx-auto px-6 sm:px-8 lg:px-16">
        <div class="mb-16">
          <p class="text-accent-red text-sm font-bold uppercase tracking-widest mb-3">Process</p>
          <h2 class="font-display text-4xl md:text-5xl uppercase text-white leading-none">
            From stream key<br />to M-Pesa
          </h2>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-4 gap-6">
          <div
            v-for="(step, i) in steps"
            :key="step.title"
            class="relative bg-bg-elevated border border-white/5 rounded-xl p-6
                   hover:border-accent-red/20 transition-colors duration-200"
          >
            <div class="flex items-start gap-4 mb-4">
              <span class="font-display text-5xl text-accent-red/25 leading-none select-none">
                {{ String(i + 1).padStart(2, '0') }}
              </span>
            </div>
            <div class="w-10 h-10 rounded-lg bg-accent-red/10 border border-accent-red/15
                        flex items-center justify-center mb-4">
              <component :is="step.icon" class="w-5 h-5 text-accent-red" />
            </div>
            <h3 class="text-white font-bold text-sm mb-2">{{ step.title }}</h3>
            <p class="text-text-muted text-sm leading-relaxed">{{ step.desc }}</p>
          </div>
        </div>
      </div>
    </section>

    <!-- ============================================================ -->
    <!-- PLATFORM FEATURES                                             -->
    <!-- ============================================================ -->
    <section class="py-24 bg-bg-surface border-y border-white/5">
      <div class="max-w-6xl mx-auto px-6 sm:px-8 lg:px-16">
        <div class="mb-14">
          <p class="text-accent-red text-sm font-bold uppercase tracking-widest mb-3">Platform</p>
          <h2 class="font-display text-4xl md:text-5xl uppercase text-white leading-none">
            Built for live events
          </h2>
          <p class="text-text-muted text-base mt-4">
            Everything you need to run a ticketed live stream.
          </p>
        </div>

        <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-px bg-white/5 rounded-xl overflow-hidden border border-white/5">
          <div
            v-for="feature in platformFeatures"
            :key="feature.title"
            class="bg-bg-surface p-7 flex flex-col gap-4 hover:bg-bg-elevated transition-colors duration-200 group"
          >
            <div class="w-11 h-11 rounded-lg bg-accent-red/10 border border-accent-red/15
                        flex items-center justify-center shrink-0 group-hover:bg-accent-red/20 transition-colors duration-200">
              <component :is="feature.icon" class="w-5 h-5 text-accent-red" />
            </div>
            <div>
              <h3 class="text-white font-bold text-sm mb-1.5">{{ feature.title }}</h3>
              <p class="text-text-muted text-sm leading-relaxed">{{ feature.desc }}</p>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- ============================================================ -->
    <!-- CTA                                                           -->
    <!-- ============================================================ -->
    <section class="py-24">
      <div class="max-w-6xl mx-auto px-6 sm:px-8 lg:px-16">
        <div class="relative bg-bg-elevated border border-white/8 rounded-xl overflow-hidden">
          <!-- Background accent -->
          <div class="absolute inset-0 pointer-events-none">
            <div
              class="absolute -top-1/2 -right-1/4 w-[500px] h-[500px] rounded-full opacity-8"
              style="background: radial-gradient(circle, #E8002D 0%, transparent 70%); filter: blur(80px);"
            />
          </div>

          <div class="relative p-10 md:p-16">
            <div class="max-w-2xl">
              <h2 class="font-display text-5xl md:text-6xl uppercase text-white leading-none mb-4">
                Ready to go live?
              </h2>
              <p class="text-text-muted text-base leading-relaxed mb-8">
                Test the platform with no account — push a real stream, see the dashboard, confirm it
                works for your setup. Then create an account and start charging.
              </p>

              <div class="flex flex-wrap items-center gap-4 mb-10">
                <RouterLink
                  to="/try"
                  class="inline-flex items-center gap-2 px-7 py-3.5 rounded-lg bg-accent-red
                         hover:bg-accent-red-hover text-white font-bold text-sm
                         transition-all duration-200 active:scale-95"
                >
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                      d="M15 10l4.553-2.276A1 1 0 0121 8.677v6.646a1 1 0 01-1.447.894L15 14M5 18h8a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v8a2 2 0 002 2z" />
                  </svg>
                  Try a Free Stream
                </RouterLink>
                <RouterLink
                  to="/register"
                  class="inline-flex items-center gap-2 px-7 py-3.5 rounded-lg border border-white/20
                         hover:border-white/40 text-white text-sm font-medium
                         transition-all hover:bg-white/5"
                >
                  Create Account
                </RouterLink>
              </div>

              <div class="flex flex-wrap items-center gap-6 text-text-muted text-sm">
                <span class="flex items-center gap-1.5">
                  <Check class="w-4 h-4 text-status-success shrink-0" />
                  No monthly fees
                </span>
                <span class="flex items-center gap-1.5">
                  <Check class="w-4 h-4 text-status-success shrink-0" />
                  No credit card required
                </span>
                <span class="flex items-center gap-1.5">
                  <Check class="w-4 h-4 text-status-success shrink-0" />
                  70% revenue share, always
                </span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>

  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { RouterLink } from 'vue-router'
import {
  Ticket, CheckCircle, Smartphone, Monitor, Shield, BarChart2,
  Target, Check, Radio, Wallet,
} from 'lucide-vue-next'

const steps = [
  {
    icon: Monitor,
    title: 'Create your event',
    desc: 'Set a title, date, ticket price, and get your RTMP stream key.',
  },
  {
    icon: Radio,
    title: 'Go live',
    desc: 'Push from OBS, your phone, or any RTMP encoder. We handle HLS delivery.',
  },
  {
    icon: Ticket,
    title: 'Viewers pay via M-Pesa',
    desc: 'STK Push to their phone. Instant access granted on payment confirmation.',
  },
  {
    icon: Wallet,
    title: 'You get paid',
    desc: '70% lands in your M-Pesa within 24 hours of the event ending.',
  },
]

const platformFeatures = [
  {
    icon: Monitor,
    title: 'HD Live Streaming',
    desc: 'Push RTMP from OBS or any encoder. We convert to HLS and serve it globally.',
  },
  {
    icon: Smartphone,
    title: 'M-Pesa Payments',
    desc: 'Viewers pay via STK Push. No card, no friction — just a PIN on their phone.',
  },
  {
    icon: Shield,
    title: 'Anti-Piracy',
    desc: 'Device fingerprinting and IP locking ensure only paying viewers can watch.',
  },
  {
    icon: BarChart2,
    title: 'Live Analytics',
    desc: 'Viewer counts, ticket sales, and revenue — live in your host dashboard.',
  },
  {
    icon: Target,
    title: 'Flexible Pricing',
    desc: 'Set any ticket price from KES 0 (free event) up to KES 100,000.',
  },
  {
    icon: CheckCircle,
    title: 'Transparent Payouts',
    desc: 'Gross, platform fee, and your cut — itemised per event, no hidden charges.',
  },
]

const calc = ref({ price: 500, viewers: 500 })
const grossRevenue = computed(() => calc.value.price * calc.value.viewers)
const platformFee = computed(() => Math.round(grossRevenue.value * 0.3))
const yourEarnings = computed(() => grossRevenue.value - platformFee.value)
</script>
